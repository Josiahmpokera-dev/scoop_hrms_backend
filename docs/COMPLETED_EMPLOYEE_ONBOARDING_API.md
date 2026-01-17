# Completed Employee Onboarding API Documentation

## Overview

The Completed Employee Onboarding API allows you to view all onboarding data for employees who have already completed the onboarding process. This includes all information they provided during the 10-step onboarding process.

## Base URL
All API endpoints use the base URL: `{{BASE_URL}}/api/v1/employees/onboarding`

## Authentication
All requests require Bearer token authentication:
```
Authorization: Bearer <access_token>
```

---

## Get Completed Employee Onboarding Data API

### Endpoint
```
POST /api/v1/employees/onboarding/completed-employee
```

### Description
Retrieves all onboarding data for a completed employee by their `employee_id`. This includes all data from all 10 onboarding steps that were filled during the onboarding process.

### Request Body
```json
{
  "employee_id": "EMP001"
}
```

### Required Fields
- `employee_id` (string, required): Employee ID (e.g., "EMP001")

### Success Response (200)
```json
{
  "success": true,
  "message": "Completed employee onboarding data retrieved successfully",
  "data": {
    "employee_id": "EMP001",
    "employee_db_id": 1,
    "progress": 100.0,
    "completed_steps": [1, 2, 3, 4, 5, 6, 7, 8, 9, 10],
    "is_completed": true,
    "steps": {
      "1": {
        "photo_url": "https://example.com/photos/john.jpg",
        "first_name": "John",
        "middle_name": "Michael",
        "last_name": "Doe",
        "date_of_birth": "1990-05-15T00:00:00Z",
        "gender": "male",
        "marital_status": "married",
        "blood_group": "O+",
        "nationality": "Tanzanian",
        "personal_email": "john.doe@example.com",
        "mobile_number": "+255123456789",
        "alternate_number": "+255987654321",
        "addresses": [
          {
            "id": 1,
            "address_type": "current",
            "address_line1": "123 Main Street",
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
            "state": "Arusha",
            "postal_code": "23100",
            "country": "Tanzania"
          }
        ]
      },
      "2": {
        "employee_id": "EMP001",
        "official_email": "john.doe@company.com",
        "date_of_joining": "2024-01-20T00:00:00Z",
        "department_id": 5,
        "position_id": 10,
        "grade": "L3",
        "reporting_manager_id": 3,
        "employment_type": "full_time",
        "location_id": 2,
        "shift": "Day",
        "work_phone": "+255123456789",
        "probation_period_days": 90,
        "expected_confirmation_date": "2024-04-20T00:00:00Z"
      },
      "3": {
        "id": 1,
        "annual_ctc": 5000000.00,
        "ctc_effective_date": "2024-01-20T00:00:00Z",
        "currency": "TZS",
        "basic_salary": 3000000.00,
        "house_rent_allowance": 500000.00,
        "transport_allowance": 200000.00,
        "special_allowance": 300000.00,
        "other_allowances": 1000000.00,
        "gross_salary": 5000000.00,
        "income_tax": 500000.00,
        "provident_fund": 250000.00,
        "professional_tax": 50000.00,
        "other_deductions": 0.00,
        "net_salary": 4200000.00
      },
      "4": [
        {
          "id": 1,
          "bank_name": "CRDB Bank",
          "account_holder_name": "John Doe",
          "account_number": "1234567890",
          "account_type": "savings",
          "branch_name": "Dar es Salaam Branch",
          "swift_code": "CRDBTZTZ"
        }
      ],
      "5": {
        "id": 1,
        "tin_number": "123456789",
        "nssf_number": "987654321",
        "nhif_number": "456789123",
        "wcf_number": "789123456",
        "sdl_number": "321654987",
        "passport_number": "A12345678",
        "passport_expiry_date": "2030-12-31T00:00:00Z",
        "work_permit_number": "WP2024001",
        "work_permit_expiry_date": "2025-12-31T00:00:00Z"
      },
      "6": [
        {
          "id": 1,
          "document_type": "national_id",
          "file_name": "national_id.pdf",
          "file_url": "https://example.com/documents/national_id.pdf",
          "issue_date": "2020-01-15T00:00:00Z",
          "expiry_date": "2030-01-15T00:00:00Z",
          "is_verified": true
        },
        {
          "id": 2,
          "document_type": "degree_certificate",
          "file_name": "degree.pdf",
          "file_url": "https://example.com/documents/degree.pdf",
          "issue_date": "2015-06-30T00:00:00Z",
          "expiry_date": null,
          "is_verified": true
        }
      ],
      "7": [
        {
          "id": 1,
          "asset_type": "laptop",
          "asset_name": "Dell Latitude 5520",
          "serial_number": "DL552012345",
          "assigned_date": "2024-01-20T00:00:00Z",
          "expected_return_date": null,
          "condition": "excellent",
          "notes": "Brand new laptop"
        }
      ],
      "8": {
        "id": 1,
        "leave_policy_id": 1,
        "attendance_policy_id": 1,
        "weekly_off_days": ["Saturday", "Sunday"],
        "effective_date": "2024-01-20T00:00:00Z"
      },
      "9": [
        {
          "id": 1,
          "contact_name": "Jane Doe",
          "relationship": "spouse",
          "phone_number": "+255987654321",
          "alternate_phone": "+255987654322",
          "email": "jane.doe@example.com",
          "address": "123 Main Street, Dar es Salaam",
          "is_primary": true
        },
        {
          "id": 2,
          "contact_name": "John Doe Sr.",
          "relationship": "father",
          "phone_number": "+255987654323",
          "alternate_phone": null,
          "email": null,
          "address": "456 Home Street, Arusha",
          "is_primary": false
        }
      ],
      "10": {
        "notes": "Employee completed all onboarding steps successfully. Ready for work."
      }
    },
    "created_at": "2024-01-15T10:30:00Z",
    "updated_at": "2024-01-20T14:30:00Z"
  }
}
```

