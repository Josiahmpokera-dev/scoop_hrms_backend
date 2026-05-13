# Recruitment Dashboard API

**Implemented:** `GET /api/v1/recruitment/dashboard/summary` (Bearer auth, permission `recruitment:read`). Optional query: `department_id` (numeric department PK, matched to requisition `department` name), `as_of` (RFC3339 or `YYYY-MM-DD` UTC).

This document specifies the **backend APIs** required to power the **Recruitment Dashboard** screen in the HRMS web app (`/recruitment/dashboard`). It is separate from the general recruitment module overview in `docs/api/RECRUITMENT_API.md` and focuses only on **aggregated metrics and lists** shown on that page.

---

## Purpose

The dashboard currently mixes **one live call** (`listOpenings` with `status=Active` for the first KPI) with **mock data** for everything else. Backend should implement the contract below so the frontend can replace mocks with a **single summary call** (recommended) or a small set of **read-only aggregate endpoints**.

**Goals:**

- Fast first paint (one round-trip preferred).
- Consistent numbers (same definitions as list/detail APIs for candidates, applications, interviews, offers, openings).
- Clear semantics for dates, stages, and “active” entities.

---

## Base path and auth

| Item | Value |
|------|--------|
| Base path | `/api/v1/recruitment` (align with existing recruitment routes). |
| Auth | `Authorization: Bearer <access_token>` |
| Permission | `recruitment:read` (or your module’s read scope for recruitment analytics). |

No public (unauthenticated) variant is required for this dashboard.

---

## Recommended: consolidated dashboard summary

### `GET /api/v1/recruitment/dashboard/summary`

Returns one JSON document that maps **directly** to the dashboard widgets: four KPI cards, pipeline bars, quick stats, top openings, upcoming interviews, and recent candidates.

**Query parameters (optional, v1 can ignore)**

| Query | Description |
|-------|-------------|
| `department_id` | Limit metrics to a department (future filter). |
| `as_of` | ISO date; default `now` (for reproducible reporting / tests). |

**Response:** `200 OK`

**Body (conceptual schema)** — field names may be `camelCase` or `snake_case` per your API standard; the frontend can adapt if documented.

```json
{
  "kpis": {
    "activeOpeningsCount": 12,
    "activeOpeningsHeadcountSum": 18,
    "totalCandidates": 240,
    "newCandidatesLast7Days": 14,
    "scheduledInterviewsCount": 8,
    "scheduledInterviewsToday": 2,
    "pendingOffersCount": 5,
    "offersExpiringWithinDays": 1,
    "offersExpiringWithinDaysThreshold": 3
  },
  "pipeline": {
    "applied": 80,
    "screening": 45,
    "assessment": 20,
    "interview": 35,
    "offer": 12,
    "hired": 8
  },
  "quickStats": {
    "avgTimeToFillDays": 32,
    "avgTimeToHireDays": 28,
    "offerAcceptanceRatePercent": 85,
    "pipelineConversionRatePercent": 12,
    "topSourceLabel": "LinkedIn",
    "topSourceSharePercent": 38,
    "referralSuccessRatePercent": 65
  },
  "topOpenings": [
    {
      "jobOpeningId": "uuid",
      "jobTitle": "Senior Engineer",
      "department": "Engineering",
      "location": "Dar es Salaam",
      "totalApplicants": 42
    }
  ],
  "upcomingInterviews": [
    {
      "interviewId": "uuid",
      "candidateId": "uuid",
      "candidateName": "Jane Doe",
      "jobOpeningId": "uuid",
      "jobTitle": "Senior Engineer",
      "interviewType": "Technical",
      "scheduledDate": "2026-05-15",
      "scheduledTime": "10:00",
      "status": "Scheduled"
    }
  ],
  "recentCandidates": [
    {
      "candidateId": "uuid",
      "fullName": "John Smith",
      "currentDesignation": "Software Developer",
      "location": "Nairobi",
      "primaryApplicationJobTitle": "Product Manager",
      "primaryApplicationAppliedAt": "2026-05-10T14:00:00Z",
      "stage": "Screening",
      "score": 78
    }
  ],
  "meta": {
    "generatedAt": "2026-05-13T10:00:00Z",
    "cacheTtlSeconds": 60
  }
}
```

### Widget mapping (acceptance criteria for backend)

