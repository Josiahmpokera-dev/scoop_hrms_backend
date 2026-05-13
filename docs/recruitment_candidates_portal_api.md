# Recruitment — Candidates & portal APIs (UI integration)

This document is **only** about listing/searching candidates, applications, interviews from the HR portal, plus realtime updates. Job openings and public apply flows stay in `docs/recruitment_job_openings_api.md`.

**Base path (authenticated):** `/api/v1/recruitment`  
**Public apply (no auth):** `/api/v1/recruitment/public`

---

## Authentication & permissions

Every `/api/v1/recruitment/*` route (except `/public/*`) requires:

`Authorization: Bearer <access_token>`

| Permission code | Typical use |
|-----------------|-------------|
| `recruitment:read` | Job openings, interviews list/detail, WebSocket |
| `recruitment:manage_candidates` | Candidates list/detail, application actions, schedule interviews, feedback, HR decisions |

---

## Standard JSON envelope

Most endpoints return:

```json
{
  "success": true,
  "message": "…",
  "data": …,
  "meta": { … }
}
```

**Pagination `meta`** (when present) follows `internal/utils/response/response.Meta`:

| Field | Meaning |
|-------|---------|
| `page` | Current page (1-based) |
| `per_page` | Page size (same as request `limit`) |
| `total` | Total rows matching filters |
| `total_pages` | Ceiling of total ÷ limit |

List endpoints cap **`limit` at 100**; default **`page=1`**, **`limit=20`** where applicable.

---

## Pipeline stages (`JobApplication.stage`)

Use these exact strings for filters (`stage=` query) and for PATCH bodies:

| Value | UI suggestion |
|-------|------------------|
| `Applied` | New applicant |
| `Screening` | CV / shortlist review |
| `Interview` | Interview scheduled or in progress |
| `InterviewCompleted` | Interviews finished; decision pending |
| `Offer` | Offer extended |
| `Hired` | Accepted |
| `Rejected` | Closed — not proceeding |
| `TalentPool` | Nurture / future roles |

---

## 1. Unified candidates list (global search & filters)

**GET** `/api/v1/recruitment/candidates`  
**Permission:** `recruitment:manage_candidates`

| Query | Required | Description |
|-------|----------|-------------|
| `jobId` | No | Opening UUID — only candidates with an application on this job |
| `stage` | No | One of the pipeline stages above (e.g. `Screening`) |
| `q` or `search` | No | Keyword: matches **first name, last name, email, phone** (case-insensitive) |
| `page` | No | Default `1` |
| `limit` | No | Default `20`, max `100` |

**Response:** `{ success, message, data: Candidate[], meta }`  
Each `Candidate` includes `applications[]` (preload). Use this for the **main CRM-style table** with filters.

---

## 2. Candidate profile (detail drawer / page)

**GET** `/api/v1/recruitment/candidates/:candidateId`  
**Permission:** `recruitment:manage_candidates`

Returns one candidate with relations loaded by the repository (`applications`, nested job openings, interviews, offers where implemented).

**Important for CV display**

- The API stores **`resumeUrl`** — a **URL** string (often S3/MinIO after upload), **not** raw base64.
- Public apply accepts base64 or URL in **`resume`**; the server uploads to storage and persists **`resumeUrl`** on the candidate.
- The UI should **`window.open(resumeUrl)`** or embed/link — do not expect base64 on GET.

---

## 3. Candidates by interview round (stage 1 / 2 / 3)

**GET** `/api/v1/recruitment/candidates/interview-stages`  
**Permission:** `recruitment:manage_candidates`

| Query | Required | Description |
|-------|----------|-------------|
| `stage` or `round` | **Yes** | **`1`**, **`2`**, or **`3`** — interview round filter |
| `jobId` | No | Limit to one opening |
| `q` or `search` | No | Same keyword search as `/candidates` |
| `page`, `limit` | No | Same pagination rules |

**Response:** `{ success, message, data: Candidate[], meta }`

Use this for **tabs**: “Round 1”, “Round 2”, “Round 3” shortlists (backend ties candidates to interviews for that round).

---

## 4. Applications for one job opening (opening detail tab)

**GET** `/api/v1/recruitment/openings/:openingId/applications`  
**Permission:** `recruitment:read`

| Query | Description |
|-------|-------------|
| `stage` | Optional pipeline filter |
| `page`, `limit` | Pagination |

**Response shape (note nesting):**

```json
{
  "success": true,
  "message": "Applications retrieved successfully",
  "data": {
    "data": [ /* JobApplication rows with Candidate preloaded */ ],
    "meta": {
      "total": 0,
      "page": 1,
      "limit": 20,
      "total_pages": 1,
      "job_opening_id": "…",
      "stage": ""
    }
  }
}
```

In the UI, read rows from **`response.data.data`** and pagination from **`response.data.meta`**.

---

## 5. Update stage (manual pipeline move)

Two routes call the **same** service logic.

### PATCH `/api/v1/recruitment/candidates/:id/stage`

### PATCH `/api/v1/recruitment/applications/:id/stage`

**Permission:** `recruitment:manage_candidates`

**Body (`UpdateStageRequest`):**

```json
{
  "stage": "Screening",
  "notes": "Optional internal note",
  "jobId": ""
}
```

**Critical — `:id` meaning**

- The implementation loads the row with **`GetApplicationByID(id)`**.  
- So **`:id` is the job application UUID**, not the candidate UUID, even on the `/candidates/:id/stage` path.

Use **`application.id`** from list/detail responses when calling PATCH.

**WebSocket:** emits `candidate.application.stage_updated` with `applicationId`, `candidateId`, `jobId`, `stage`.

---

