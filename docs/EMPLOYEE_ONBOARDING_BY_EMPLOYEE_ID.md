# Employee Onboarding API - Employee ID Based Tracking

## Overview

The Employee Onboarding API now tracks onboarding progress by **Employee ID** instead of Draft ID. This makes it much easier and more intuitive to manage employee onboarding.

### Key Features

✅ **Employee ID Auto-Generated**: Automatically generated when saving Step 1 (Personal Information)  
✅ **Employee ID Based Tracking**: Use employee ID for all operations  
✅ **Progress Tracking**: Automatic progress calculation (0-100%)  
✅ **Finished/Unfinished Steps**: Clear visibility of which steps are done  
✅ **Easy Resume**: Simply use employee ID to continue onboarding  

---

## How It Works

### Flow Diagram

```
┌─────────────────────────────────────────────────────────┐
│         Employee Onboarding Flow (Employee ID)          │
├─────────────────────────────────────────────────────────┤
│                                                         │
│  1. Create Draft → Get draft_id (temporary)             │
│  2. Save Step 1 → Employee ID auto-generated (EMP001)   │
│  3. Use Employee ID for all subsequent steps            │
│  4. View Progress → GET /onboarding/{employee_id}       │
│  5. Save Steps → POST /onboarding/{employee_id}/step/X  │
│  6. Complete → POST /onboarding/{employee_id}/complete  │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Employee ID Generation

- **Format**: `EMP` + 6-digit number (e.g., `EMP000001`, `EMP000002`)
- **When Generated**: Automatically when saving Step 1 (Personal Information)
- **Uniqueness**: System ensures each employee ID is unique
- **Persistence**: Once generated, it's stored in the draft and used for all operations

---

## API Endpoints

### 1. Create Draft (Step 1)

**Endpoint:** `POST /api/v1/employees/onboarding/draft`

**Description:** Creates a new onboarding draft. When you save step 1 data, the employee ID is automatically generated.

**Request Body:**

```json
{
  "step": 1,
  "data": {
    "first_name": "John",
    "last_name": "Doe",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "gender": "male",
    "personal_email": "john.doe@example.com",
    "mobile_number": "+255712345678"
  }
}
```

**Response:**

```json
{
  "success": true,
  "message": "Onboarding draft created and step 1 saved successfully. Employee ID: EMP000001",
  "data": {
    "draft_id": 1,
    "employee_id": "EMP000001",
    "progress": 10.0,
    "completed_steps": [1],
    "finished_steps": [1],
    "unfinished_steps": [2, 3, 4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "steps": {
      "1": {
        "first_name": "John",
        "last_name": "Doe",
        ...
      }
    }
  }
}
```

**Important:** Save the `employee_id` from the response! You'll use it for all subsequent operations.

---

### 2. Get Onboarding Status by Employee ID

**Endpoint:** `GET /api/v1/employees/onboarding/{employee_id}`

**Description:** Get the complete onboarding status, progress, and all saved step data for an employee.

**Path Parameters:**
- `employee_id` (required): The employee ID (e.g., `EMP000001`)

**Response:**

```json
{
  "success": true,
  "message": "Draft retrieved successfully",
  "data": {
    "draft_id": 1,
    "employee_id": "EMP000001",
    "progress": 30.0,
    "completed_steps": [1, 2, 3],
    "finished_steps": [1, 2, 3],
    "unfinished_steps": [4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "steps": {
      "1": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "date_of_birth": "1990-01-15T00:00:00Z",
        "gender": "male",
        "personal_email": "john.doe@example.com",
        "mobile_number": "+255712345678",
        "addresses": [...]
      },
      "2": {
        "id": 1,
        "employee_id": "EMP000001",
        "official_email": "john.doe@company.com",
        "department_id": 5,
        "position_id": 10,
        "employment_type": "full_time"
      },
      "3": {
        "id": 1,
        "annual_ctc": 24000000.00,
        "basic_salary": 15000000.00,
        "currency": "TZS"
      }
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**Response Fields Explained:**

- `employee_id`: The employee ID (use this for all operations)
- `progress`: Completion percentage (0-100%)
- `completed_steps`: Array of step numbers that are completed
- `finished_steps`: Array of steps that have data saved (same as completed_steps)
- `unfinished_steps`: Array of steps that don't have data yet
- `steps`: Object containing all saved step data, keyed by step number

**cURL Example:**

```bash
curl -X GET http://localhost:8080/api/v1/employees/onboarding/EMP000001 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

### 3. Save Step by Employee ID

**Endpoint:** `POST /api/v1/employees/onboarding/{employee_id}/step/{step}`

**Description:** Save data for a specific onboarding step using employee ID.

**Path Parameters:**
- `employee_id` (required): The employee ID (e.g., `EMP000001`)
- `step` (required): Step number (1-10)

**Request Body:**

```json
{
  "step": 2,
  "data": {
    "official_email": "john.doe@company.com",
    "date_of_joining": "2024-01-15T00:00:00Z",
    "department_id": 5,
    "position_id": 10,
    "employment_type": "full_time",
    "location_id": 3
  }
}
```

**Response:**

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 20.0,
    "completed_steps": [1, 2],
    "is_completed": false,
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 2,
    "data": {
      "official_email": "john.doe@company.com",
      "department_id": 5,
      "employment_type": "full_time"
    }
  }'
```

---

### 4. Complete Onboarding by Employee ID

**Endpoint:** `POST /api/v1/employees/onboarding/{employee_id}/complete`

**Description:** Complete the onboarding process and create the final employee record.

**Path Parameters:**
- `employee_id` (required): The employee ID (e.g., `EMP000001`)

**Request Body:** None (employee_id in path is sufficient)

**Response:**

```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "id": 100,
    "employee_id": "EMP000001",
    "first_name": "John",
    "last_name": "Doe",
    "work_email": "john.doe@company.com",
    "department_id": 5,
    "position_id": 10,
    "status": "active",
    "created_at": "2024-01-15T12:00:00Z"
  }
}
```

**cURL Example:**

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/complete \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**Requirements:**
- Steps 1 and 2 must be completed (minimum)
- All step data will be migrated to the final employee record

---




**Save Step 1:**
```
POST /api/v1/employees/onboarding/draft
```

**Or if employee_id already exists:**
```
POST /api/v1/employees/onboarding/{employee_id}/step/1
```

### Request Body

```json
{
  "step": 1,
  "data": {
    "photo_url": "https://example.com/photos/john-doe.jpg",
    "first_name": "John",
    "middle_name": "Michael",
    "last_name": "Doe",
    "date_of_birth": "1990-01-15T00:00:00Z",
    "gender": "male",
    "marital_status": "single",
    "blood_group": "O+",
    "nationality": "Tanzanian",
    "personal_email": "john.doe@example.com",
    "mobile_number": "+255712345678",
    "alternate_number": "+255712345679",
    "current_address": "123 Main Street, Kinondoni",
    "permanent_address": "456 Home Street, Arusha",
    "city": "Dar es Salaam",
    "state": "Dar es Salaam",
    "postal_code": "11101",
    "country": "Tanzania"
  }
}
```

### Required Fields

- `first_name` (required): Minimum 2 characters
- `last_name` (required): Minimum 2 characters

### Optional Fields

- `photo_url`: URL to employee photo
- `middle_name`: Middle name
- `date_of_birth`: ISO 8601 date format
- `gender`: "male", "female", or "other"
- `marital_status`: "single", "married", "divorced", or "widowed"
- `blood_group`: "A+", "A-", "B+", "B-", "AB+", "AB-", "O+", or "O-"
- `nationality`: Country name
- `personal_email`: Valid email address
- `mobile_number`: Phone number
- `alternate_number`: Alternate phone number
- `current_address`: Current residential address
- `permanent_address`: Permanent residential address
- `city`: City name
- `state`: State/region name
- `postal_code`: Postal/ZIP code
- `country`: Country name

### Response

```json
{
  "success": true,
  "message": "Onboarding draft created and step 1 saved successfully. Employee ID: EMP000001",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 10.0,
    "completed_steps": [1],
    "finished_steps": [1],
    "unfinished_steps": [2, 3, 4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T10:30:00Z"
  }
}
```

### Get Step 1 Data

**Endpoint:**
```
GET /api/v1/employees/onboarding/EMP000001
```

**Response:**
```json
{
  "success": true,
  "message": "Draft retrieved successfully",
  "data": {
    "draft_id": 1,
    "employee_id": "EMP000001",
    "progress": 10.0,
    "finished_steps": [1],
    "unfinished_steps": [2, 3, 4, 5, 6, 7, 8, 9, 10],
    "steps": {
      "1": {
        "id": 1,
        "photo_url": "https://example.com/photos/john-doe.jpg",
        "first_name": "John",
        "middle_name": "Michael",
        "last_name": "Doe",
        "date_of_birth": "1990-01-15T00:00:00Z",
        "gender": "male",
        "marital_status": "single",
        "blood_group": "O+",
        "nationality": "Tanzanian",
        "personal_email": "john.doe@example.com",
        "mobile_number": "+255712345678",
        "alternate_number": "+255712345679",
        "addresses": [
          {
            "id": 1,
            "address_type": "current",
            "address_line1": "123 Main Street, Kinondoni",
            "city": "Dar es Salaam",
            "state": "Dar es Salaam",
            "postal_code": "11101",
            "country": "Tanzania"
          },
          {
            "id": 2,
            "address_type": "permanent",
            "address_line1": "456 Home Street",
            "city": "Arusha",
            "country": "Tanzania"
          }
        ]
      }
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/draft \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 1,
    "data": {
      "first_name": "John",
      "last_name": "Doe",
      "date_of_birth": "1990-01-15T00:00:00Z",
      "gender": "male",
      "personal_email": "john.doe@example.com",
      "mobile_number": "+255712345678"
    }
  }'
```

**Important Notes:**
- Employee ID (`EMP000001`) is **auto-generated** and returned in the response
- Save the `employee_id` for all subsequent operations
- Addresses are automatically saved as separate records (current and permanent)

---

## Step 2: Employment Details

**Purpose:** Set employment information including department, position, and work location.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/2
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/2
```

### Request Body

```json
{
  "step": 2,
  "data": {
    "official_email": "john.doe@company.com",
    "date_of_joining": "2024-01-15T00:00:00Z",
    "department_id": 5,
    "position_id": 10,
    "grade": "G5",
    "reporting_manager_id": 20,
    "employment_type": "full_time",
    "location_id": 3,
    "shift": "Day",
    "work_phone": "+255222345678",
    "probation_period_days": 90,
    "expected_confirmation_date": "2024-04-15T00:00:00Z"
  }
}
```

### Required Fields

- None (all fields are optional, but `department_id` is recommended)

### Optional Fields

- `official_email`: Work email address (must be valid email format)
- `date_of_joining`: ISO 8601 date format
- `department_id`: Department ID (recommended - automatically sets organization_id)
- `position_id`: Job position ID (optional)
- `grade`: Employee grade (e.g., "G5", "G6")
- `reporting_manager_id`: Manager's employee ID
- `employment_type`: "full_time", "part_time", "contract", or "intern"
- `location_id`: Office location ID
- `shift`: "Day", "Night", "Rotating", etc.
- `work_phone`: Work phone number
- `probation_period_days`: Minimum 90 days
- `expected_confirmation_date`: ISO 8601 date format

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 20.0,
    "completed_steps": [1, 2],
    "finished_steps": [1, 2],
    "unfinished_steps": [3, 4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T11:00:00Z"
  }
}
```

### Get Step 2 Data

**Endpoint:**
```
GET /api/v1/employees/onboarding/EMP000001
```

**Response includes step 2:**
```json
{
  "steps": {
    "2": {
      "id": 1,
      "employee_id": "EMP000001",
      "official_email": "john.doe@company.com",
      "date_of_joining": "2024-01-15T00:00:00Z",
      "department_id": 5,
      "position_id": 10,
      "grade": "G5",
      "reporting_manager_id": 20,
      "employment_type": "full_time",
      "location_id": 3,
      "shift": "Day",
      "work_phone": "+255222345678",
      "probation_period_days": 90,
      "expected_confirmation_date": "2024-04-15T00:00:00Z"
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/2 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 2,
    "data": {
      "official_email": "john.doe@company.com",
      "department_id": 5,
      "position_id": 10,
      "employment_type": "full_time",
      "location_id": 3
    }
  }'
```

**Important Notes:**
- `employee_id` is NOT required in the request - it's already set from step 1
- `department_id` automatically sets `organization_id` and `organization_unit_id`
- `position_id` is optional

---

## Step 3: Salary & CTC (Cost to Company)

**Purpose:** Define salary structure, allowances, and deductions.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/3
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/3
```

### Request Body

```json
{
  "step": 3,
  "data": {
    "annual_ctc": 24000000.00,
    "ctc_effective_date": "2024-01-15T00:00:00Z",
    "currency": "TZS",
    "basic_salary": 15000000.00,
    "house_rent_allowance": 3000000.00,
    "transport_allowance": 2000000.00,
    "special_allowance": 2000000.00,
    "other_allowances": 2000000.00,
    "income_tax": 2400000.00,
    "provident_fund": 1200000.00,
    "professional_tax": 50000.00,
    "other_deductions": 0.00
  }
}
```

### Required Fields

- None (all fields are optional)

### Optional Fields

- `annual_ctc`: Annual Cost to Company (decimal)
- `ctc_effective_date`: ISO 8601 date format
- `currency`: Currency code (default: "TZS", must be 3 characters)
- `basic_salary`: Basic salary amount (decimal)
- `house_rent_allowance`: HRA amount (decimal)
- `transport_allowance`: Transport allowance (decimal)
- `special_allowance`: Special allowance (decimal)
- `other_allowances`: Other allowances (decimal)
- `income_tax`: Income tax deduction (decimal)
- `provident_fund`: PF deduction (decimal)
- `professional_tax`: Professional tax (decimal)
- `other_deductions`: Other deductions (decimal)

**Note:** Gross salary and net monthly salary are calculated automatically.

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 30.0,
    "completed_steps": [1, 2, 3],
    "finished_steps": [1, 2, 3],
    "unfinished_steps": [4, 5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T11:15:00Z"
  }
}
```

### Get Step 3 Data

**Response includes step 3:**
```json
{
  "steps": {
    "3": {
      "id": 1,
      "annual_ctc": 24000000.00,
      "ctc_effective_date": "2024-01-15T00:00:00Z",
      "currency": "TZS",
      "basic_salary": 15000000.00,
      "house_rent_allowance": 3000000.00,
      "transport_allowance": 2000000.00,
      "special_allowance": 2000000.00,
      "other_allowances": 2000000.00,
      "gross_salary": 24000000.00,
      "income_tax": 2400000.00,
      "provident_fund": 1200000.00,
      "professional_tax": 50000.00,
      "other_deductions": 0.00,
      "net_monthly_salary": 20350000.00
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/3 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 3,
    "data": {
      "annual_ctc": 24000000.00,
      "currency": "TZS",
      "basic_salary": 15000000.00,
      "house_rent_allowance": 3000000.00,
      "transport_allowance": 2000000.00
    }
  }'
```

---

## Step 4: Bank Account Details

**Purpose:** Store bank account information for salary payments.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/4
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/4
```

### Request Body

```json
{
  "step": 4,
  "data": {
    "bank_name": "CRDB Bank",
    "account_holder_name": "John Doe",
    "account_number": "1234567890",
    "account_type": "savings",
    "branch_name": "Kinondoni Branch",
    "swift_code": "CRDBTZTZ"
  }
}
```

### Required Fields

- None (all fields are optional)

### Optional Fields

- `bank_name`: Name of the bank
- `account_holder_name`: Account holder's name
- `account_number`: Bank account number
- `account_type`: "savings" or "current"
- `branch_name`: Bank branch name
- `swift_code`: SWIFT/BIC code

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 40.0,
    "completed_steps": [1, 2, 3, 4],
    "finished_steps": [1, 2, 3, 4],
    "unfinished_steps": [5, 6, 7, 8, 9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T11:30:00Z"
  }
}
```

### Get Step 4 Data

**Response includes step 4:**
```json
{
  "steps": {
    "4": {
      "id": 1,
      "bank_name": "CRDB Bank",
      "account_holder_name": "John Doe",
      "account_number": "1234567890",
      "account_type": "savings",
      "branch_name": "Kinondoni Branch",
      "swift_code": "CRDBTZTZ",
      "is_primary": true
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/4 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 4,
    "data": {
      "bank_name": "CRDB Bank",
      "account_holder_name": "John Doe",
      "account_number": "1234567890",
      "account_type": "savings",
      "branch_name": "Kinondoni Branch"
    }
  }'
```

---

## Step 5: Statutory Requirements

**Purpose:** Collect statutory compliance information (TIN, NSSF, NHIF, etc.).

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/5
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/5
```

### Request Body

```json
{
  "step": 5,
  "data": {
    "tin_number": "123456789",
    "nssf_number": "987654321",
    "nhif_number": "456789123",
    "wcf_number": "789123456",
    "sdl_number": "321654987",
    "passport_number": "A12345678",
    "passport_expiry_date": "2030-01-15T00:00:00Z",
    "work_permit_number": "WP123456",
    "work_permit_expiry_date": "2025-01-15T00:00:00Z"
  }
}
```

### Required Fields

- None (all fields are optional)

### Optional Fields

- `tin_number`: Tax Identification Number
- `nssf_number`: National Social Security Fund number
- `nhif_number`: National Health Insurance Fund number
- `wcf_number`: Workers Compensation Fund number
- `sdl_number`: Skills Development Levy number
- `passport_number`: Passport number (for expatriates)
- `passport_expiry_date`: ISO 8601 date format
- `work_permit_number`: Work permit number (for expatriates)
- `work_permit_expiry_date`: ISO 8601 date format

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 50.0,
    "completed_steps": [1, 2, 3, 4, 5],
    "finished_steps": [1, 2, 3, 4, 5],
    "unfinished_steps": [6, 7, 8, 9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T11:45:00Z"
  }
}
```

### Get Step 5 Data

**Response includes step 5:**
```json
{
  "steps": {
    "5": {
      "id": 1,
      "tin_number": "123456789",
      "nssf_number": "987654321",
      "nhif_number": "456789123",
      "wcf_number": "789123456",
      "sdl_number": "321654987",
      "passport_number": "A12345678",
      "passport_expiry_date": "2030-01-15T00:00:00Z",
      "work_permit_number": "WP123456",
      "work_permit_expiry_date": "2025-01-15T00:00:00Z"
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/5 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 5,
    "data": {
      "tin_number": "123456789",
      "nssf_number": "987654321",
      "nhif_number": "456789123"
    }
  }'
```

---

## Step 6: Documents

**Purpose:** Upload and manage employee documents (resume, certificates, contracts, etc.).

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/6
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/6
```

### Request Body

```json
{
  "step": 6,
  "data": {
    "documents": [
      {
        "document_type": "identity",
        "file_name": "national_id.pdf",
        "file_url": "https://example.com/files/national_id.pdf",
        "file_size": 245760,
        "mime_type": "application/pdf",
        "description": "National ID card"
      },
      {
        "document_type": "education",
        "file_name": "degree_certificate.pdf",
        "file_url": "https://example.com/files/degree.pdf",
        "file_size": 512000,
        "mime_type": "application/pdf",
        "description": "Bachelor's Degree Certificate"
      },
      {
        "document_type": "contract",
        "file_name": "employment_contract.pdf",
        "file_url": "https://example.com/files/contract.pdf",
        "file_size": 1024000,
        "mime_type": "application/pdf",
        "description": "Employment Contract"
      },
      {
        "document_type": "tax_statutory",
        "file_name": "tin_certificate.pdf",
        "file_url": "https://example.com/files/tin.pdf",
        "file_size": 128000,
        "mime_type": "application/pdf",
        "description": "TIN Certificate"
      }
    ]
  }
}
```

### Required Fields

- `documents` (array): Array of document objects
  - `document_type` (required): "identity", "work_permit", "education", "contract", "tax_statutory", or "other"
  - `file_name` (required): Name of the document file
  - `file_url` (required): URL to the uploaded file

### Optional Fields

- `file_size`: File size in bytes
- `mime_type`: MIME type of the file
- `description`: Description of the document

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 60.0,
    "completed_steps": [1, 2, 3, 4, 5, 6],
    "finished_steps": [1, 2, 3, 4, 5, 6],
    "unfinished_steps": [7, 8, 9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T12:00:00Z"
  }
}
```

### Get Step 6 Data

**Response includes step 6:**
```json
{
  "steps": {
    "6": [
      {
        "id": 1,
        "document_type": "identity",
        "file_name": "national_id.pdf",
        "file_url": "https://example.com/files/national_id.pdf",
        "file_size": 245760,
        "mime_type": "application/pdf",
        "description": "National ID card",
        "uploaded_at": "2024-01-15T12:00:00Z"
      },
      {
        "id": 2,
        "document_type": "education",
        "file_name": "degree_certificate.pdf",
        "file_url": "https://example.com/files/degree.pdf",
        "file_size": 512000,
        "mime_type": "application/pdf",
        "description": "Bachelor's Degree Certificate",
        "uploaded_at": "2024-01-15T12:00:00Z"
      }
    ]
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/6 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 6,
    "data": {
      "documents": [
        {
          "document_type": "identity",
          "file_name": "national_id.pdf",
          "file_url": "https://example.com/files/national_id.pdf"
        }
      ]
    }
  }'
```

**Important Notes:**
- Documents must be uploaded to a file storage service first
- Provide the `file_url` in the request
- Multiple documents can be uploaded in one request

---

## Step 7: Company Assets

**Purpose:** Assign company assets to the employee (laptop, phone, vehicle, etc.).

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/7
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/7
```

### Request Body

```json
{
  "step": 7,
  "data": {
    "assets": [
      {
        "asset_type": "laptop",
        "asset_name": "Dell Latitude 5520",
        "serial_number": "DL123456789",
        "asset_tag": "LAP001",
        "assigned_date": "2024-01-15T00:00:00Z",
        "expected_return_date": null,
        "condition": "new",
        "notes": "New laptop for development work"
      },
      {
        "asset_type": "phone",
        "asset_name": "iPhone 14 Pro",
        "serial_number": "IP987654321",
        "asset_tag": "PHN001",
        "assigned_date": "2024-01-15T00:00:00Z",
        "expected_return_date": null,
        "condition": "new",
        "notes": "Company phone for business use"
      },
      {
        "asset_type": "vehicle",
        "asset_name": "Toyota Corolla 2023",
        "serial_number": "VH456789123",
        "asset_tag": "VEH001",
        "assigned_date": "2024-01-15T00:00:00Z",
        "expected_return_date": null,
        "condition": "good",
        "notes": "Company vehicle for field work"
      }
    ]
  }
}
```

### Required Fields

- `assets` (array): Array of asset objects
  - `asset_type` (required): Type of asset (e.g., "laptop", "phone", "vehicle")
  - `asset_name` (required): Name/description of the asset

### Optional Fields

- `serial_number`: Serial number of the asset
- `asset_tag`: Asset tag/identifier
- `assigned_date`: ISO 8601 date format
- `expected_return_date`: ISO 8601 date format (null if not expected to return)
- `condition`: "new", "good", "fair", or "poor"
- `notes`: Additional notes about the asset

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 70.0,
    "completed_steps": [1, 2, 3, 4, 5, 6, 7],
    "finished_steps": [1, 2, 3, 4, 5, 6, 7],
    "unfinished_steps": [8, 9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T12:15:00Z"
  }
}
```

### Get Step 7 Data

**Response includes step 7:**
```json
{
  "steps": {
    "7": [
      {
        "id": 1,
        "asset_type": "laptop",
        "asset_name": "Dell Latitude 5520",
        "serial_number": "DL123456789",
        "asset_tag": "LAP001",
        "assigned_date": "2024-01-15T00:00:00Z",
        "expected_return_date": null,
        "condition": "new",
        "notes": "New laptop for development work"
      },
      {
        "id": 2,
        "asset_type": "phone",
        "asset_name": "iPhone 14 Pro",
        "serial_number": "IP987654321",
        "asset_tag": "PHN001",
        "assigned_date": "2024-01-15T00:00:00Z",
        "expected_return_date": null,
        "condition": "new",
        "notes": "Company phone for business use"
      }
    ]
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/7 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 7,
    "data": {
      "assets": [
        {
          "asset_type": "laptop",
          "asset_name": "Dell Latitude 5520",
          "serial_number": "DL123456789",
          "assigned_date": "2024-01-15T00:00:00Z"
        }
      ]
    }
  }'
```

---

## Step 8: Policies

**Purpose:** Assign leave and attendance policies to the employee.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/8
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/8
```

### Request Body

```json
{
  "step": 8,
  "data": {
    "leave_policy_id": 1,
    "attendance_policy_id": 2,
    "weekly_off_days": "Saturday,Sunday",
    "effective_date": "2024-01-15T00:00:00Z"
  }
}
```

### Required Fields

- None (all fields are optional)

### Optional Fields

- `leave_policy_id`: Leave policy ID (references leave_policies table)
- `attendance_policy_id`: Attendance policy ID (references attendance_policies table)
- `weekly_off_days`: Weekly off days (e.g., "Saturday,Sunday" or "Friday")
- `effective_date`: ISO 8601 date format

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 80.0,
    "completed_steps": [1, 2, 3, 4, 5, 6, 7, 8],
    "finished_steps": [1, 2, 3, 4, 5, 6, 7, 8],
    "unfinished_steps": [9, 10],
    "is_completed": false,
    "updated_at": "2024-01-15T12:30:00Z"
  }
}
```

### Get Step 8 Data

**Response includes step 8:**
```json
{
  "steps": {
    "8": {
      "id": 1,
      "leave_policy_id": 1,
      "attendance_policy_id": 2,
      "weekly_off_days": "Saturday,Sunday",
      "effective_date": "2024-01-15T00:00:00Z"
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/8 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 8,
    "data": {
      "weekly_off_days": "Saturday,Sunday",
      "effective_date": "2024-01-15T00:00:00Z"
    }
  }'
```

---

## Step 9: Emergency Contacts

**Purpose:** Add emergency contact information for the employee.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/9
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/9
```

### Request Body

```json
{
  "step": 9,
  "data": {
    "contacts": [
      {
        "contact_name": "Jane Doe",
        "relationship": "Spouse",
        "phone_number": "+255712345680",
        "alternate_phone": "+255712345681",
        "email": "jane.doe@example.com",
        "address": "123 Main Street, Dar es Salaam",
        "is_primary": true
      },
      {
        "contact_name": "John Smith",
        "relationship": "Father",
        "phone_number": "+255712345682",
        "alternate_phone": null,
        "email": null,
        "address": "456 Home Street, Arusha",
        "is_primary": false
      },
      {
        "contact_name": "Mary Johnson",
        "relationship": "Sister",
        "phone_number": "+255712345683",
        "alternate_phone": "+255712345684",
        "email": "mary.johnson@example.com",
        "address": "789 Family Street, Mwanza",
        "is_primary": false
      }
    ]
  }
}
```

### Required Fields

- `contacts` (array): Array of contact objects
  - `contact_name` (required): Name of emergency contact
  - `phone_number` (required): Primary phone number

### Optional Fields

- `relationship`: Relationship to employee (e.g., "Spouse", "Father", "Sister", "Friend")
- `alternate_phone`: Alternate phone number
- `email`: Email address (must be valid email format)
- `address`: Physical address
- `is_primary`: Boolean - true for primary emergency contact

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 90.0,
    "completed_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9],
    "finished_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9],
    "unfinished_steps": [10],
    "is_completed": false,
    "updated_at": "2024-01-15T12:45:00Z"
  }
}
```

### Get Step 9 Data

**Response includes step 9:**
```json
{
  "steps": {
    "9": [
      {
        "id": 1,
        "contact_name": "Jane Doe",
        "relationship": "Spouse",
        "phone_number": "+255712345680",
        "alternate_phone": "+255712345681",
        "email": "jane.doe@example.com",
        "address": "123 Main Street, Dar es Salaam",
        "is_primary": true
      },
      {
        "id": 2,
        "contact_name": "John Smith",
        "relationship": "Father",
        "phone_number": "+255712345682",
        "alternate_phone": null,
        "email": null,
        "address": "456 Home Street, Arusha",
        "is_primary": false
      }
    ]
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/9 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 9,
    "data": {
      "contacts": [
        {
          "contact_name": "Jane Doe",
          "relationship": "Spouse",
          "phone_number": "+255712345680",
          "is_primary": true
        }
      ]
    }
  }'
```

**Important Notes:**
- At least one emergency contact is recommended
- One contact should be marked as `is_primary: true`
- Multiple contacts can be added in one request

---

## Step 10: Internal Notes

**Purpose:** Add internal notes and comments about the employee onboarding.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/step/10
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/step/10
```

### Request Body

```json
{
  "step": 10,
  "data": {
    "notes": "Employee has previous experience in similar role. Recommended for fast-track onboarding. All documents verified. Ready for final approval."
  }
}
```

### Required Fields

- None (all fields are optional)

### Optional Fields

- `notes`: Internal notes about the employee (text)

### Response

```json
{
  "success": true,
  "message": "Step saved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP000001",
    "progress": 100.0,
    "completed_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    "finished_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    "unfinished_steps": [],
    "is_completed": false,
    "updated_at": "2024-01-15T13:00:00Z"
  }
}
```

### Get Step 10 Data

**Response includes step 10:**
```json
{
  "steps": {
    "10": {
      "notes": "Employee has previous experience in similar role. Recommended for fast-track onboarding. All documents verified. Ready for final approval."
    }
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/step/10 \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "step": 10,
    "data": {
      "notes": "Employee has previous experience in similar role. Recommended for fast-track onboarding."
    }
  }'
```

**Important Notes:**
- Step 10 is optional but recommended
- Notes are stored in the draft's `step_data` JSON field
- Progress becomes 100% after step 10 is saved

---

## Complete Onboarding

**Purpose:** Finalize onboarding and create the employee record.

### Endpoint

```
POST /api/v1/employees/onboarding/{employee_id}/complete
```

**Example:**
```
POST /api/v1/employees/onboarding/EMP000001/complete
```

### Request Body

None (employee_id in path is sufficient)

### Requirements

- **Minimum:** Steps 1 and 2 must be completed
- **Recommended:** Complete all 10 steps before finalizing

### Response

```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "id": 100,
    "employee_id": "EMP000001",
    "first_name": "John",
    "middle_name": "Michael",
    "last_name": "Doe",
    "work_email": "john.doe@company.com",
    "personal_email": "john.doe@example.com",
    "phone_number": "+255712345678",
    "department_id": 5,
    "position_id": 10,
    "location_id": 3,
    "hire_date": "2024-01-15T00:00:00Z",
    "employment_type": "full_time",
    "status": "active",
    "is_active": true,
    "created_at": "2024-01-15T13:00:00Z",
    "updated_at": "2024-01-15T13:00:00Z"
  }
}
```

### cURL Example

```bash
curl -X POST http://localhost:8080/api/v1/employees/onboarding/EMP000001/complete \
  -H "Authorization: Bearer YOUR_TOKEN"
```

**What Happens:**
1. Validates required steps (1 and 2) are completed
2. Creates final employee record in `employees` table
3. Migrates all draft data to employee record
4. Sets `is_completed = true` in draft
5. Sets `employee_id_final` in draft
6. Returns the created employee object

---

## Get Complete Onboarding Status

**Endpoint:**
```
GET /api/v1/employees/onboarding/{employee_id}
```

**Example:**
```
GET /api/v1/employees/onboarding/EMP000001
```

### Complete Response Example (All Steps Completed)

```json
{
  "success": true,
  "message": "Draft retrieved successfully",
  "data": {
    "draft_id": 1,
    "employee_id": "EMP000001",
    "progress": 100.0,
    "completed_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    "finished_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    "unfinished_steps": [],
    "is_completed": false,
    "steps": {
      "1": {
        "id": 1,
        "first_name": "John",
        "last_name": "Doe",
        "date_of_birth": "1990-01-15T00:00:00Z",
        "gender": "male",
        "personal_email": "john.doe@example.com",
        "mobile_number": "+255712345678",
        "addresses": [...]
      },
      "2": {
        "id": 1,
        "official_email": "john.doe@company.com",
        "department_id": 5,
        "position_id": 10,
        "employment_type": "full_time"
      },
      "3": {
        "id": 1,
        "annual_ctc": 24000000.00,
        "basic_salary": 15000000.00,
        "currency": "TZS"
      },
      "4": {
        "id": 1,
        "bank_name": "CRDB Bank",
        "account_number": "1234567890"
      },
      "5": {
        "id": 1,
        "tin_number": "123456789",
        "nssf_number": "987654321"
      },
      "6": [
        {
          "id": 1,
          "document_type": "identity",
          "file_name": "national_id.pdf",
          "file_url": "https://example.com/files/national_id.pdf"
        }
      ],
      "7": [
        {
          "id": 1,
          "asset_type": "laptop",
          "asset_name": "Dell Latitude 5520",
          "serial_number": "DL123456789"
        }
      ],
      "8": {
        "id": 1,
        "weekly_off_days": "Saturday,Sunday"
      },
      "9": [
        {
          "id": 1,
          "contact_name": "Jane Doe",
          "relationship": "Spouse",
          "phone_number": "+255712345680",
          "is_primary": true
        }
      ],
      "10": {
        "notes": "Employee has previous experience in similar role."
      }
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-15T13:00:00Z"
  }
}
```

### cURL Example

```bash
curl -X GET http://localhost:8080/api/v1/employees/onboarding/EMP000001 \
  -H "Authorization: Bearer YOUR_TOKEN"
```

---

## Quick Reference: All Steps Summary

| Step | Name | Endpoint | Required Fields | Progress |
|------|------|----------|----------------|----------|
| 1 | Personal Information | `POST /onboarding/draft` or `POST /onboarding/{id}/step/1` | `first_name`, `last_name` | 10% |
| 2 | Employment Details | `POST /onboarding/{id}/step/2` | None (recommended: `department_id`) | 20% |
| 3 | Salary & CTC | `POST /onboarding/{id}/step/3` | None | 30% |
| 4 | Bank Account | `POST /onboarding/{id}/step/4` | None | 40% |
| 5 | Statutory | `POST /onboarding/{id}/step/5` | None | 50% |
| 6 | Documents | `POST /onboarding/{id}/step/6` | `documents[].document_type`, `documents[].file_name`, `documents[].file_url` | 60% |
| 7 | Assets | `POST /onboarding/{id}/step/7` | `assets[].asset_type`, `assets[].asset_name` | 70% |
| 8 | Policies | `POST /onboarding/{id}/step/8` | None | 80% |
| 9 | Emergency Contacts | `POST /onboarding/{id}/step/9` | `contacts[].contact_name`, `contacts[].phone_number` | 90% |
| 10 | Internal Notes | `POST /onboarding/{id}/step/10` | None | 100% |
| Complete | Finalize | `POST /onboarding/{id}/complete` | Steps 1 & 2 minimum | - |

---

## Complete Example: Full Onboarding Flow

Here's a complete example showing all 10 steps:

```javascript
const BASE_URL = 'http://localhost:8080/api/v1/employees/onboarding';
const token = 'YOUR_JWT_TOKEN';

// ============================================
// STEP 1: Personal Information
// ============================================
const step1Response = await fetch(`${BASE_URL}/draft`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 1,
    data: {
      first_name: "John",
      middle_name: "Michael",
      last_name: "Doe",
      date_of_birth: "1990-01-15T00:00:00Z",
      gender: "male",
      marital_status: "single",
      blood_group: "O+",
      nationality: "Tanzanian",
      personal_email: "john.doe@example.com",
      mobile_number: "+255712345678",
      alternate_number: "+255712345679",
      current_address: "123 Main Street, Kinondoni",
      permanent_address: "456 Home Street, Arusha",
      city: "Dar es Salaam",
      state: "Dar es Salaam",
      postal_code: "11101",
      country: "Tanzania"
    }
  })
});

