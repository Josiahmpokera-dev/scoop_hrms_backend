# Employee Onboarding: Existing User vs New Person (Frontend Guide)

This document describes how to implement the **two onboarding flows** on the frontend:

1. **New person (non-user)** – At the end, the system creates login credentials and returns them.
2. **Existing user (already in system, not yet an employee)** – At the end, no credentials are created; the employee record is linked to their existing account.

---

## Summary

| Flow | Create draft | Step 2 (Employment) official_email | Complete onboarding response |
|------|----------------|------------------------------------|------------------------------|
| **New person** | `POST .../draft` (no `linked_user_id` / `existing_user_email`) | Any work email | `credentials_created: true`, `credentials: { email, password, username }` |
| **Existing user** | `POST .../draft` with `linked_user_id` or `existing_user_email` | Must match that user’s email (pre-filled from API) | `credentials_created: false`, `linked_to_existing_user: true`, no `credentials` |

---

## Existing user flow: end-to-end (with samples)

This section walks through **onboarding an existing user** (someone who already has a login but is not yet an employee) with endpoints and sample request/response bodies.

**Base URL:** `{{BASE_URL}}` (e.g. `https://api.example.com`)  
**Auth:** All requests require a valid Bearer token (HR/Admin).

---

### Step A: List users who are not yet employees

Get users you can onboard as employees. Use the returned `id` in Step B.

**Endpoint:** `GET {{BASE_URL}}/api/v1/employees/onboarding/non-employee-users`

**Query (optional):** `page=1`, `page_size=20`, `search=jane`

**Request example:**
```http
GET {{BASE_URL}}/api/v1/employees/onboarding/non-employee-users?page=1&page_size=20
Authorization: Bearer <token>
```

**Response example:**
```json
{
  "success": true,
  "message": "Non-employee users retrieved successfully",
  "data": [
    {
      "id": 42,
      "username": "jane",
      "email": "jane@company.com",
      "first_name": "Jane",
      "last_name": "Doe",
      "tenant_id": 1,
      "role": "user",
      "is_active": true,
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 1,
    "total_pages": 1
  }
}
```

**Take note:** Use `data[0].id` (e.g. `42`) and `data[0].email` (`jane@company.com`) for the next step.

---

### Step B: Create onboarding draft for that user

Create a draft that is linked to the existing user. **No credentials will be created** when you complete onboarding.

**Endpoint:** `POST {{BASE_URL}}/api/v1/employees/onboarding/draft`

**Request body (option 1 – by user ID):**
```json
{
  "linked_user_id": 42
}
```

**Request body (option 2 – by email):**
```json
{
  "existing_user_email": "jane@company.com"
}
```

**Request example:**
```http
POST {{BASE_URL}}/api/v1/employees/onboarding/draft
Authorization: Bearer <token>
Content-Type: application/json

{
  "linked_user_id": 42
}
```

**Response example:**
```json
{
  "success": true,
  "message": "Onboarding draft created successfully",
  "data": {
    "id": 5,
    "tenant_id": 1,
    "employee_id": null,
    "linked_user_id": 42,
    "completed_steps": "[]",
    "progress": 0,
    "is_completed": false,
    "created_at": "2025-01-20T09:00:00Z",
    "updated_at": "2025-01-20T09:00:00Z"
  }
}
```

**Take note:** You now have a `draft_id` (e.g. `5`). After Step 1 (personal info) is saved, you will also have an `employee_id` (e.g. `EMP001`) to use in step URLs.

---

### Step C: Get draft (to pre-fill Step 2 and show "existing user" message)

Load the draft to get `linked_user_email` for pre-filling **Official email** in Step 2 and to know that no credentials will be created.

**Endpoint:** `GET {{BASE_URL}}/api/v1/employees/onboarding/draft/5`  
(or `GET {{BASE_URL}}/api/v1/employees/onboarding/EMP001` once you have `employee_id`)

**Request example:**
```http
GET {{BASE_URL}}/api/v1/employees/onboarding/draft/5
Authorization: Bearer <token>
```

**Response example (relevant fields for existing user):**
```json
{
  "success": true,
  "message": "Draft retrieved successfully",
  "data": {
    "draft_id": 5,
    "employee_id": "EMP001",
    "linked_user_id": 42,
    "linked_user_email": "jane@company.com",
    "is_existing_user_onboarding": true,
    "progress": 20,
    "completed_steps": [1],
    "finished_steps": [1],
    "unfinished_steps": [2, 3, 4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "steps": {
      "1": { "first_name": "Jane", "last_name": "Doe" },
      "2": null
    },
    "created_at": "2025-01-20T09:00:00Z",
    "updated_at": "2025-01-20T09:05:00Z"
  }
}
```