### Response Fields

| Field | Type | Description |
|-------|------|-------------|
| `employee_id` | string | Employee ID (e.g., "EMP001") |
| `employee_db_id` | integer | Database ID of the employee record |
| `progress` | float | Completion percentage (always 100.0 for completed employees) |
| `completed_steps` | array | List of completed step numbers (always [1,2,3,4,5,6,7,8,9,10]) |
| `is_completed` | boolean | Whether onboarding is completed (always true) |
| `steps` | object | Step data keyed by step number (1-10) |
| `created_at` | datetime | When employee record was created |
| `updated_at` | datetime | When employee record was last updated |

### Step Data Structure

The `steps` object contains data for each completed step:

#### Step 1: Personal Information
- Personal details: name, date of birth, gender, marital status, blood group, nationality
- Contact information: personal email, mobile number, alternate number
- Addresses: current and permanent addresses

#### Step 2: Employment Details
- Employee ID, official email, date of joining
- Department, position, grade
- Reporting manager, employment type
- Location, shift, work phone
- Probation period information

#### Step 3: Salary & CTC
- Annual CTC, currency
- Salary components: basic, allowances
- Deductions: tax, provident fund, etc.
- Calculated gross and net salary

#### Step 4: Bank Account Details
- Array of bank accounts
- Bank name, account details, branch, SWIFT code

#### Step 5: Statutory Requirements
- TIN, NSSF, NHIF, WCF, SDL numbers
- Passport and work permit information (for expatriates)

#### Step 6: Documents
- Array of uploaded documents
- Document type, file URL, issue/expiry dates, verification status

#### Step 7: Company Assets
- Array of assigned assets
- Asset type, name, serial number, assignment dates, condition

#### Step 8: Policies
- Leave and attendance policies
- Weekly off days, effective date

#### Step 9: Emergency Contacts
- Array of emergency contacts
- Contact details, relationship, primary contact flag

#### Step 10: Internal Notes
- Internal notes and comments

### Example Request (cURL)
```bash
curl -X POST "http://localhost:8080/api/v1/employees/onboarding/completed-employee" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "employee_id": "EMP001"
  }'
```

