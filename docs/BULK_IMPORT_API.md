# Bulk Import API — Employees, Departments & Positions

> Download Excel templates, fill them in, and upload to bulk-create departments, positions, and employees.

---

## Recommended Upload Order

```
 Step 1: Departments   →   Step 2: Positions   →   Step 3: Employees
     ▲                         ▲                        ▲
     │                         │                        │
  No dependencies      Needs dept codes         Needs dept & position codes
```

Upload **departments first**, then **positions** (which reference department codes), then **employees** (which reference both department and position codes).

---

## API Endpoints

| # | Method | Endpoint                                  | Description                     | Access   |
|---|--------|-------------------------------------------|---------------------------------|----------|
| 1 | `GET`  | `/api/v1/bulk-import/templates/departments` | Download department Excel template | HR/Admin |
| 2 | `GET`  | `/api/v1/bulk-import/templates/positions`   | Download position Excel template   | HR/Admin |
| 3 | `GET`  | `/api/v1/bulk-import/templates/employees`   | Download employee Excel template   | HR/Admin |
| 4 | `POST` | `/api/v1/bulk-import/departments`           | Import departments from Excel      | HR/Admin |
| 5 | `POST` | `/api/v1/bulk-import/positions`             | Import positions from Excel        | HR/Admin |
| 6 | `POST` | `/api/v1/bulk-import/employees`             | Import employees from Excel        | HR/Admin |

---

## 1. Download Templates

Each template is a styled `.xlsx` file with:
- **Header row** (row 1) — column names, required fields in red, optional in blue
- **Sample data** (row 2+) — example values to show the expected format
- **Instructions sheet** — field descriptions, valid values, and tips

### Department Template

```
GET /api/v1/bulk-import/templates/departments
```

Downloads `department_bulk_upload_template.xlsx` with columns:

| Column                | Required | Example              | Description                      |
|-----------------------|----------|----------------------|----------------------------------|
| Code                  | Yes      | ENG                  | Unique short code (max 20 chars) |
| Name                  | Yes      | Engineering          | Department name (max 100 chars)  |
| Description           | No       | Software engineering | Free text                        |
| Level                 | No       | department           | company, business_unit, department, team |
| Department Type       | No       | core                 | core, support, operational, strategic |
| Parent Department Code| No       | TECH                 | Code of parent dept (for hierarchy) |
| Employee Capacity     | No       | 50                   | Expected headcount               |

### Position Template

```
GET /api/v1/bulk-import/templates/positions
```

Downloads `position_bulk_upload_template.xlsx` with columns:

| Column                    | Required | Example            | Description                         |
|---------------------------|----------|--------------------|-------------------------------------|
| Code                      | Yes      | SE-001             | Unique short code (max 20 chars)    |
| Title                     | Yes      | Software Engineer  | Position title (max 100 chars)      |
| Grade                     | No       | L3                 | Grade level                         |
| Department Code           | No       | ENG                | Must match existing department code |
| Reports To Position Code  | No       | SSE-001            | Parent position code                |
| Budgeted Headcount        | No       | 10                 | Expected number of hires            |
| Key Competencies          | No       | Go, Docker         | Comma-separated skills              |

### Employee Template

```
GET /api/v1/bulk-import/templates/employees
```

Downloads `employee_bulk_upload_template.xlsx` with columns:

| Column                     | Required | Example              | Description                        |
|----------------------------|----------|----------------------|------------------------------------|
| Employee ID                | Yes      | EMP001               | Unique employee code               |
| First Name                 | Yes      | John                 | First name                         |
| Middle Name                | No       | Michael              | Middle name                        |
| Last Name                  | Yes      | Doe                  | Last name                          |
| Work Email                 | No       | john.doe@company.com | Must be unique                     |
| Personal Email             | No       | john@gmail.com       |                                    |
| Phone Number               | No       | +255712345678        |                                    |
| Gender                     | No       | Male                 | Male, Female, Other                |
| Date of Birth              | No       | 1990-05-15           | YYYY-MM-DD format                  |
| Nationality                | No       | Tanzanian            |                                    |
| Marital Status             | No       | Single               | Single, Married, Divorced, Widowed |
| Department Code            | No       | ENG                  | Must match existing dept code      |
| Position Code              | No       | SE-001               | Must match existing position code  |
| Employment Type            | No       | full_time            | full_time, part_time, contract, intern |
| Hire Date                  | No       | 2024-01-15           | YYYY-MM-DD format                  |
| Grade                      | No       | L3                   | Grade level                        |
| Shift                      | No       | Day                  | Day, Night, Rotating               |
| Manager Employee ID        | No       | EMP000               | Employee ID of manager             |
| Salary                     | No       | 3500000              | Numeric value                      |
| Currency                   | No       | TZS                  | 3-letter code (default: TZS)       |
| Emergency Contact Name     | No       | Jane Doe             |                                    |
| Emergency Contact Phone    | No       | +255712345679        |                                    |
| Emergency Contact Relation | No       | Spouse               |                                    |
| Notes                      | No       | From Arusha branch   |                                    |

---

## 2. Import Data

Upload a filled Excel file to import records.

### Import Departments