## 6. Application actions (accept / reject / move to interview round 1)

**POST** `/api/v1/recruitment/applications/:applicationId/action`  
**Permission:** `recruitment:manage_candidates`

**Body (`ApplicationActionRequest`):**

| Field | When |
|-------|------|
| `action` | **`accept`** \| **`reject`** \| **`move_to_interview_stage_one`** |
| `notes` | Optional; appended to application notes |

For **`move_to_interview_stage_one`** only:

| Field | Rule |
|-------|------|
| `scheduledDate`, `scheduledTime` | **Required** |
| `interviewType` | Optional; default **`Stage 1`** |
| `duration`, `mode` | Optional |
| `interviewerEmployeeIds` | Preferred — employee IDs from HR directory |
| `interviewers` | Optional email list (legacy / notifications) |

**Behaviour**

- **`accept`** → `stage` = **`Screening`**
- **`reject`** → **`Rejected`**
- **`move_to_interview_stage_one`** → schedules **round 1** interview, then sets `stage` = **`Interview`**

**WebSocket:** `candidate.application.actioned` with `action`, `stage`, `applicationId`, `candidateId`, `jobId`.

---

## 7. Move application to talent pool

**POST** `/api/v1/recruitment/applications/:applicationId/move-to-talent-pool`  
**Permission:** `recruitment:manage_candidates`

See `MoveApplicationToTalentPoolRequest` in code — optional `syncTalentPoolRow`, `internalPoolNotes`, etc.

---

## 8. Interviews (global list & detail)

### Global list

**GET** `/api/v1/recruitment/interviews`  
**Permission:** `recruitment:read`

| Query | Description |
|-------|-------------|
| `stage` or `round` | **`1`**, **`2`**, or **`3`** — omit for all rounds |
| `status` | e.g. `Scheduled`, `PanelCompleted`, `HrApprovedNextRound` |
| `jobId` | Opening UUID |
| `candidateId` | Candidate UUID |
| `page`, `limit` | Pagination |

**Response:** `{ success, message, data: Interview[], meta }` — flat array in **`data`**, no extra nesting.

### Interview detail

**GET** `/api/v1/recruitment/interviews/:interviewId`  
**Permission:** `recruitment:read`

### Interviews for one opening

**GET** `/api/v1/recruitment/openings/:openingId/interviews`  
**Permission:** `recruitment:read`

Returns an array (no pagination in handler).

### Schedule / feedback / HR decision

- **POST** `/interviews` — schedule (`ScheduleInterviewRequest`, `round` 1–3, `interviewerEmployeeIds`, `isReinterview`)
- **POST** `/interviews/:id/feedback` — panel feedback
- **POST** `/interviews/:id/hr-decision` — HR gate (`proceed_next_round` | `schedule_reinterview` | `reject_candidate`), **`remarks`** required; **`decidedByEmployeeId`** if JWT user has no employee row

Full field lists and status workflow: see `docs/recruitment_job_openings_api.md` § Interviews.

---

## 9. Realtime — WebSocket

**GET** `/api/v1/recruitment/ws/candidates`  
**Permission:** `recruitment:read` (via middleware)

The HTTP upgrade runs **after** `AuthMiddleware`, which currently accepts **only** the header:

`Authorization: Bearer <token>`

**Browser note:** the native `WebSocket` API cannot set custom headers. Options:

- Use a **desktop/mobile app** or server-side client that can send the header, or
- Put a **small same-origin proxy** that adds `Authorization`, or
- Extend the server later to accept `?access_token=` on the upgrade (not implemented today).

**Message envelope:**

```json
{
  "type": "candidate.application.created",
  "timestamp": "2026-05-13T12:00:00Z",
  "data": { }
}
```

**Event types emitted today**

| `type` | When | Typical `data` keys |
|--------|------|---------------------|
| `candidate.application.created` | New application persisted | `applicationId`, `candidateId`, `jobId`, `stage` |
| `candidate.application.stage_updated` | PATCH stage | `applicationId`, `candidateId`, `jobId`, `stage` |
| `candidate.application.actioned` | POST action | `applicationId`, `candidateId`, `jobId`, `action`, `stage` |
| `interview.scheduled` | New interview row | `interviewId`, `applicationId`, `candidateId`, `jobId`, `round`, `status` |
| `interview.feedback_submitted` | Panel submitted feedback | (see service broadcast payload) |
| `interview.hr_decision` | HR decision recorded | `interviewId`, `applicationId`, `candidateId`, `jobId`, `round`, `decision`, `status`, `remarks`, `applicationStage` |

**UI pattern**

1. After login, open the socket (with auth as your stack allows).
2. On any `candidate.*` event: refetch **`GET /candidates`** and/or **`GET /openings/:id/applications`** for the active job.
3. On `interview.*`: refetch **`GET /interviews`** and/or **`GET /candidates/interview-stages`**.

---

## 10. Quick UI route map

| Screen | Primary endpoints |
|--------|-------------------|
| All candidates | `GET /candidates` |
| Candidate 360° | `GET /candidates/:id` |
| Job applicants tab | `GET /openings/:id/applications` (note nested `data.data`) |
| Pipeline actions | `POST /applications/:id/action`, `PATCH /applications/:id/stage` |
| Interview rounds tabs | `GET /candidates/interview-stages?stage=1|2|3` |
| Interview calendar / table | `GET /interviews` |
| Live refresh | `GET /ws/candidates` + refetch lists |

---

## 11. Public apply (external careers site only)

Not authenticated:

- **GET** `/api/v1/recruitment/public/openings/by-token/:token`
- **POST** `/api/v1/recruitment/public/apply`

Details and CV encoding rules: **`docs/recruitment_job_openings_api.md`**.