### Example Request (JavaScript)
```javascript
const response = await fetch('http://localhost:8080/api/v1/employees/onboarding/completed-employee', {
  method: 'POST',
  headers: {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    employee_id: 'EMP001'
  })
});

const data = await response.json();
console.log(data.data); // Completed employee onboarding data
```

### Example Request (React)
```javascript
const getCompletedEmployeeData = async (employeeID) => {
  try {
    const response = await fetch(`${BASE_URL}/api/v1/employees/onboarding/completed-employee`, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${token}`,
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        employee_id: employeeID
      })
    });

    const result = await response.json();
    if (result.success) {
      return result.data;
    } else {
      throw new Error(result.message);
    }
  } catch (error) {
    console.error('Error fetching completed employee data:', error);
    throw error;
  }
};

// Usage
const employeeData = await getCompletedEmployeeData('EMP001');
console.log('Step 1 (Personal Info):', employeeData.steps['1']);
console.log('Step 2 (Employment):', employeeData.steps['2']);
```

### Error Responses

#### Employee Not Found
```json
{
  "success": false,
  "message": "employee not found"
}
```

#### Tenant Mismatch
```json
{
  "success": false,
  "message": "employee does not belong to your tenant"
}
```

#### Validation Error
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "employee_id is required in request body"
}
```

---

## Use Cases

### 1. View Employee Onboarding History
Use this API to view all information that was collected during an employee's onboarding process:
- Personal information and contact details
- Employment details and department assignment
- Salary and compensation structure
- Bank account information
- Statutory compliance documents
- Assigned assets and equipment
- Policy acknowledgments
- Emergency contacts

### 2. Employee Profile Review
HR can review complete employee profiles including:
- All personal and professional information
- Documents and certifications
- Asset assignments
- Policy compliance

### 3. Data Verification
Verify that all onboarding data was correctly migrated from draft to employee record:
- Compare draft data with completed employee data
- Ensure data integrity after onboarding completion
- Audit trail for compliance

### 4. Reporting and Analytics
Extract onboarding data for:
- Compliance reporting
- HR analytics
- Employee data exports
- Integration with other systems

---

## Data Source Priority

The API loads data from dedicated onboarding tables first, then falls back to the main employee table if needed:

1. **Step 1 (Personal Info)**: `employee_basic_information` → `employees` table
2. **Step 2 (Employment)**: `employee_employment_details` → `employees` table
3. **Step 3 (Salary)**: `employee_salary_components` table
4. **Step 4 (Bank)**: `employee_bank_accounts` table
5. **Step 5 (Statutory)**: `employee_statutory_info` table
6. **Step 6 (Documents)**: `employee_documents` table
7. **Step 7 (Assets)**: `employee_assets` table
8. **Step 8 (Policies)**: `employee_policies` table
9. **Step 9 (Emergency Contacts)**: `employee_emergency_contacts` → `employees` table (legacy fields)
10. **Step 10 (Notes)**: `employees.notes` field

---

## Integration with Other APIs

This API works alongside:

1. **List Draft Employees**: `GET /onboarding/draft-employees` - View incomplete onboardings
2. **Get Draft Employee**: `GET /onboarding/draft-employee/:employee_id` - View incomplete onboarding data
3. **List Employees**: `GET /employees` - View all employees (includes onboarding status)
4. **Get Employee**: `GET /employees/:id` - View employee basic information

---

## Quick Reference

| Operation | Endpoint | Method | Description |
|-----------|----------|--------|-------------|
| Get Completed Employee Onboarding | `/onboarding/completed-employee` | POST | Get all onboarding data for completed employee (employee_id in body) |
| Get Draft Employee | `/onboarding/draft-employee/:employee_id` | GET | Get draft employee data to continue onboarding |
| List Draft Employees | `/onboarding/draft-employees` | GET | List all incomplete draft employees |

---

## Notes

- Only employees who have **completed** onboarding will return data
- If an employee is still in draft, use the draft employee API instead
- All step data is loaded from dedicated database tables
- Progress is always 100% for completed employees
- All 10 steps are always marked as completed
- Data is loaded from both onboarding tables and employee table (with fallback)

---

**Last Updated:** 2024-01-15  
**Version:** 1.0.0