```
POST /api/v1/bulk-import/departments
Content-Type: multipart/form-data
```

**Form field:** `file` — the filled `.xlsx` file

**Success Response (200):**
```json
{
  "success": true,
  "message": "Department import completed: 5 imported, 0 skipped out of 5 rows",
  "data": {
    "total_rows": 5,
    "imported": 5,
    "skipped": 0
  }
}
```

**Partial Success (200):**
```json
{
  "success": true,
  "message": "Department import completed: 3 imported, 2 skipped out of 5 rows",
  "data": {
    "total_rows": 5,
    "imported": 3,
    "skipped": 2,
    "errors": [
      { "row": 4, "field": "Code", "message": "Department code 'HR' already exists" },
      { "row": 5, "field": "Parent Department Code", "message": "Parent 'UNKNOWN' not found" }
    ]
  }
}
```

**All Failed (422):**
```json
{
  "success": false,
  "message": "No departments imported — check errors",
  "data": {
    "total_rows": 2,
    "imported": 0,
    "skipped": 2,
    "errors": [
      { "row": 2, "field": "Code", "message": "Department code 'ENG' already exists" },
      { "row": 3, "message": "Code and Name are required" }
    ]
  }
}
```

### Import Positions

```
POST /api/v1/bulk-import/positions
Content-Type: multipart/form-data
```

Same response structure as departments.

### Import Employees

```
POST /api/v1/bulk-import/employees
Content-Type: multipart/form-data
```

Same response structure as departments.

---

## Validation Rules

### Departments
- `Code` and `Name` are required
- `Code` must be unique (not already in the database)
- `Parent Department Code`, if provided, must reference an existing or earlier-in-file department

### Positions
- `Code` and `Title` are required
- `Code` must be unique
- `Department Code`, if provided, must reference an existing department
- `Reports To Position Code`, if provided, must reference an existing or earlier-in-file position

### Employees
- `Employee ID`, `First Name`, and `Last Name` are required
- `Employee ID` must be unique
- `Work Email`, if provided, must be unique
- `Department Code` and `Position Code`, if provided, must reference existing records
- `Manager Employee ID`, if provided, must reference an existing or earlier-in-file employee
- Date fields accept: `YYYY-MM-DD`, `MM/DD/YYYY`, `DD/MM/YYYY`

---

## Frontend Integration

### Workflow

```
┌─────────────────────────────────────────────┐
│          Bulk Import Page                    │
│                                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ Download  │  │ Download  │  │ Download  │  │
│  │ Dept      │  │ Position  │  │ Employee  │  │
│  │ Template  │  │ Template  │  │ Template  │  │
│  └─────┬────┘  └─────┬────┘  └─────┬────┘  │
│        │              │              │        │
│        ▼              ▼              ▼        │
│   User fills     User fills     User fills    │
│   in Excel       in Excel       in Excel      │
│        │              │              │        │
│        ▼              ▼              ▼        │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  │
│  │ Upload   │  │ Upload   │  │ Upload   │  │
│  │ Depts    │  │ Positions│  │ Employees│  │
│  └─────┬────┘  └─────┬────┘  └─────┬────┘  │
│        │              │              │        │
│        ▼              ▼              ▼        │
│   Show results   Show results   Show results  │
│   (imported,     (imported,     (imported,     │
│    skipped,       skipped,       skipped,      │
│    errors)        errors)        errors)       │
└─────────────────────────────────────────────┘
```

### Download Template (JavaScript)

```javascript
const downloadTemplate = async (type) => {
  // type: "employees", "departments", or "positions"
  const res = await fetch(`/api/v1/bulk-import/templates/${type}`, {
    headers: { 'Authorization': `Bearer ${token}` },
  });
  const blob = await res.blob();
  const url = window.URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${type}_bulk_upload_template.xlsx`;
  a.click();
};
```

### Upload Filled Excel (JavaScript)

```javascript
const importData = async (type, file) => {
  const formData = new FormData();
  formData.append('file', file);

  const res = await fetch(`/api/v1/bulk-import/${type}`, {
    method: 'POST',
    headers: { 'Authorization': `Bearer ${token}` },
    body: formData,
  });
  return await res.json();
  // result.data.imported → number created
  // result.data.errors → array of { row, field, message }
};
```

---

## Error Responses

| Code | Scenario                                  |
|------|-------------------------------------------|
| 400  | No file uploaded                          |
| 400  | Not an Excel file (.xlsx / .xls)          |
| 400  | File too large (> 10 MB)                  |
| 400  | No data rows found                        |
| 200  | Partial success (some imported, some skipped) |
| 422  | All rows failed validation                |

---

## Limits

| Limit            | Value     |
|------------------|-----------|
| Max file size    | 10 MB     |
| Max rows (employees)    | 500       |
| Max rows (departments)  | 200       |
| Max rows (positions)    | 200       |
| File types       | .xlsx, .xls |

---

## File Structure

```
internal/modules/bulk_import/
├── handlers/
│   └── bulk_import_handler.go    ← 6 HTTP handlers
└── services/
    ├── template_service.go       ← Excel template generation
    └── import_service.go         ← Excel parsing + validation + DB insert
```