const { data: step1Data } = await step1Response.json();
const employeeID = step1Data.employee_id; // "EMP000001"
console.log('✅ Step 1 completed. Employee ID:', employeeID);
console.log('Progress:', step1Data.progress); // 10%

// ============================================
// STEP 2: Employment Details
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/2`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 2,
    data: {
      official_email: "john.doe@company.com",
      date_of_joining: "2024-01-15T00:00:00Z",
      department_id: 5,
      position_id: 10,
      grade: "G5",
      reporting_manager_id: 20,
      employment_type: "full_time",
      location_id: 3,
      shift: "Day",
      work_phone: "+255222345678",
      probation_period_days: 90,
      expected_confirmation_date: "2024-04-15T00:00:00Z"
    }
  })
});
console.log('✅ Step 2 completed. Progress: 20%');

// ============================================
// STEP 3: Salary & CTC
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/3`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 3,
    data: {
      annual_ctc: 24000000.00,
      ctc_effective_date: "2024-01-15T00:00:00Z",
      currency: "TZS",
      basic_salary: 15000000.00,
      house_rent_allowance: 3000000.00,
      transport_allowance: 2000000.00,
      special_allowance: 2000000.00,
      other_allowances: 2000000.00,
      income_tax: 2400000.00,
      provident_fund: 1200000.00,
      professional_tax: 50000.00,
      other_deductions: 0.00
    }
  })
});
console.log('✅ Step 3 completed. Progress: 30%');

