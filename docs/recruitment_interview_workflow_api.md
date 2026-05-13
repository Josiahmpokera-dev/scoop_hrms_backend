# Structured interview workflow API

This module implements a **stage-driven interview process** aligned with the existing HRMS stack (Gin, JWT auth, `recruitment:*` permissions, employees linked to users, WebSocket hub events). It lives **alongside** the legacy `Interview` round model (`POST /recruitment/interviews`, panel feedback, HR gate) documented in `docs/recruitment_job_openings_api.md`. Use **one** model per candidate/application according to your hiring policy.

**Base path:** `/api/v1/recruitment/interview-workflows`

**Auth:** `Authorization: Bearer <token>` (same as other recruitment routes).

---

## Design summary

| Concept | Behaviour |
|---------|-----------|
| **Definition** | Per **job opening**: ordered **stages**, each with a **structured scoring rubric** (dimensions with min/max). Optional **CEO** in the final approval chain. |
| **Process** | One running workflow per **job application** (enforced). Snapshots freeze rubric/titles at start. |
| **Stage attempt** | Each scheduled run is a row: **one assigned interviewer** (`employees.employee_id`). **Sequential**: HR schedules stage *N+1* only after stage *N* is **Completed** with outcome **`proceed`**. |
| **Attendance** | Recorded once per attempt (`pending` → `present` \| `absent` \| `excused`). Cannot change after **evaluation** is submitted. |
| **Interview start** | Assignee only; requires **`present`**. |
| **Evaluation** | Assignee only; **mandatory remarks**; **scores** must match snapshot dimensions; outcome **`proceed`** \| **`repeat_stage`** \| **`fail`**. Locks the attempt (immutable). |
| **Repeat / reschedule** | **`repeat_stage`**: new attempt row, same stage sequence, `attemptNumber+1`. **`absence` + `reschedule`**: previous attempt **`Superseded`**, new attempt created. **`fail_candidate`**: process **`failed`**, application **`Rejected`**. |
| **Final approvals** | After **last** stage **`proceed`**: sequential **HR → Manager → (optional CEO)**. Each step only by **assigned** approver; any **reject** rejects the candidate; all **approve** sets process **`approved`** and application **`InterviewCompleted`**. |
| **Audit** | Append-only **`interview_workflow_audit_events`** + immutable fields on completed attempts. **`GET .../timeline`** returns process + audit. |

**Authorization model**

- **HR configuration & scheduling:** `recruitment:manage_candidates` (definitions, start process, schedule next stage).
- **Interviewer / approver actions:** `recruitment:read` **plus** the server checks the JWT user’s **employee** (`employees.user_id` → `employee_id`) matches **`assignedEmployeeId`** on the attempt or approval row.

---

## Permissions

| Permission | Routes |
|------------|--------|
| `recruitment:manage_candidates` | Create/update definitions, list definitions, start process, schedule stage |
| `recruitment:read` | Get process, timeline, definition; assignee/approver POST actions (enforced in service) |

---

## Standard response

Same envelope as the rest of the API: `{ "success", "message", "data", "error?", "meta?" }`.

---

## 1. Definitions (HR configures stages + rubric)

### List definitions for a job

**GET** `/api/v1/recruitment/interview-workflows/definitions?jobOpeningId=<uuid>`

### Create definition

**POST** `/api/v1/recruitment/interview-workflows/definitions`  
**Permission:** `recruitment:manage_candidates`

```json
{
  "jobOpeningId": "<opening-uuid>",
  "name": "Engineering standard",
  "includeCeoApproval": true,
  "stages": [
    {
      "sequence": 1,
      "title": "Technical screen",
      "description": "Depth on stack and problem solving.",
      "scoringDimensions": [
        { "key": "technical", "label": "Technical depth", "min": 1, "max": 5, "required": true },
        { "key": "communication", "label": "Communication", "min": 1, "max": 5, "required": true }
      ]
    },
    {
      "sequence": 2,
      "title": "Culture & leadership",
      "description": "",
      "scoringDimensions": [
        { "key": "values", "label": "Values alignment", "min": 1, "max": 5, "required": true }
      ]
    }
  ]
}
```

### Get definition

