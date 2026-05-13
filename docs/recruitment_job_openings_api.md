# Job openings API

Base path: `/api/v1/recruitment` (authenticated) and `/api/v1/recruitment/public` (no auth).

All job openings are created from an **approved** job requisition (`requisitionId`).

## Status lifecycle

| Status   | Meaning |
|----------|---------|
| `Draft`  | Editable; not visible on public list; not accepting applications. |
| `Active` | Live role; appears on public list (if not expired); accepts applications. Legacy rows may still show `Published` — they are treated the same as `Active` for listing and applying. |
| `Closed` | No longer accepting applications; not shown on public list. |

Flow: **`Draft`** → activate → **`Active`** → close → **`Closed`** → reopen → **`Draft`**.

---

## Configuration

| Environment variable | Purpose |
|---------------------|---------|
| `PUBLIC_RECRUITMENT_CAREERS_URL` | Optional. Front-end careers site origin without trailing slash, e.g. `https://careers.example.com`. Used in **share-link** responses as `careers_apply_url` (`…/apply?token=…`). If unset, only API URLs are returned. |

Proxies should send `X-Forwarded-Proto` and `X-Forwarded-Host` where applicable so generated links use the correct scheme and host.

---

## Authenticated APIs

Send header: `Authorization: Bearer <access_token>`.

Permissions use the existing recruitment codes (`recruitment:read`, `recruitment:create`, `recruitment:update`, `recruitment:publish`).

### Create job requisition (for hiring approval)

**POST** `/api/v1/recruitment/requisitions`

**Permission:** `recruitment:create`

**Body (JSON):**

```json
{
  "jobTitle": "Software Developer",
  "department": "Engineering",
  "location": "Dar es Salaam",
  "headcount": 2,
  "employmentType": "full_time",
  "targetStartDate": "2026-07-01",
  "justification": "We need to backfill attrition for the new product squad.",
  "priority": "High",
  "status": "Draft"
}
```

Notes on defaults:

- You can omit `hiringManager` — it defaults to **`HR`**.
- You can omit `grade`, `estimatedBudget`, and `currency` — they default to empty/zero (currency defaults to **`TZS`**).

**Response:** `200 OK`

```json
{
  "success": true,
  "message": "Requisition created successfully",
  "data": {
    "id": "9d20cfbc-0b38-49ea-924b-c13530b4b752",
    "message": "Requisition created successfully"
  }
}
```

### Create job opening

**POST** `/api/v1/recruitment/openings`

**Permission:** `recruitment:create`

**Body (JSON):**

```json
{
  "requisitionId": "<uuid-of-approved-requisition>",
  "jobTitle": "Software Developer",
  "jobDescription": "…",
  "responsibilities": ["…"],
  "mustHaveSkills": ["Go", "PostgreSQL"],
  "niceToHaveSkills": ["Docker"],
  "experience": { "min": 2, "max": 8 },
  "education": ["Bachelor's"],
  "expiryDate": "2026-12-31"
}
```

`expiryDate` is optional while **Draft**; supports `YYYY-MM-DD` or RFC3339.
`salaryRange` is no longer accepted on create; salary fields default internally (`salaryMin=0`, `salaryMax=0`, `salaryCurrency="TZS"`).

**Response:** `201 Created` — body includes generated `applyToken`.

---

### List job openings

**GET** `/api/v1/recruitment/openings?status=&department=&requisition_id=&page=1&limit=20`

**Permission:** `recruitment:read`

| Query | Description |
|-------|--------------|
| `status` | `Draft`, `Active`, `Closed`, or `Published` (legacy). Filtering by `Active` also includes legacy `Published` rows. |
| `department` | Filter via joined requisition. |
| `requisition_id` | Filter by requisition UUID. Alias: `requisitionId`. |
| `page`, `limit` | Pagination (default `page=1`, `limit=20`, max `100`). |

Returns top-level `data` + `meta` (no nested `data.data`):
`{ "success": true, "message": "...", "data": [...], "meta": { "page", "per_page", "total", "total_pages" } }`.

---

### List job openings grouped by status (UI dashboard)

**GET** `/api/v1/recruitment/openings/grouped?page=1&limit=20&department=&requisition_id=`

**Permission:** `recruitment:read`

This endpoint now returns a **flat list** (not grouped), same response format as `GET /openings` for UI consistency.

**Sample response:**