// ============================================
// STEP 4: Bank Account
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/4`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 4,
    data: {
      bank_name: "CRDB Bank",
      account_holder_name: "John Doe",
      account_number: "1234567890",
      account_type: "savings",
      branch_name: "Kinondoni Branch",
      swift_code: "CRDBTZTZ"
    }
  })
});
console.log('✅ Step 4 completed. Progress: 40%');

// ============================================
// STEP 5: Statutory Requirements
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/5`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 5,
    data: {
      tin_number: "123456789",
      nssf_number: "987654321",
      nhif_number: "456789123",
      wcf_number: "789123456",
      sdl_number: "321654987"
    }
  })
});
console.log('✅ Step 5 completed. Progress: 50%');

// ============================================
// STEP 6: Documents
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/6`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 6,
    data: {
      documents: [
        {
          document_type: "identity",
          file_name: "national_id.pdf",
          file_url: "https://example.com/files/national_id.pdf",
          file_size: 245760,
          mime_type: "application/pdf",
          description: "National ID card"
        },
        {
          document_type: "education",
          file_name: "degree_certificate.pdf",
          file_url: "https://example.com/files/degree.pdf",
          file_size: 512000,
          mime_type: "application/pdf",
          description: "Bachelor's Degree Certificate"
        }
      ]
    }
  })
});
console.log('✅ Step 6 completed. Progress: 60%');

// ============================================
// STEP 7: Company Assets
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/7`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 7,
    data: {
      assets: [
        {
          asset_type: "laptop",
          asset_name: "Dell Latitude 5520",
          serial_number: "DL123456789",
          asset_tag: "LAP001",
          assigned_date: "2024-01-15T00:00:00Z",
          condition: "new",
          notes: "New laptop for development work"
        },
        {
          asset_type: "phone",
          asset_name: "iPhone 14 Pro",
          serial_number: "IP987654321",
          asset_tag: "PHN001",
          assigned_date: "2024-01-15T00:00:00Z",
          condition: "new",
          notes: "Company phone for business use"
        }
      ]
    }
  })
});
console.log('✅ Step 7 completed. Progress: 70%');