**GET** `/api/v1/recruitment/interview-workflows/definitions/:id`  
**Permission:** `recruitment:read`

### Update definition (full replace of stages)

**PUT** `/api/v1/recruitment/interview-workflows/definitions/:id`  
**Permission:** `recruitment:manage_candidates`  
Body: same shape as create (including `jobOpeningId` — must match existing).

> **Note:** Running processes use **snapshots** created at start; editing a definition does not rewrite in-flight processes.

---

## 2. Start a process (bind to one application)

**POST** `/api/v1/recruitment/interview-workflows/processes`  
**Permission:** `recruitment:manage_candidates`

- `definitionId` must belong to the **same job opening** as the application.
- **At most one** process per `applicationId`.
- If `includeCeoApproval` is **true** on the definition, **`finalApprovals.ceoEmployeeId`** is **required** and must exist in `employees`.

```json
{
  "applicationId": "<application-uuid>",
  "definitionId": "<definition-uuid>",
  "firstStage": {
    "assignedEmployeeId": "EMP001",
    "scheduledDate": "2026-06-15",
    "scheduledTime": "09:30",
    "durationMin": 60,
    "mode": "Video"
  },
  "finalApprovals": {
    "hrEmployeeId": "EMP-HR-01",
    "managerEmployeeId": "EMP-MGR-02",
    "ceoEmployeeId": "EMP-CEO-01"
  }
}
```

Creates: process (`in_progress`), **stage snapshots**, first **stage attempt** (sequence `1`, `Scheduled`), audit events, WebSocket `interview_workflow.process_started`.

---

## 3. Read process & timeline (UI)

**GET** `/api/v1/recruitment/interview-workflows/processes/:id`  
**GET** `/api/v1/recruitment/interview-workflows/processes/by-application/:applicationId`  
**Permission:** `recruitment:read`

Response includes `snapshots`, `attempts` (ordered by `stageSequence`, `attemptNumber`), `approvals` (when in final phase).

**GET** `/api/v1/recruitment/interview-workflows/processes/:id/timeline`  
Returns `{ "process": { ... }, "audit": [ ... ] }` for a **card + timeline** UI.

---

## 4. Schedule the next stage (HR)

**POST** `/api/v1/recruitment/interview-workflows/processes/:processId/stages`  
**Permission:** `recruitment:manage_candidates`

Rules:

- Process status must be **`in_progress`**.
- For `stageSequence > 1`, the **previous** stage’s latest non-superseded attempt must be **`Completed`** with outcome **`proceed`**.
- No **open** attempt (`Scheduled` or `InProgress`) may exist for that `stageSequence`.

```json
{
  "stageSequence": 2,
  "assignedEmployeeId": "EMP042",
  "scheduledDate": "2026-06-20",
  "scheduledTime": "14:00",
  "durationMin": 45,
  "mode": "On-site"
}
```

---

## 5. Assignee: attendance → start → evaluation

All use **stage-attempt** id from `process.attempts[].id`.

### Record attendance (once)

**POST** `/api/v1/recruitment/interview-workflows/stage-attempts/:attemptId/attendance`  
**Permission:** `recruitment:read` + must be **assigned** interviewer.

```json
{ "status": "present", "notes": "Candidate joined Zoom on time." }
```

`status`: `present` | `absent` | `excused`.

Cannot change after evaluation is finalized.

### Mark in progress

**POST** `/api/v1/recruitment/interview-workflows/stage-attempts/:attemptId/start`  
**Permission:** `recruitment:read` + assignee.  
Requires **`present`** attendance; moves **`Scheduled` → `InProgress`**.

### Submit evaluation (locks attempt)

**POST** `/api/v1/recruitment/interview-workflows/stage-attempts/:attemptId/evaluation`  
**Permission:** `recruitment:read` + assignee.  
Requires **`InProgress`** and **`present`**.

```json
{
  "scores": { "technical": 4, "communication": 5 },
  "remarks": "Solid system design answers; minor gaps in observability.",
  "outcome": "proceed"
}
```

**Outcomes**