| UI section | Response path | Rules |
|------------|---------------|--------|
| **Active Openings** (big number) | `kpis.activeOpeningsCount` | Count of job openings in **`Active`** status (include legacy **`Published`** if your data model still uses it). Same rule as `GET /openings?status=Active`. |
| **“X positions”** under openings | `kpis.activeOpeningsHeadcountSum` | Sum of **approved headcount / openings slots** per opening (whatever field represents “number of roles” on the opening; align with openings list UI). |
| **Total Candidates** | `kpis.totalCandidates` | Distinct candidates (or distinct person records) in recruitment scope; define whether archived / deleted are excluded. |
| **“Y this week”** | `kpis.newCandidatesLast7Days** | New candidates **or** new applications in the last 7×24 hours—**pick one** and document it; recommend **first application created** per candidate in window. |
| **Upcoming Interviews** (big number) | `kpis.scheduledInterviewsCount` | Interviews with status **`Scheduled`** and `scheduledDateTime` **≥ start of today** (or **> now**—document choice). |
| **“Z today”** | `kpis.scheduledInterviewsToday` | Subset of those with local calendar date = today. |
| **Pending Offers** | `kpis.pendingOffersCount` | Offers in **`Sent`**, **`Pending Approval`**, or equivalent “not accepted / not rejected / not expired” states—list allowed statuses in your OpenAPI. |
| **“Expiring soon”** | `kpis.offersExpiringWithinDays` | Count of pending offers with `expiryDate` within **`offersExpiringWithinDaysThreshold`** days (default **3**); echo threshold in response for UI copy. |
| **Pipeline** bars | `pipeline.*` | Integer counts per stage for **primary application** (or latest application)—must use the **same stage enum** as `PATCH …/applications/.../stage` elsewhere. Denominator for progress bars on the client is `max(totalCandidates, 1)` unless you return `pipelineDenominator`. |
| **Quick Stats** | `quickStats.*` | All **server-computed** from historical data (see definitions below). |
| **Top Performing Jobs** | `topOpenings` | Top **5** by `totalApplicants` among **Active** openings; tie-break by `totalApplicants` then `jobTitle`. |
| **Upcoming Interviews** list | `upcomingInterviews` | Next **5** by ascending schedule after now; only **Scheduled**. |
| **Recent Candidates** | `recentCandidates` | Last **5** by `primaryApplicationAppliedAt` descending. `score` optional (omit or `null` if not scored). |

---

## Stage model alignment

The dashboard UI today shows pipeline labels: **Applied**, **Screening**, **Assessment**, **Interview**, **Offer**, **Hired**.

Your applications API may use a slightly different set (e.g. **InterviewCompleted**, **TalentPool**, no **Assessment**). Backend should either:

1. **Return pipeline keys exactly as the UI expects** (with `assessment` possibly always `0` until the product adds that stage), or  
2. **Return a `pipeline` object with canonical backend keys** plus a **`pipelineDisplayMap`** for the frontend—avoid if possible; prefer (1) with stable enums documented in OpenAPI.

**Rejected / Withdrawn** candidates should **not** be included in pipeline funnel counts unless product explicitly wants “all ever” counts—default recommendation: **count only non-terminal applications** in funnel, and document that.

---

## Quick stats — suggested definitions

These are **analytic** fields; backend owns the formulas.

| Field | Suggested definition |
|-------|----------------------|
| `avgTimeToFillDays` | Mean days from **opening activated** to **first hire** on that opening (only closed/filled openings in window, e.g. last 12 months). |
| `avgTimeToHireDays` | Mean days from **candidate’s first application** to **Hired** for hired rows in same window. |
| `offerAcceptanceRatePercent` | `accepted_offers / (accepted_offers + declined_offers)` in window, ×100, rounded. |
| `pipelineConversionRatePercent` | e.g. `hired / applications_created` in window ×100—**document exact numerator/denominator**. |
| `topSourceLabel` + `topSourceSharePercent` | Dominant `application.source` (or candidate source) in last N days. |
| `referralSuccessRatePercent` | e.g. offers accepted where source = referral / total referral applications—**document**. |

If some metrics are not yet implementable, return `null` for that field and the UI can show “—” or hide the row.

---

## Alternative: multiple endpoints (not preferred)

If you do **not** ship a consolidated endpoint initially, the frontend can compose the dashboard from existing resources (with more latency and risk of inconsistent filters):

| Data | Existing / proposed endpoint |
|------|-------------------------------|
| Active openings count | `GET /openings?status=Active&page=1&limit=1` → `meta.total` (+ separate sum for headcount if not in list row). |
| Candidates / pipeline | `GET /candidates` with aggregations **not** in current list API → likely inefficient; **prefer SQL aggregates** behind `dashboard/summary`. |
| Interviews | `GET /interviews?status=Scheduled&…` |
| Offers | `GET /offers?status=…` |

Document pagination and filter parity so dashboard numbers match drill-down pages.

---

## Errors

| HTTP | When |
|------|------|
| `401` | Missing or invalid token. |
| `403` | Missing `recruitment:read` (or configured permission). |
| `500` | Aggregation failure; body should include `message` / `code` per your error standard. |

---

## Performance and caching

- Target **p95 &lt; 500 ms** for `GET …/dashboard/summary` with typical data volumes; use indexed queries and pre-aggregated tables if needed.
- Optional `meta.cacheTtlSeconds` hints short CDN / client cache; **do not** cache longer than business tolerates (stale interview counts annoy users).

---

## Versioning

- Bump or add `Accept-Version` / `?apiVersion=2` if you change pipeline key names or quick-stat definitions.
- Changelog entry when **Assessment** stage is added or merged on the backend.

---

## Frontend follow-up (out of scope for this doc)

After backend ships `GET /api/v1/recruitment/dashboard/summary`, the Next.js page `src/app/(protected-pages)/recruitment/dashboard/page.tsx` should:

1. Call the new endpoint (via `src/services/recruitment/recruitmentApi.ts` client).
2. Remove dependency on `@/mock/data/recruitment/recruitmentData` for those sections.
3. Keep `listOpenings` only if redundant with `kpis`—prefer **single** source from summary.

---

## Summary for backend ticket

**Implement:** `GET /api/v1/recruitment/dashboard/summary`  
**Auth:** Bearer + `recruitment:read`  
**Returns:** `kpis`, `pipeline`, `quickStats`, `topOpenings` (5), `upcomingInterviews` (5), `recentCandidates` (5), `meta`  
**Align:** Opening status (`Active`/`Published`), application stages, interview `Scheduled` semantics, offer pending + expiry rules with existing recruitment entities.