// ============================================
// STEP 8: Policies
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/8`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 8,
    data: {
      leave_policy_id: 1,
      attendance_policy_id: 2,
      weekly_off_days: "Saturday,Sunday",
      effective_date: "2024-01-15T00:00:00Z"
    }
  })
});
console.log('✅ Step 8 completed. Progress: 80%');

// ============================================
// STEP 9: Emergency Contacts
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/9`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 9,
    data: {
      contacts: [
        {
          contact_name: "Jane Doe",
          relationship: "Spouse",
          phone_number: "+255712345680",
          alternate_phone: "+255712345681",
          email: "jane.doe@example.com",
          address: "123 Main Street, Dar es Salaam",
          is_primary: true
        },
        {
          contact_name: "John Smith",
          relationship: "Father",
          phone_number: "+255712345682",
          address: "456 Home Street, Arusha",
          is_primary: false
        }
      ]
    }
  })
});
console.log('✅ Step 9 completed. Progress: 90%');

// ============================================
// STEP 10: Internal Notes
// ============================================
await fetch(`${BASE_URL}/${employeeID}/step/10`, {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    step: 10,
    data: {
      notes: "Employee has previous experience in similar role. Recommended for fast-track onboarding. All documents verified. Ready for final approval."
    }
  })
});
console.log('✅ Step 10 completed. Progress: 100%');



