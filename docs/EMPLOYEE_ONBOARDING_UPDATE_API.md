# Employee Onboarding — Update Step API

This document describes the endpoints to edit specific onboarding steps for an employee draft. These endpoints are idempotent and designed for stepper-based UIs where a user can revisit any step, modify fields, and save progress.

## Overview

- Update a specific step using the employee_id (recommended) or draft_id (legacy support).
- Steps supported: 1–10.
- Upsert semantics per step:
  - Steps 1,2,3,4,5,8,9,10: existing records are updated; if none exists, a new record is created.
  - Step 6 (Documents): replacing the list via JSON will delete and recreate draft documents; multipart upload replaces an individual document by type.
  - Step 7 (Assets): replacing the list removes prior draft assets and saves new ones.
- Progress is recalculated after each update.
- Required steps for completion: 1, 2, and 9.

## Authentication & Roles

- Requires a valid JWT.
- HR/Admin access enforced by middleware.

## Endpoints

### Update Step by Employee ID

- Method: PUT  
- Path: `/api/v1/employees/onboarding/:employee_id/step/:step`  
- Body (JSON):
```json
{
  "step": 1,
  "data": { }
}
```
- Notes:
  - `step` in body must match the `:step` in the path (1–10).
  - `:employee_id` is case-insensitive; it is normalized to uppercase internally.
  - Creates a draft automatically when updating step 1 and a draft doesn’t exist.

### Update Step by Draft ID (Legacy)

- Method: PUT  
- Path: `/api/v1/employees/onboarding/draft/:draft_id/step/:step`  
- Body (JSON):
```json
{
  "step": 2,
  "data": { }
}
```
- Notes:
  - `step` in body must match the `:step` in the path (1–10).

### Multipart Upload (Step 1 and Step 6)

For photo/documents use the existing upload endpoints:

- Photo (Step 1):  
  `POST /api/v1/employees/onboarding/upload-photo`  
  form-data: `employee_id`, `file`

- Documents (Step 6):  
  `POST /api/v1/employees/onboarding/upload-document`  
  form-data: `employee_id`, `document_type`, `file`, `description?`

These calls integrate with drafts and recalculate progress. A document of the same type replaces the previous file.

## Request Payloads Per Step (Data Field)

Below are representative fields expected in `data` for each step. Only relevant fields need to be sent. Fields omitted are left unchanged or cleared as per step semantics.

- Step 1 — Personal Info
  - `first_name`, `last_name` (required)
  - `middle_name?`, `date_of_birth?`, `gender?`, `marital_status?`, `blood_group?`, `nationality?`
  - `personal_email?`, `mobile_number?`, `alternate_number?`
  - `current_address?`, `permanent_address?`, `city?`, `state?`, `postal_code?`, `country?`
  - `employee_id?` (normalized and saved; auto-generated if missing on step 1)

- Step 2 — Employment
  - `official_email?`, `date_of_joining?`
  - `department_id?`, `position_id?`, `grade?`
  - `reporting_manager_id?` (auto-suggested from department head if not provided)
  - `employment_type?`, `location_id?`, `shift?`, `work_phone?`
  - `probation_period_days?`, `expected_confirmation_date?`
  - `employee_id?` (optional override; normalized)

- Step 3 — Salary
  - `annual_ctc?`, `ctc_effective_date?`, `currency?`
  - `basic_salary?`, `house_rent_allowance?`, `transport_allowance?`, `special_allowance?`, `other_allowances?`
  - `income_tax?`, `provident_fund?`, `professional_tax?`, `other_deductions?`
  - Gross/net computed automatically

- Step 4 — Bank
  - `bank_name?`, `account_holder_name?`, `account_number?`, `account_type?`, `branch_name?`, `swift_code?`

- Step 5 — Statutory
  - `tin_number?`, `nssf_number?`, `nhif_number?`, `wcf_number?`, `sdl_number?`
  - `passport_number?`, `passport_expiry_date?`, `work_permit_number?`, `work_permit_expiry_date?`

- Step 6 — Documents
  - JSON update (replacement):
    ```json
    {
      "documents": [
        {
          "document_type": "identity",
          "file_name": "id.pdf",
          "file_url": "https://.../id.pdf",
          "mime_type": "application/pdf",
          "file_size": 12345,
          "description": "National ID"
        }
      ]
    }
    ```
  - Multipart update: use `/upload-document` to replace a single type.

- Step 7 — Assets
  - `assets`: array of assets with fields `asset_type`, `asset_name`, `serial_number?`, `asset_tag?`, `assigned_date?`, `expected_return_date?`, `condition?`, `notes?`

- Step 8 — Policies
  - `leave_policy_id?`, `attendance_policy_id?`, `weekly_off_days?`, `effective_date?`

- Step 9 — Emergency Contacts
  - `contacts`: array of objects with `contact_name`, `relationship?`, `phone_number`, `alternate_phone?`, `email?`, `address?`, `is_primary`

- Step 10 — Notes
  - `notes?`

## Examples

### Update Step 1 (Personal Info)
PUT `/api/v1/employees/onboarding/EMP001/step/1`
```json
{
  "step": 1,
  "data": {
    "first_name": "Jane",
    "last_name": "Doe",
    "personal_email": "jane.doe@example.com",
    "mobile_number": "+255700000001",
    "current_address": "Mbezi Beach",
    "city": "Dar es Salaam",
    "country": "TZ"
  }
}
```

### Update Step 2 (Employment)
PUT `/api/v1/employees/onboarding/EMP001/step/2`
```json
{
  "step": 2,
  "data": {
    "official_email": "jane.doe@company.com",
    "department_id": 12,
    "position_id": 34,
    "reporting_manager_id": 56,
    "employment_type": "full_time",
    "date_of_joining": "2026-04-01T00:00:00Z",
    "location_id": 2
  }
}
```

### Update Step 6 (Documents) via JSON (replace list)
PUT `/api/v1/employees/onboarding/EMP001/step/6`
```json
{
  "step": 6,
  "data": {
    "documents": [
      {
        "document_type": "identity",
        "file_name": "passport.pdf",
        "file_url": "https://s3-region/bucket/.../passport.pdf",
        "mime_type": "application/pdf"
      }
    ]
  }
}
```

## Responses

Successful update:
```json
{
  "success": true,
  "message": "Step updated successfully",
  "data": {
    "draft_id": 123,
    "employee_id": "EMP001",
    "progress": 60,
    "completed_steps": [1,2,6],
    "finished_steps": [1,2,6,...],
    "unfinished_steps": [3,4,5,7,8,9,10],
    "is_completed": false,
    "steps": {
      "1": { },
      "2": { },
      "6": { }
    },
    "created_at": "2026-03-30T10:00:00Z",
    "updated_at": "2026-03-30T10:05:00Z"
  }
}
```

Common error cases:
- 400: invalid step, missing IDs, invalid references, or “draft not found for employee ID 'X'. Please save step 1 first”.
- 422: validation errors for invalid field formats.
- 404: draft/employee not found where applicable.

## Notes

- Employee ID is normalized to uppercase and trimmed.
- Step completion and overall progress recalculate on each update.
- Completion requires steps 1, 2, and 9 to be completed.