**Frontend behaviour:**
- Pre-fill **Step 2 → Official email** with `linked_user_email` (`jane@company.com`) and keep it matching the existing account.
- Show: **"This employee will use their existing account. No new credentials will be created."**

---

### Step D: Save steps (1–10) as usual

Use the same save-step endpoints as for a new person. For **existing user**, Step 2 **official_email** must match the user's email (e.g. `jane@company.com`).

**Endpoint (by draft_id):** `POST {{BASE_URL}}/api/v1/employees/onboarding/draft/5/step/2`  
**Endpoint (by employee_id):** `POST {{BASE_URL}}/api/v1/employees/onboarding/EMP001/step/2`

**Request body example (Step 2 – Employment details):**
```json
{
  "step": 2,
  "data": {
    "official_email": "jane@company.com",
    "date_of_joining": "2025-01-20",
    "department_id": 1,
    "position_id": 2,
    "employment_type": "full_time"
  }
}
```

**Response:** Same structure as other step saves (draft with updated progress and completed steps). Repeat for steps 3–10 until the draft is complete.

---

### Step E: Complete onboarding (no credentials returned)

Finalise onboarding. Because this draft was created with `linked_user_id`, the backend **links the employee to the existing user** and **does not create or return credentials**.

**Endpoint (by draft_id):** `POST {{BASE_URL}}/api/v1/employees/onboarding/complete`  
**Endpoint (by employee_id):** `POST {{BASE_URL}}/api/v1/employees/onboarding/EMP001/complete`

**Request body (when using complete-by-draft_id):**
```json
{
  "draft_id": 5
}
```

**Request example:**
```http
POST {{BASE_URL}}/api/v1/employees/onboarding/complete
Authorization: Bearer <token>
Content-Type: application/json

{
  "draft_id": 5
}
```

**Response example (existing user – no credentials):**
```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "employee": {
      "id": 101,
      "employee_id": "EMP001",
      "first_name": "Jane",
      "last_name": "Doe",
      "user_id": 42,
      "tenant_id": 1,
      "status": "active",
      "created_at": "2025-01-20T10:00:00Z"
    },
    "credentials_created": false,
    "linked_to_existing_user": true
  }
}
```

**Important:** There is **no `credentials`** object. The user (e.g. Jane) continues to log in with their **existing** email and password.

**Frontend behaviour:**
- Check `data.credentials_created === false` and `data.linked_to_existing_user === true`.
- Do **not** show or store any credentials.
- Show: **"Employee record created and linked to their existing account. They can log in with their current credentials."**

---

### Process summary (existing user)

| Step | Action | Endpoint | Body / key response |
|------|--------|----------|---------------------|
| A | List non-employee users | `GET .../onboarding/non-employee-users` | Query: `page`, `page_size`, `search`. Response: `data[].id`, `data[].email` |
| B | Create draft for existing user | `POST .../onboarding/draft` | Body: `{ "linked_user_id": 42 }` or `{ "existing_user_email": "jane@company.com" }`. Response: `draft_id`, later `employee_id` |
| C | Get draft (pre-fill + hint) | `GET .../onboarding/draft/:draft_id` or `.../onboarding/:employee_id` | Response: `linked_user_email`, `is_existing_user_onboarding: true` |
| D | Save steps 1–10 | `POST .../onboarding/draft/:draft_id/step/:step` or `.../onboarding/:employee_id/step/:step` | Step 2 `official_email` must match `linked_user_email` |
| E | Complete onboarding | `POST .../onboarding/complete` or `.../onboarding/:employee_id/complete` | Body: `{ "draft_id": 5 }` or none. Response: `credentials_created: false`, `linked_to_existing_user: true`, **no** `credentials` |

---

## 1. Create draft

**Endpoint:** `POST {{BASE_URL}}/api/v1/employees/onboarding/draft`

**Auth:** Required (e.g. HR/Admin)

### New person (current behaviour)

Send an empty body or only step data. Do **not** send `linked_user_id` or `existing_user_email`.

```json
{}
```

Or with optional step 1 data:

```json
{
  "step": 1,
  "data": {
    "first_name": "Jane",
    "last_name": "Doe",
    ...
  }
}
```

### Existing user (onboard a user who is not yet an employee)

Send **either**:

- `linked_user_id` – ID of the user to onboard as employee, or  
- `existing_user_email` – email of that user (backend resolves to user and sets `linked_user_id`).

Example by user ID:

```json
{
  "linked_user_id": 42
}
```

Example by email:

```json
{
  "existing_user_email": "jane@company.com"
}
```

If the user is already an employee, the API returns **400** with a message like:  
`"user is already onboarded as an employee"`.

---

## 2. Get draft (for pre-fill and UI hints)

**Endpoints:**

- `GET {{BASE_URL}}/api/v1/employees/onboarding/draft/:draft_id`
- `GET {{BASE_URL}}/api/v1/employees/onboarding/:employee_id` (by employee_id string)

**Auth:** Required

Response includes:

| Field | Type | Meaning |
|-------|------|--------|
| `linked_user_id` | number \| null | Set when draft is for an existing user. |
| `linked_user_email` | string \| null | User’s email; use to **pre-fill Step 2 “Official email”** and to show a notice. |
| `is_existing_user_onboarding` | boolean | `true` = no credentials at end; show “Employee will use existing login”. |

**Frontend behaviour:**

- If `is_existing_user_onboarding` is `true`:
  - Pre-fill **Step 2 → Official email** with `linked_user_email` (and optionally make it read-only or show a note that it must match the existing account).
  - In the completion step or summary, show: **“This employee will use their existing account to log in. No new credentials will be created.”**
- If `is_existing_user_onboarding` is `false`:
  - Show the usual message that credentials will be generated at the end and must be shared securely.

---

## 3. Save steps (unchanged)

Same as today:

- `POST {{BASE_URL}}/api/v1/employees/onboarding/draft/:draft_id/step/:step`  
- or `POST {{BASE_URL}}/api/v1/employees/onboarding/:employee_id/step/:step`

For **existing user** onboarding, Step 2 **official_email** should match `linked_user_email` (the existing user’s email). The backend uses this to link the employee to the user at completion.

---

## 4. Complete onboarding

**Endpoints:**

- `POST {{BASE_URL}}/api/v1/employees/onboarding/complete`  
  Body: `{ "draft_id": 123 }`
- `POST {{BASE_URL}}/api/v1/employees/onboarding/:employee_id/complete`

**Auth:** Required

### Response shape (same for both endpoints)

```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "employee": { ... },
    "credentials_created": true,
    "linked_to_existing_user": false,
    "credentials": {
      "email": "newuser@company.com",
      "username": "newuser",
      "password": "generatedSecurePassword"
    }
  }
}
```

When the employee was **linked to an existing user** (no new account created):

```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "employee": { ... },
    "credentials_created": false,
    "linked_to_existing_user": true
  }
}
```

**Frontend behaviour:**

1. Read `data.credentials_created` (or `data.linked_to_existing_user`).
2. **If `credentials_created === true`:**
   - Show the “Credentials created” UI.
   - Display `data.credentials` (email, username, password) and advise the user to store/share them securely.
3. **If `linked_to_existing_user === true`:**
   - Do **not** show or store any credentials.
   - Show a message like: **“Employee record created and linked to their existing account. They can log in with their current credentials.”**

---

## 5. Optional: “Invite existing user to be onboarded”

To let HR choose “Onboard existing user” from a list:

1. **List users (and optionally mark who is already an employee)**  
   Use your existing users API; if you have an endpoint that returns “is_employee” or “employee_id”, use it to hide or label users who are already employees.

2. **Create draft with that user**  
   On “Start onboarding” for the selected user, call:

   `POST .../employees/onboarding/draft`  
   with `{ "linked_user_id": <user_id> }` or `{ "existing_user_email": "<email>" }`.

3. **Redirect to onboarding wizard**  
   Use the returned `draft_id` (or the `employee_id` after Step 1) and load draft with `GET .../draft/:draft_id` or `GET .../onboarding/:employee_id`. Use `linked_user_email` and `is_existing_user_onboarding` as above for pre-fill and messaging.

**Tip:** Use the dedicated API **`GET /api/v1/employees/onboarding/non-employee-users`** to list only users who are not yet employees (with optional `page`, `page_size`, `search`). Use each returned user's `id` as `linked_user_id` when creating the draft.

---

## 7. Quick reference

| Action | New person | Existing user |
|--------|------------|----------------|
| Create draft | No `linked_user_id` / `existing_user_email` | Send `linked_user_id` or `existing_user_email` |
| Step 2 official_email | Any work email | Pre-fill with `linked_user_email` (must match) |
| After complete | Show `credentials`, “Share securely” | Show “Use existing login”, no credentials |

Backend behaviour:

- **New person:** Creates employee + new user account, returns credentials.
- **Existing user:** Creates employee, links to existing user, assigns “employee” role if needed, returns no credentials. If official_email does not match an existing user, backend creates a new user and credentials as in the “new person” flow.