## Field Validation Rules

### Step 1: Personal Information
- `first_name`: Required, minimum 2 characters
- `last_name`: Required, minimum 2 characters
- `gender`: Must be "male", "female", or "other"
- `marital_status`: Must be "single", "married", "divorced", or "widowed"
- `blood_group`: Must be "A+", "A-", "B+", "B-", "AB+", "AB-", "O+", or "O-"
- `personal_email`: Must be valid email format

### Step 2: Employment Details
- `official_email`: Must be valid email format
- `employment_type`: Must be "full_time", "part_time", "contract", or "intern"
- `probation_period_days`: Minimum 90 days

### Step 3: Salary & CTC
- `currency`: Must be exactly 3 characters (e.g., "TZS", "USD")

### Step 4: Bank Account
- `account_type`: Must be "savings" or "current"

### Step 6: Documents
- `document_type`: Must be "identity", "work_permit", "education", "contract", "tax_statutory", or "other"
- `file_name`: Required
- `file_url`: Required

### Step 7: Assets
- `asset_type`: Required
- `asset_name`: Required

### Step 9: Emergency Contacts
- `contact_name`: Required
- `phone_number`: Required
- `email`: Must be valid email format if provided

---

## Error Responses

### Common Error Format

```json
{
  "success": false,
  "message": "Error message here",
  "error": "Detailed error description"
}
```