| `outcome` | Effect |
|-----------|--------|
| `proceed` | Attempt **`Completed`**. If **last** stage → create **final approval** rows; process → **`pending_final_approvals`**. Else HR schedules next stage. |
| `repeat_stage` | Current attempt **`Completed`**; new attempt same `stageSequence`, higher `attemptNumber`. Requires **`nextScheduledDate`**, **`nextScheduledTime`**. |
| `fail` | Process **`failed`**; application **`Rejected`**. |

### Absence handling (after marking `absent` or `excused`)

**POST** `/api/v1/recruitment/interview-workflows/stage-attempts/:attemptId/absence-action`  
**Permission:** `recruitment:read` + assignee.

```json
{
  "action": "reschedule",
  "nextScheduledDate": "2026-06-22",
  "nextScheduledTime": "11:00",
  "durationMin": 60,
  "mode": "Video",
  "notes": "Candidate notified; agreed new slot."
}
```

| `action` | Effect |
|----------|--------|
| `fail_candidate` | Attempt **`Missed`**; process **`failed`**; application **`Rejected`**. |
| `reschedule` | Attempt **`Superseded`**; new **`Scheduled`** attempt (same stage, next `attemptNumber`). |

---

## 6. Final approval chain

**POST** `/api/v1/recruitment/interview-workflows/final-approvals/:approvalId/decide`  
**Permission:** `recruitment:read` + must match **`assignedEmployeeId`** on that approval row.

```json
{ "decision": "approve", "remarks": "Aligned with panel recommendation." }
```

- Steps run in **`stepOrder`**: HR (1), Manager (2), CEO (3) if enabled.
- Cannot act until **all lower** steps are **`approved`**.
- **`reject`** on any step → process **`rejected`**, application **`Rejected`**.
- Last **`approve`** → process **`approved`**, application **`InterviewCompleted`**.

---

## Process & attempt statuses (UI)

**Process:** `in_progress` | `pending_final_approvals` | `approved` | `rejected` | `failed`

**Attempt:** `Scheduled` | `InProgress` | `Completed` | `Missed` | `Superseded`

Map to UI cards:

| UI label | API |
|----------|-----|
| Scheduled | `Scheduled`, attendance `pending` |
| In progress | `InProgress` |
| Completed | `Completed` |
| Missed | `Missed` (failed no-show path) |
| Repeated / superseded | `Superseded` or older `Completed` with `repeat_stage` + newer attempt exists |

**Locks (out-of-sequence prevention)**

- Only **HR** can schedule a stage until prerequisites are met.
- Only **assignee** can record attendance / start / evaluate / absence actions.
- **Evaluation** and **final decisions** are immutable once saved.

---

## WebSocket events (optional live UI)

Payloads are JSON; same auth limitations as `/recruitment/ws/candidates` (browser WebSocket + `Authorization` header).

| Event | When |
|-------|------|
| `interview_workflow.process_started` | Process created |
| `interview_workflow.stage_scheduled` | New attempt row |
| `interview_workflow.attendance` | Attendance saved |
| `interview_workflow.evaluation` | Evaluation submitted |
| `interview_workflow.absence` | Absence fail or reschedule |
| `interview_workflow.final_approval` | Approver decided |

---

## UI implementation notes

1. **Timeline page:** `GET .../timeline` for audit dots; `GET .../processes/:id` for cards per attempt (`stageSequence` + `attemptNumber`).
2. **Interviewer panel:** Show actions only if `currentUser.employeeId === attempt.assignedEmployeeId` and attempt is in an actionable state.
3. **Progress bar:** `completedStages / snapshot.length` using latest non-superseded **`Completed`** + **`proceed`** per sequence.
4. **Approval dashboard:** Filter approvals where `assignedEmployeeId === me` and `status === pending` and prior steps approved (or query full process and compute client-side).

---

## Relationship to legacy interviews

| Feature | Legacy `/recruitment/interviews` | This workflow |
|---------|----------------------------------|----------------|
| Panel | Multiple `interviewerEmployeeIds` | **Single** assignee per attempt (exclusive) |
| HR between rounds | `hr-decision` on interview row | **Sequential stages** + optional **final** HR/Manager/CEO chain |
| Configurable stages | Fixed rounds 1–3 | **HR-defined** definition per job |

Choose per integration; data is **not** auto-synced between the two models.