```json
{
  "success": true,
  "message": "Job openings retrieved successfully",
  "data": [
    {
      "id": "5a2e2c9b-0c26-4e0e-8cd4-3ff9b17b9d22",
      "requisitionId": "9d20cfbc-0b38-49ea-924b-c13530b4b752",
      "applyToken": "AbCdEf123456",
      "jobTitle": "Backend Engineer",
      "jobDescription": "Build internal APIs",
      "status": "Draft",
      "expiryDate": "2026-12-31T00:00:00Z",
      "salaryCurrency": "TZS",
      "requisition": {
        "id": "9d20cfbc-0b38-49ea-924b-c13530b4b752",
        "requisitionNumber": "REQ-2026-101",
        "department": "Engineering",
        "location": "Dar es Salaam",
        "hiringManager": "HR",
        "employmentType": "full_time",
        "headcount": 2,
        "priority": "High",
        "status": "Pending Approval"
      }
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

---

### Search job openings (same response shape as list-all)

**GET** `/api/v1/recruitment/openings/search?q=&status=&department=&requisition_id=&page=1&limit=20`

**Permission:** `recruitment:read`

Keyword `q` searches `jobTitle` and `jobDescription` (case-insensitive).  
Response shape is exactly the same as `GET /openings`: `{ data: [...], meta: {...} }`.

**Sample response:**

```json
{
  "success": true,
  "message": "Job openings retrieved successfully",
  "data": [
    { "id": "…", "jobTitle": "Senior Backend Engineer", "status": "Active" }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

Usage examples:

```bash
# Flat all-status list (legacy grouped route, now consistent format)
curl -sS "http://localhost:8080/api/v1/recruitment/openings/grouped?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"

# Search by keyword + status/department with same response shape
curl -sS "http://localhost:8080/api/v1/recruitment/openings/search?q=backend&status=Active&department=Engineering&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"
```

---

### Get one job opening

**GET** `/api/v1/recruitment/openings/:id`

**Permission:** `recruitment:read`

---

### Update job opening (Draft only)

**PUT** `/api/v1/recruitment/openings/:id`

**Permission:** `recruitment:update`

**Body:** Partial JSON (`UpdateJobOpeningRequest`) — any of: `jobTitle`, `jobDescription`, `responsibilities`, `mustHaveSkills`, `niceToHaveSkills`, `experienceMin`, `experienceMax`, `education`, `salaryMin`, `salaryMax`, `salaryCurrency`, `expiryDate`, `platforms`.

---

### Activate (publish)

**POST** `/api/v1/recruitment/openings/:id/publish`  
**POST** `/api/v1/recruitment/openings/:id/activate` _(alias)_

**Permission:** `recruitment:publish`

**Body:**

```json
{
  "platforms": ["careers_site", "linkedin"],
  "expiryDate": "2026-12-31"
}
```

- `platforms` optional (channels / labels for reporting).
- `expiryDate` **required unless** already set on the opening — determines last day the post accepts applications.

Sets status to **`Active`**.

---

### Close

**POST** `/api/v1/recruitment/openings/:id/close`

**Permission:** `recruitment:update`

Sets status to **`Closed`**.

---

### Reopen to Draft

**POST** `/api/v1/recruitment/openings/:id/reopen`

**Permission:** `recruitment:update`

Sets status from **`Closed`** back to **`Draft`** for edits.

---

### Share / apply links (for recruiters)

**GET** `/api/v1/recruitment/openings/:id/share-link`

**Permission:** `recruitment:read`

Returns JSON suitable for emailing or configuring the careers site:

```json
{
  "opening_id": "...",
  "apply_token": "AbCdEf123456",
  "careers_apply_url": "https://careers.example.com/apply?token=AbCdEf123456",
  "public_job_api_url": "https://api.example.com/api/v1/recruitment/public/openings/by-token/AbCdEf123456",
  "public_apply_url": "https://api.example.com/api/v1/recruitment/public/apply",
  "instructions": "..."
}
```

If `apply_token` was missing on old rows, the server generates one and persists it.

---

## Frontend: candidate apply form (which URL does what?)

Candidates **never** call authenticated `/api/v1/recruitment/...` routes. Their browser only talks to **`/api/v1/recruitment/public/...`** (no `Authorization` header).

There are **three** concepts recruiters and frontend developers mix up:

| Concept | What it is | Typical use |
|--------|------------|----------------|
| **`careers_apply_url`** | A **human-facing link** to **your** careers site | You implement a page at e.g. `/apply?token=…`. Requires **`PUBLIC_RECRUITMENT_CAREERS_URL`** in server config. |
| **`public_job_api_url`** | A **backend GET URL** (same API host as HRMS) | Your SPA calls this with **`fetch`** to load job title/description for the form. |
| **`public_apply_url`** | The **backend POST URL** where the form submits the application | Your SPA **`fetch(POST, JSON)`** here — this is what actually creates the application. |

So: the “public URL” for **accepting applications** is **`POST {API_ORIGIN}/api/v1/recruitment/public/apply`**. The pretty link (`careers_apply_url`) only routes the candidate to **your** UI; the UI then calls the API.

### End-to-end usage (recommended)

1. **HR creates an opening** and gets links (either **`GET …/openings/:id/share-link`** or **`POST …/openings?include_share_link=true`** after create). Response includes **`apply_token`**, **`public_job_api_url`**, **`public_apply_url`**, and optionally **`careers_apply_url`**.
2. **Publish** the opening (**activate**) so it is **`Active`** and not expired — otherwise public GET/apply return errors.
3. **Candidate opens your careers page**  
   - Either they use **`careers_apply_url`** (you host `https://careers.example.com/apply?token=AbCdEf…`), **or** you embed the token in any route you want (`/jobs/xyz?token=…`).  
   - Read **`token`** from the query string in your SPA/router.
4. **Your frontend loads job details** (to show title, JD, expiry):  
   **`GET`** `{API_ORIGIN}/api/v1/recruitment/public/openings/by-token/{token}`  
   Same URL as **`public_job_api_url`** from the share-link payload (token substituted).
5. **Your frontend submits the form** to:  
   **`POST`** `{API_ORIGIN}/api/v1/recruitment/public/apply`  
   Same as **`public_apply_url`** — **always this single endpoint** for every job; identify the role with **`applyToken`** (or **`jobId`**) in the JSON body.

### Frontend config

Point your UI at the API origin your users can reach (through your gateway):

```text
API_ORIGIN=https://hrms-api.example.com
```

- Job JSON: `GET ${API_ORIGIN}/api/v1/recruitment/public/openings/by-token/${token}`
- Submit: `POST ${API_ORIGIN}/api/v1/recruitment/public/apply`

If the careers site is on **another domain**, enable **CORS** on the API for that origin, or proxy **`/api/v1/recruitment/public/**`** through your careers domain.

### Minimal example (browser `fetch`)

Replace `API_ORIGIN` and read `token` from `new URLSearchParams(window.location.search).get('token')`.

```javascript
const API_ORIGIN = 'https://hrms-api.example.com';
const token = new URLSearchParams(window.location.search).get('token');
if (!token) { /* show error */ }

// 1) Load job for the form labels / JD
const jobRes = await fetch(
  `${API_ORIGIN}/api/v1/recruitment/public/openings/by-token/${encodeURIComponent(token)}`
);
const jobPayload = await jobRes.json();
const job = jobPayload.data ?? jobPayload; // adjust to your API envelope

// 2) User fills form; encode CV as data URL or base64 (see “Submit application” below)
const body = {
  applyToken: token,
  firstName: 'Jane',
  lastName: 'Doe',
  email: 'jane@example.com',
  phone: '+255700000000',
  coverLetter: '…',
  resume: dataUrlFromFileInput, // data:application/pdf;base64,...
  resumeName: 'cv.pdf',
};

const applyRes = await fetch(`${API_ORIGIN}/api/v1/recruitment/public/apply`, {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify(body),
});
```

### HR: get apply links when creating an opening

**POST** `/api/v1/recruitment/openings?include_share_link=true` — response includes **`share_link`** with the same fields as **`GET …/share-link`**, so you can show copy-paste URLs immediately after draft create.

---

## Public APIs (no auth)

### List active openings

**GET** `/api/v1/recruitment/public/openings?department=`

Returns openings with status **`Active`** or legacy **`Published`**, excluding expired posts.

---

### Job detail by apply token

**GET** `/api/v1/recruitment/public/openings/by-token/:token`

Returns the opening when it is active and not expired — for the careers page to render the form.

---

### Submit application

**POST** `/api/v1/recruitment/public/apply`  
**Content-Type:** `application/json`

Either **`jobId`** (opening UUID) or **`applyToken`** (from share link / opening record) must be provided.

```json
{
  "applyToken": "AbCdEf123456",
  "firstName": "Jane",
  "lastName": "Doe",
  "email": "jane@example.com",
  "phone": "+255700000000",
  "resume": "<base64 or https URL>",
  "resumeName": "cv.pdf",
  "linkedInUrl": "",
  "portfolioUrl": "",
  "coverLetter": "…"
}
```

**Response:** `201 Created` with application record.

---

### External careers application form (outside the HR portal)

Use this when you host a **separate** site or page (e.g. `careers.example.com`, a static landing page, or a minimal SPA) that is **not** the internal HR portal. The backend still exposes only **JSON** over HTTPS; there is **no** `multipart/form-data` upload endpoint — attachments are sent as **base64** (or an existing **HTTPS URL**) inside the JSON body.

#### End-to-end flow for the UI

1. **Discover the job**  
   - From the recruiter share payload, use **`careers_apply_url`** (see `PUBLIC_RECRUITMENT_CAREERS_URL`) or build your own route, e.g. `https://careers.example.com/apply?token=<applyToken>`.  
   - Read **`applyToken`** from the query string (or use **`jobId`** if you prefer the opening UUID from internal tools — candidates usually only see the token).

2. **Load job details for the form header / JD**  
   - **GET** `/api/v1/recruitment/public/openings/by-token/:token`  
   - Render title, description, and expiry so the candidate knows what they are applying for. If this returns **404**, the post is closed, expired, or inactive — show a friendly message.

3. **Submit the application**  
   - **POST** `/api/v1/recruitment/public/apply`  
   - **Content-Type:** `application/json`  
   - **No `Authorization` header** (public).

#### Form fields → JSON (`ApplyJobRequest`)

| UI control | JSON field | Required |
|------------|------------|----------|
| — | `applyToken` **or** `jobId` | One of them (token is typical for public pages) |
| First name | `firstName` | Yes |
| Last name | `lastName` | Yes |
| Email | `email` | Yes (valid email) |
| Phone | `phone` | No |
| Cover / motivation letter (textarea) | `coverLetter` | No — if empty, you may send **`applicationLetter`** instead (same storage) |
| LinkedIn | `linkedInUrl` | No |
| Portfolio | `portfolioUrl` | No |
| CV / résumé file | `resume` | No — but strongly recommended for recruitment |
| Original filename (helps storage naming) | `resumeName` | No — e.g. `Jane_Doe_CV.pdf` |

#### Attachment (CV) — how the API expects it

- **`resume` as HTTPS URL:** If the file is already hosted (CDN, temporary upload elsewhere), pass that URL string. The server stores it as-is on the candidate profile when it starts with `http`.  
- **`resume` as file upload from your form:** Use the browser **File** input → read the file as a **base64** string and send it in JSON:  
  - **Recommended:** send a full **data URL** so MIME is preserved, e.g. `data:application/pdf;base64,JVBERi0xLjQK...` — the API detects MIME from the prefix and picks an extension (pdf, doc, docx, jpg, png).  
  - **Alternative:** send **raw base64** only (no `data:` prefix); extension defaults to **`.pdf`** unless you set **`resumeName`** with a proper extension.  
- After upload, the backend stores the file via configured object storage and saves the resulting **`resumeUrl`** on the candidate — candidates do not POST multipart directly to MinIO/S3; everything goes through **`public/apply`**.

#### Minimal browser sketch (conceptual)

```text
On file selected:
  fileReader.readAsDataURL(file) → string → assign to payload.resume
  payload.resumeName = file.name
On submit:
  fetch(API_BASE + '/api/v1/recruitment/public/apply', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      applyToken: tokenFromQuery,
      firstName, lastName, email, phone,
      coverLetter,
      resume: dataUrlOrHttpsUrl,
      resumeName: originalFileName,
      linkedInUrl, portfolioUrl
    })
  })
```

Keep **`resume`** strings reasonably sized: very large payloads depend on your reverse proxy / API gateway body limits — enforce a **client-side max file size** (e.g. 5–10 MB) for CVs.

#### Cross-origin (careers site ≠ API host)

If the careers UI is on another origin than the API (e.g. UI on `careers.example.com`, API on `api.example.com`), the API must allow **CORS** for `POST` (and `OPTIONS` preflight) on `/api/v1/recruitment/public/apply`, plus **`GET`** for listing/detail if the browser calls those directly. Configure that on the API gateway or application server; otherwise host a **small backend-for-frontend** on the same origin as the form that proxies to the API.

#### Success / errors

- **201 Created:** application created; use the response body for confirmation UX (optional application id).  
- **4xx:** validation (missing names/email, invalid token/job), job not accepting applications, duplicate application rules — surface messages returned by the API.

---

## Example curls

Replace `TOKEN`, `OPENING_ID`, `REQ_ID`, `HOST`.

```bash
# Create (approved requisition)
curl -sS -X POST "http://localhost:8080/api/v1/recruitment/openings" \
  -H "Authorization: Bearer YOUR_JWT" \
  -H "Content-Type: application/json" \
  -d "{\"requisitionId\":\"$REQ_ID\",\"jobTitle\":\"Developer\",\"jobDescription\":\"We are hiring\",\"experience\":{\"min\":1,\"max\":5}}"

# Activate
curl -sS -X POST "http://localhost:8080/api/v1/recruitment/openings/$OPENING_ID/activate" \
  -H "Authorization: Bearer YOUR_JWT" \
  -H "Content-Type: application/json" \
  -d '{"expiryDate":"2026-12-31","platforms":["internal"]}'

# Share link
curl -sS "http://localhost:8080/api/v1/recruitment/openings/$OPENING_ID/share-link" \
  -H "Authorization: Bearer YOUR_JWT"

# Public list
curl -sS "http://localhost:8080/api/v1/recruitment/public/openings"

# Public detail
curl -sS "http://localhost:8080/api/v1/recruitment/public/openings/by-token/$TOKEN"

# Apply
curl -sS -X POST "http://localhost:8080/api/v1/recruitment/public/apply" \
  -H "Content-Type: application/json" \
  -d "{\"applyToken\":\"$TOKEN\",\"firstName\":\"Jane\",\"lastName\":\"Doe\",\"email\":\"jane@example.com\"}"
```

---

## Database

Run **`hrms-migrator`** after deploy so `job_openings` includes **`apply_token`** (and aligns with the `JobOpening` model). Existing **`Published`** rows continue to work; new activations set **`Active`**.

If `talent_pool_candidates` existed before new fields (`source_job_opening_id`, `source_application_id`, `internal_notes`, `last_contacted_at`), run a schema sync / migration so those columns exist.

---

## Applicants per job opening

**GET** `/api/v1/recruitment/openings/:id/applications?stage=&page=1&limit=20`

**Permission:** `recruitment:read`

Returns `{ data: [ JobApplication with candidate embedded ], meta }`. Filter by pipeline **`stage`** (see constants below).

---

## Candidates APIs (global list + filter + search + details)

### Unified candidates list (global + filters + search)

**GET** `/api/v1/recruitment/candidates?jobId=&stage=&q=&page=1&limit=20`

**Permission:** `recruitment:manage_candidates`

Use one endpoint for all scenarios:

- **Global list:** no filters
- **Filter by opening:** `jobId=<opening-uuid>`
- **Filter by stage:** `stage=Interview`
- **Search:** `q=jane` (matches first name, last name, email, phone)
- **Combined:** `jobId + stage + q + pagination`

**Sample response:**

```json
{
  "success": true,
  "message": "Candidates retrieved successfully",
  "data": [
    {
      "id": "cand-1",
      "firstName": "Jane",
      "lastName": "Doe",
      "email": "jane@example.com",
      "phone": "+255700000000",
      "resumeUrl": "https://cdn.example.com/recruitment/resumes/jane_cv.pdf",
      "linkedInUrl": "https://linkedin.com/in/jane",
      "portfolioUrl": "",
      "coverLetter": "…",
      "applications": [{ "id": "app-1", "jobOpeningId": "open-1", "stage": "Interview" }]
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

### Candidate details (profile + CV + applications + interviews + offers)

**GET** `/api/v1/recruitment/candidates/:id`

**Permission:** `recruitment:manage_candidates`

Returns full candidate profile and recruitment details:

- Profile fields (`firstName`, `lastName`, `email`, `phone`, `linkedInUrl`, `portfolioUrl`)
- **CV link** (`resumeUrl`)
- Cover letter
- Applications with linked job openings
- Interviews per application
- Offer (if exists)

**Sample response:**

```json
{
  "success": true,
  "message": "Candidate retrieved successfully",
  "data": {
    "id": "cand-1",
    "firstName": "Jane",
    "lastName": "Doe",
    "email": "jane@example.com",
    "phone": "+255700000000",
    "resumeUrl": "https://cdn.example.com/recruitment/resumes/jane_cv.pdf",
    "linkedInUrl": "https://linkedin.com/in/jane",
    "portfolioUrl": "https://jane.dev",
    "coverLetter": "…",
    "applications": [
      {
        "id": "app-1",
        "stage": "Interview",
        "jobOpening": { "id": "open-1", "jobTitle": "Backend Engineer", "status": "Active" },
        "interviews": [{ "id": "int-1", "round": 1, "status": "PanelCompleted" }],
        "offer": null
      }
    ]
  }
}
```

---

### Candidates by interview stage (round 1/2/3)

**GET** `/api/v1/recruitment/candidates/interview-stages?stage=1&jobId=&q=&page=1&limit=20`

**Permission:** `recruitment:manage_candidates`

Filter candidates selected/scheduled for interview by round:

- `stage=1` → Stage 1 interview candidates
- `stage=2` → Stage 2 interview candidates
- `stage=3` → Stage 3 interview candidates

You may also use `round` as an alias of `stage`.

**Sample response:**

```json
{
  "success": true,
  "message": "Candidates retrieved successfully",
  "data": [
    {
      "id": "cand-1",
      "firstName": "Jane",
      "lastName": "Doe",
      "email": "jane@example.com",
      "resumeUrl": "https://cdn.example.com/recruitment/resumes/jane_cv.pdf",
      "applications": [
        {
          "id": "app-1",
          "stage": "Interview",
          "jobOpening": { "id": "open-1", "jobTitle": "Backend Engineer" },
          "interviews": [{ "id": "int-1", "round": 1, "status": "Scheduled" }]
        }
      ]
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

Usage examples:

```bash
# Stage 1
curl -sS "http://localhost:8080/api/v1/recruitment/candidates/interview-stages?stage=1&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"

# Stage 2, one opening only, keyword search
curl -sS "http://localhost:8080/api/v1/recruitment/candidates/interview-stages?stage=2&jobId=$OPENING_ID&q=jane&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"
```

---

### Realtime updates via WebSocket (candidates + applications + interviews)

**GET** `/api/v1/recruitment/ws/candidates`  
**Permission:** `recruitment:read` (send `Authorization: Bearer <JWT>`)

Use this socket in the portal to auto-refresh list pages when new applications arrive or statuses change.

Common event types:

- `candidate.application.created` — new candidate application submitted
- `candidate.application.stage_updated` — stage changed (e.g. via PATCH stage API)
- `candidate.application.actioned` — accept/reject/move_to_interview_stage_one action applied
- `interview.scheduled` — interview created
- `interview.feedback_submitted` — panel feedback/marks submitted
- `interview.hr_decision` — HR decision + remarks captured

Message format:

```json
{
  "type": "candidate.application.actioned",
  "timestamp": "2026-05-13T12:20:00Z",
  "data": {
    "applicationId": "app-1",
    "candidateId": "cand-1",
    "jobId": "open-1",
    "action": "accept",
    "stage": "Screening"
  }
}
```

Frontend usage idea:

1. Connect to websocket after login.
2. On `candidate.*` events, refetch:
   - `GET /recruitment/candidates`
   - `GET /recruitment/openings/:id/applications` (if opening page is active)
3. On `interview.*` events, refetch:
   - `GET /recruitment/interviews`
   - `GET /recruitment/candidates/interview-stages?stage=1|2|3`

This keeps “all list APIs” consistent immediately after actions like Accept/Reject/Move to Stage 1/Interview updates.

---

## Pipeline stages (`JobApplication.stage`)

Use these values for filters and UI tabs:

| Stage | Typical meaning |
|-------|------------------|
| `Applied` | Submitted via public apply |
| `Screening` | HR reviewing CV |
| `Interview` | Interview scheduled / in progress |
| `InterviewCompleted` | Interviews done; decision pending |
| `Offer` | Offer extended |
| `Hired` | Accepted |
| `Rejected` | Not proceeding |
| `TalentPool` | Good profile but not advancing on this role; kept for outreach |

**PATCH** `/api/v1/recruitment/applications/:id/stage` — same body as `/candidates/:id/stage`: `{ "stage": "Screening", "notes": "…" }`  
**Permission:** `recruitment:manage_candidates`

### Application actions (accept / reject / move to interview stage 1)

**POST** `/api/v1/recruitment/applications/:id/action`  
**Permission:** `recruitment:manage_candidates`

Use this endpoint for explicit recruiter actions:

- `accept` → sets stage to `Screening`
- `reject` → sets stage to `Rejected`
- `move_to_interview_stage_one` → schedules round 1 interview and sets stage to `Interview`

#### 1) Accept application

```json
{
  "action": "accept",
  "notes": "Shortlisted after CV screening."
}
```

#### 2) Reject application

```json
{
  "action": "reject",
  "notes": "Profile does not match role requirements."
}
```

#### 3) Move to interview stage 1

```json
{
  "action": "move_to_interview_stage_one",
  "notes": "Proceed to first-round interview.",
  "interviewType": "Stage 1",
  "scheduledDate": "2026-06-01",
  "scheduledTime": "10:00",
  "duration": 60,
  "mode": "Video Call",
  "interviewerEmployeeIds": ["EMP001", "EMP042"]
}
```

Response:

```json
{
  "success": true,
  "message": "Application action processed successfully",
  "data": null
}
```

**POST** `/api/v1/recruitment/applications/:id/move-to-talent-pool` — body:

```json
{
  "notes": "Strong candidate; keep for future roles",
  "syncTalentPoolRow": true,
  "internalPoolNotes": "Prefers remote; ping Q3 hiring"
}
```

`syncTalentPoolRow` defaults to **true** when omitted (upserts `talent_pool_candidates` for email/contact). Set **`false`** to only change stage.

---

## Create opening + apply link in one response

**POST** `/api/v1/recruitment/openings?include_share_link=true`

Returns `{ "opening": { … }, "share_link": { … } }` so recruiters can publish the external link immediately after drafting requirements.

---

## Interviews

Interviews are organised in **up to three rounds** (`round`: **1**, **2**, or **3**) per candidate/job application. Each scheduled session has an **`attempt`** counter within that round (starts at **1**); HR can request a **same-round re-interview**, which schedules another session with `attempt` incremented and **`isReinterview`: true**.

### Interviewer panel (employees)

Prefer **`interviewerEmployeeIds`**: an array of **`employees.employee_id`** values (managers or any employee). The API resolves names/emails and stores a snapshot on **`interviewerPanel`**. You may still send **`interviewers`** as a legacy list of notification emails; at least one of **`interviewerEmployeeIds`** or **`interviewers`** is required.

### Status workflow

| Status | Meaning |
|--------|---------|
| `Scheduled` | Slot booked; panel has not submitted feedback |
| `PanelCompleted` | Interviewers submitted feedback — **HR must decide** before another interview is scheduled |
| `HrApprovedNextRound` | HR approved advancing — you may schedule the **next round** (or finish after round 3) |
| `HrRequestedReinterview` | HR wants another session **same round** — call schedule again with **`isReinterview`: true** |
| `HrRejected` | Candidate rejected at this stage |
| `Cancelled` | Reserved / optional cancellation flow |

Legacy rows may still show **`Completed`** instead of `PanelCompleted`; treat them like panel-complete for HR decisions.

### Schedule when the candidate already has an application

**POST** `/api/v1/recruitment/interviews` — **Permission:** `recruitment:manage_candidates`

Body **`ScheduleInterviewRequest`** (required unless noted):

| Field | Notes |
|-------|--------|
| `candidateId`, `jobId`, `interviewType` | Required |
| `round` | **1–3** (required) |
| `scheduledDate`, `scheduledTime` | Required |
| `interviewerEmployeeIds` | Preferred: employee IDs from `/employees` |
| `interviewers` | Optional legacy email list for notifications |
| `isReinterview` | **`true`** only after HR chose **`schedule_reinterview`** for the latest session of that round |

Rules:

- You cannot schedule **any** new interview while another row for the same application is **`Scheduled`** or **`PanelCompleted`** (waiting on feedback or HR).
- **Round 2** requires round **1**’s latest session to be **`HrApprovedNextRound`**. Same for **round 3** vs round **2**.
- **Re-interview** same round: latest session for that round must be **`HrRequestedReinterview`**, then **`isReinterview`: true** on the next **`POST /interviews`**.

Example (round 1, employee panel):

```json
{
  "candidateId": "<uuid>",
  "jobId": "<opening-uuid>",
  "interviewType": "Technical",
  "round": 1,
  "scheduledDate": "2026-06-01",
  "scheduledTime": "10:00",
  "duration": 60,
  "mode": "Video Call",
  "interviewerEmployeeIds": ["EMP001", "EMP042"]
}
```

### Panel feedback

**POST** `/api/v1/recruitment/interviews/:id/feedback` — existing **`SubmitFeedbackRequest`**. Sets status to **`PanelCompleted`** (not final hire/reject — HR still decides).

### HR decision (remarks required)

**POST** `/api/v1/recruitment/interviews/:id/hr-decision` — **Permission:** `recruitment:manage_candidates`

Requires authenticated user linked to an **employee** record (for audit), or pass **`decidedByEmployeeId`** explicitly.

```json
{
  "decision": "proceed_next_round",
  "remarks": "Strong technical fit; proceed to hiring manager round.",
  "decidedByEmployeeId": "EMP-HR-01"
}
```

**`decision`** (required):

- **`proceed_next_round`** — unlocks the next **`round`** schedule. If current **`round`** is **3**, sets application stage to **`InterviewCompleted`**.
- **`schedule_reinterview`** — unlocks another **`POST /interviews`** for the **same** **`round`** with **`isReinterview`: true**.
- **`reject_candidate`** — sets interview to **`HrRejected`** and application stage to **`Rejected`**.

**`remarks`** is required for every decision.

### Manual (walk-in / referral — no application yet)

**POST** `/api/v1/recruitment/interviews/schedule-manual`

Creates **Candidate** + **JobApplication** (stage `Interview`) + **Interview**. Same scheduling fields as above, including **`round`**, **`interviewerEmployeeIds`**, **`isReinterview`**.

### Batch from pipeline (many applicants, same slot)

**POST** `/api/v1/recruitment/interviews/schedule-batch`

```json
{
  "jobOpeningId": "<opening-uuid>",
  "applicationIds": ["app-1", "app-2"],
  "interviewType": "Technical",
  "round": 1,
  "scheduledDate": "2026-06-01",
  "scheduledTime": "10:00",
  "duration": 60,
  "mode": "Video Call",
  "interviewerEmployeeIds": ["EMP001"],
  "interviewers": ["hr@company.com"],
  "isReinterview": false,
  "updateStage": true
}
```

`updateStage` omitted or **true** sets each application to **`Interview`**. Set **`false`** to only create interview rows.

### Global interview list (all stages by default)

**GET** `/api/v1/recruitment/interviews?page=1&limit=20`

**Permission:** `recruitment:read`

By default this lists interviews across **all rounds/stages**.

Optional filters:

- `stage` (or `round`) = `1`, `2`, `3`
- `status` (e.g. `Scheduled`, `PanelCompleted`, `HrApprovedNextRound`)
- `jobId`
- `candidateId`

This response includes:

- interviewer info (`interviewerPanel`, `interviewers`)
- panel/HR remarks (`feedback.comments`, `hrRemarks`)
- marks/score (`feedback.rating`, 1-5)

**Sample response:**

```json
{
  "success": true,
  "message": "Interviews retrieved successfully",
  "data": [
    {
      "id": "int-1",
      "candidateId": "cand-1",
      "jobId": "open-1",
      "round": 2,
      "status": "PanelCompleted",
      "interviewerPanel": [
        { "employeeId": "EMP001", "displayName": "John HR", "workEmail": "john@company.com" }
      ],
      "interviewers": ["john@company.com"],
      "hrRemarks": "Good communication and technical depth",
      "feedback": {
        "rating": 4,
        "recommendation": "Hire",
        "comments": "Strong backend fundamentals"
      }
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 1, "total_pages": 1 }
}
```

Examples:

```bash
# All interview stages (default)
curl -sS "http://localhost:8080/api/v1/recruitment/interviews?page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"

# Stage 1 only
curl -sS "http://localhost:8080/api/v1/recruitment/interviews?stage=1&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"

# Stage 3 only for one opening
curl -sS "http://localhost:8080/api/v1/recruitment/interviews?stage=3&jobId=$OPENING_ID&page=1&limit=20" \
  -H "Authorization: Bearer YOUR_JWT"
```

### Interview details

**GET** `/api/v1/recruitment/interviews/:id` — **Permission:** `recruitment:read`

Returns one interview with full panel, remarks, and marks details.

### List interviews for an opening

**GET** `/api/v1/recruitment/openings/:id/interviews` — **Permission:** `recruitment:read` — includes **`feedback`** when present.

---

## Offers (salary & package)

Existing APIs unchanged:

- **POST** `/api/v1/recruitment/offers` — `CreateOfferRequest` (candidateId, jobId, joiningDate, compensation, benefits, …).
- **POST** `/api/v1/recruitment/offers/:id/approval`, **POST** `/api/v1/recruitment/offers/:id/send`

**GET** `/api/v1/recruitment/openings/:id/offers` — list offers for this role (**Permission:** `recruitment:read`).

---

## Talent pool (non‑shortlisted / nurture list)

**GET** `/api/v1/recruitment/talent-pool?page=1&limit=20` — paginated directory (**Permission:** `recruitment:manage_candidates`).  
**GET** `/api/v1/recruitment/talent-pool/search?skills=&location=&experienceMin=` — existing search.

**POST** `/api/v1/recruitment/talent-pool/:id/contact` — send email (SMTP via `.env`):

```json
{ "subject": "Future opportunities", "message": "<p>Hello …</p>" }
```

Updates **`lastContactedAt`** on success.

---

## Suggested UI screens (for your separate frontend)

Map screens to the endpoints above:

1. **Public apply (standalone careers site)** — page at **`/apply?token=`** (or your route): **GET** `public/openings/by-token/:token` for JD; **POST** `public/apply` with JSON + base64 CV as in **External careers application form** above.
2. **Requisition approved → Job opening wizard** — fields map to `CreateJobOpeningRequest`; call create with `?include_share_link=true`; show **share_link** on success.
3. **Opening detail** — tabs: **Applicants** (`GET …/applications`), **Interviews** (`GET …/interviews`), **Offers** (`GET …/offers`).
4. **Applicant pipeline** — table bound to `applications`; actions call **PATCH …/applications/:id/stage** or **move-to-talent-pool**.
5. **Schedule interviews** — “From applicants” multi-select → **schedule-batch**; “Manual entry” form → **schedule-manual**.
6. **Offer** — form → **POST /offers**; list via opening offers GET.
7. **Talent pool** — **GET /talent-pool** + **contact** action → **POST …/contact**.

This repository contains **APIs only**; wire these routes in your UI client (React/Vue/etc.).