### Step-Specific Errors

#### Step 1 Errors
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Key: 'Step1PersonalInfoRequest.FirstName' Error:Field validation for 'FirstName' failed on the 'required' tag"
}
```

#### Step 2 Errors
```json
{
  "success": false,
  "message": "department not found",
  "error": "department not found"
}
```

#### Step 6 Errors
```json
{
  "success": false,
  "message": "Validation failed",
  "error": "Key: 'DocumentUpload.FileName' Error:Field validation for 'FileName' failed on the 'required' tag"
}
```

---

## Tips and Best Practices

1. **Always Save Employee ID**: After step 1, save the `employee_id` immediately
2. **Check Progress**: Use `GET /onboarding/{employee_id}` to check progress before proceeding
3. **Validate Data**: Validate all required fields before sending requests
4. **Handle Errors**: Implement proper error handling for validation and network errors
5. **Resume Capability**: Use `GET /onboarding/{employee_id}` to resume onboarding and pre-fill forms
6. **Progress Indicators**: Show progress bar and step status to users
7. **Step Navigation**: Use `finished_steps` and `unfinished_steps` to enable/disable step navigation

---

**Last Updated:** 2024-01-15  
**Version:** 2.0.0

## Best Practices

### 1. Always Store Employee ID

```javascript
// After step 1, save employee_id
const employeeID = response.data.employee_id;
localStorage.setItem('onboarding_employee_id', employeeID);
// Or in your state management
```

### 2. Use Employee ID for All Operations

```javascript
// ✅ Good - Use employee_id
POST /api/v1/employees/onboarding/EMP000001/step/2

// ❌ Bad - Don't use draft_id
POST /api/v1/employees/onboarding/draft/1/step/2
```


---

## Error Handling

### Common Errors

#### 1. Employee ID Not Found

```json
{
  "success": false,
  "message": "draft not found for this employee ID",
  "error": "draft not found for this employee ID"
}
```

**Solution:** Verify employee ID is correct. Employee ID is generated in step 1.

#### 2. Invalid Step Number

```json
{
  "success": false,
  "message": "Invalid step number. Must be between 1 and 10",
  "error": "invalid step number"
}
```

**Solution:** Use step numbers 1-10 only.

#### 3. Step Already Completed

You can save the same step multiple times - it will update the existing data.

---

## Migration from Draft ID

If you were using draft_id before, here's how to migrate:

### Old Way (Draft ID)
```javascript
// Create draft
POST /onboarding/draft
// Get draft_id: 1

// Save step
POST /onboarding/draft/1/step/2

// Get status
GET /onboarding/draft/1
```

### New Way (Employee ID)
```javascript
// Create draft and save step 1
POST /onboarding/draft (with step 1 data)
// Get employee_id: "EMP000001"

// Save step
POST /onboarding/EMP000001/step/2

// Get status
GET /onboarding/EMP000001
```

**Benefits:**
- ✅ More intuitive (employee ID is meaningful)
- ✅ Easier to track (no need to remember draft_id)
- ✅ Better for users (employee ID is visible)
- ✅ Easier to resume (just use employee ID)

---

## Summary

### Key Points

1. **Employee ID is Auto-Generated**: When you save step 1, employee ID is automatically generated (format: EMP000001)

2. **Use Employee ID for Everything**: All operations use employee ID:
   - `GET /onboarding/{employee_id}` - Get status
   - `POST /onboarding/{employee_id}/step/{step}` - Save step
   - `POST /onboarding/{employee_id}/complete` - Complete onboarding

3. **Progress is Automatic**: Progress, finished steps, and unfinished steps are automatically calculated and returned

4. **Easy Resume**: Just use employee ID to resume onboarding from where you left off

5. **Clear Status**: Response includes:
   - `progress`: Completion percentage
   - `finished_steps`: Steps with data
   - `unfinished_steps`: Steps without data
   - `steps`: All saved step data

### Quick Reference

| Operation | Endpoint | Method |
|-----------|----------|--------|
| Create draft + Step 1 | `/onboarding/draft` | POST |
| Get status | `/onboarding/{employee_id}` | GET |
| Save step | `/onboarding/{employee_id}/step/{step}` | POST |
| Complete | `/onboarding/{employee_id}/complete` | POST |

---

**Last Updated:** 2024-01-15  
**Version:** 2.0.0
