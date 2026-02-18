# Side Menu Visibility — Implementation Documentation

## Overview

The Side Menu Visibility feature lets **Admin / Super Admin** control which navigation items are shown or hidden in the application. Changes apply globally to all users.

**Base URL:** `/api/v1/settings`

---

## Architecture

```
internal/modules/settings/
├── models/
│   └── menu_visibility.go        # DB model, DTOs, full menu tree definition, cascade helpers
├── repositories/
│   └── menu_visibility_repository.go   # GORM CRUD with upsert (ON CONFLICT)
├── services/
│   └── menu_visibility_service.go      # Business logic, cascade rules, validation
└── handlers/
    └── menu_visibility_handler.go      # 5 HTTP handlers
```

### Database Table

**`menu_visibility_settings`** (auto-migrated)

| Column | Type | Description |
|--------|------|-------------|
| `id` | bigint (PK) | Auto-increment |
| `menu_key` | varchar(100) | Unique menu item key |
| `visible` | boolean | `true` = shown, `false` = hidden |
| `updated_by` | bigint | FK → users.id |
| `created_at` | timestamp | Row creation |
| `updated_at` | timestamp | Last modification |

Only items that have been explicitly toggled are stored. Items not in the table default to **visible**.

---

## API Endpoints (5 total)

### 1. Get All Menu Visibility Settings

**`GET /api/v1/settings/menu-visibility`**

**Auth:** Admin / Super Admin

Returns the full menu tree (78 items) with visibility flags and a summary count.

**Response:**

```json
{
  "success": true,
  "message": "Menu visibility retrieved successfully",
  "data": {
    "items": [
      {
        "menu_key": "category.overview",
        "title": "Overview",
        "type": "title",
        "visible": true,
        "children": [
          {
            "menu_key": "dashboard",
            "title": "Dashboard",
            "type": "item",
            "path": "/dashboard",
            "visible": true
          }
        ]
      },
      {
        "menu_key": "category.hrPeople",
        "title": "HR & People",
        "type": "title",
        "visible": true,
        "children": [
          {
            "menu_key": "employeeManagement",
            "title": "Employee Management",
            "type": "collapse",
            "visible": true,
            "children": [
              { "menu_key": "employees.directory", "title": "Employee Directory", "type": "item", "path": "/employees/directory", "visible": true },
              { "menu_key": "employees.assets", "title": "Assets & Equipment", "type": "item", "path": "/employees/assets", "visible": false }
            ]
          }
        ]
      }
    ],
    "summary": {
      "total": 78,
      "visible": 72,
      "hidden": 6
    }
  }
}
```

---

### 2. Bulk Update Menu Visibility

**`PUT /api/v1/settings/menu-visibility`**

**Auth:** Admin / Super Admin

Update visibility for multiple items at once. The backend enforces **cascade rules** automatically.

**Request:**

```json
{
  "items": [
    { "menu_key": "employees.assets", "visible": false },
    { "menu_key": "payrollManagement", "visible": false },
    { "menu_key": "training.certifications", "visible": true }
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `items` | array | Yes | Array of visibility changes (min 1) |
| `items[].menu_key` | string | Yes | Menu item key |
| `items[].visible` | boolean | Yes | New visibility state |

**Response:**

```json
{
  "success": true,
  "message": "Menu visibility updated for 12 items",
  "data": {
    "updated_count": 12,
    "items": [
      { "menu_key": "employees.assets", "visible": false },
      { "menu_key": "payrollManagement", "visible": false },
      { "menu_key": "payroll.dashboard", "visible": false },
      { "menu_key": "payroll.run", "visible": false },
      { "menu_key": "payroll.salaryStructure", "visible": false },
      { "menu_key": "payroll.payslips", "visible": false },
      { "menu_key": "payroll.loans", "visible": false },
      { "menu_key": "payroll.compliance", "visible": false },
      { "menu_key": "payroll.reports", "visible": false },
      { "menu_key": "training.certifications", "visible": true },
      { "menu_key": "trainingDevelopment", "visible": true },
      { "menu_key": "category.talent", "visible": true }
    ]
  }
}
```

Note how hiding `payrollManagement` cascaded to all 7 children, and showing `training.certifications` auto-showed its parent group and category.

**Error — invalid keys:**

```json
{
  "success": false,
  "message": "invalid menu keys: foo.bar, nonexistent.item"
}
```

**Error — protected items:**

```json
{
  "success": false,
  "message": "protected menu items cannot be hidden: dashboard, settings.system"
}
```

---

### 3. Toggle Single Menu Item

**`PATCH /api/v1/settings/menu-visibility/:menuKey`**

**Auth:** Admin / Super Admin

Quick toggle for a single item. Cascade rules apply.

**URL Parameter:** `menuKey` — e.g., `performance.goals`, `payrollManagement`

**Request:**

```json
{
  "visible": false
}
```

**Response (hiding a leaf item):**

```json
{
  "success": true,
  "message": "Menu item 'performance.goals' is now hidden",
  "data": {
    "menu_key": "performance.goals",
    "visible": false,
    "cascaded": [
      { "menu_key": "performance.goals", "visible": false }
    ]
  }
}
```

**Response (hiding a parent — cascades to children):**

```json
{
  "success": true,
  "message": "Menu item 'payrollManagement' and 7 children is now hidden",
  "data": {
    "menu_key": "payrollManagement",
    "visible": false,
    "cascaded": [
      { "menu_key": "payrollManagement", "visible": false },
      { "menu_key": "payroll.dashboard", "visible": false },
      { "menu_key": "payroll.run", "visible": false },
      { "menu_key": "payroll.salaryStructure", "visible": false },
      { "menu_key": "payroll.payslips", "visible": false },
      { "menu_key": "payroll.loans", "visible": false },
      { "menu_key": "payroll.compliance", "visible": false },
      { "menu_key": "payroll.reports", "visible": false }
    ]
  }
}
```

---

### 4. Reset All to Default

**`POST /api/v1/settings/menu-visibility/reset`**

**Auth:** Admin / Super Admin

Resets all items to visible. No request body needed.

**Response:**

```json
{
  "success": true,
  "message": "Menu visibility reset to defaults. All 78 items are now visible.",
  "data": {
    "reset_count": 78
  }
}
```

---

### 5. Get Active (Hidden Keys) — Frontend Endpoint

**`GET /api/v1/settings/menu-visibility/active`**

**Auth:** Any authenticated user

This is the lightweight endpoint the **frontend calls on every page load**. Returns only the list of hidden menu keys.

**Response (some items hidden):**

```json
{
  "success": true,
  "message": "Hidden menu keys retrieved successfully",
  "data": {
    "hidden_keys": [
      "employees.assets",
      "payrollManagement",
      "payroll.dashboard",
      "payroll.run",
      "payroll.salaryStructure",
      "payroll.payslips",
      "payroll.loans",
      "payroll.compliance",
      "payroll.reports"
    ]
  }
}
```

**Response (nothing hidden):**

```json
{
  "success": true,
  "message": "Hidden menu keys retrieved successfully",
  "data": {
    "hidden_keys": []
  }
}
```

---

## Cascade Rules

| Rule | Trigger | Effect |
|------|---------|--------|
| **Rule 1** | Hide a group | All child items are hidden |
| **Rule 2** | Hide a category | All child groups and items are hidden |
| **Rule 3** | Show a child item | Parent group + grandparent category auto-shown |
| **Rule 4** | Hide a protected key | Request rejected with 400 error |

### Protected Keys (cannot be hidden)

| Key | Reason |
|-----|--------|
| `category.overview` | Core navigation |
| `dashboard` | Landing page |
| `category.administration` | Admin access |
| `settingsAdmin` | Admin access |
| `settings.system` | Prevents lockout |

---

## Role Access

| Endpoint | Employee | Manager | HR | Admin | Super Admin |
|----------|----------|---------|-----|-------|-------------|
| GET `/menu-visibility` | No | No | No | **Yes** | **Yes** |
| PUT `/menu-visibility` | No | No | No | **Yes** | **Yes** |
| PATCH `/menu-visibility/:menuKey` | No | No | No | **Yes** | **Yes** |
| POST `/menu-visibility/reset` | No | No | No | **Yes** | **Yes** |
| GET `/menu-visibility/active` | **Yes** | **Yes** | **Yes** | **Yes** | **Yes** |

---

## Frontend Integration

```typescript
// On app load (layout component or context provider)
const res = await fetch("/api/v1/settings/menu-visibility/active", {
  headers: { Authorization: `Bearer ${token}` }
});
const { data } = await res.json();
const hiddenKeys: string[] = data.hidden_keys;

// Filter navigation — remove hidden items at every level
const filteredNav = navigationConfig
  .filter(cat => !hiddenKeys.includes(cat.key))
  .map(cat => ({
    ...cat,
    subMenu: cat.subMenu
      .filter(group => !hiddenKeys.includes(group.key))
      .map(group => ({
        ...group,
        subMenu: group.subMenu.filter(item => !hiddenKeys.includes(item.key))
      }))
  }));
```

**Caching:** Use SWR/React Query with 5-minute TTL, invalidate after admin saves.

---

## All 78 Menu Keys

### Categories (8)
`category.overview`, `category.hrPeople`, `category.timeLeave`, `category.payroll`, `category.talent`, `category.mySpace`, `category.supportReports`, `category.administration`

### Groups (12)
`employeeManagement`, `timeAttendance`, `leaveManagement`, `payrollManagement`, `recruitmentManagement`, `performanceManagement`, `trainingDevelopment`, `projectManagement`, `selfService`, `helpdeskSupport`, `reportsAnalytics`, `settingsAdmin`

### Items (58)
`dashboard`, `employees.directory`, `employees.organization`, `employees.assets`, `employees.documents`, `employees.activeOnboarding`, `employees.offboarding`, `attendance.daily`, `attendance.shifts`, `attendance.timesheets`, `attendance.overtime`, `attendance.reports`, `leave.apply`, `leave.requests`, `leave.calendar`, `leave.policies`, `leave.holidays`, `leave.reports`, `payroll.dashboard`, `payroll.run`, `payroll.salaryStructure`, `payroll.payslips`, `payroll.loans`, `payroll.compliance`, `payroll.reports`, `recruitment.dashboard`, `recruitment.requisitions`, `recruitment.openings`, `recruitment.candidates`, `recruitment.interviews`, `recruitment.offers`, `recruitment.talentPool`, `performance.dashboard`, `performance.goals`, `performance.appraisals`, `performance.feedback360`, `performance.talentReview`, `performance.reports`, `training.dashboard`, `training.catalog`, `training.myLearning`, `training.sessions`, `training.certifications`, `training.reports`, `projects.allProjects`, `projects.myProjects`, `projects.dailyTasks`, `selfService.myProfile`, `selfService.myPayslips`, `selfService.requests`, `selfService.directory`, `selfService.myAssets`, `selfService.requestLeave`, `helpdesk.myTickets`, `helpdesk.knowledgeBase`, `helpdesk.dashboard`, `reports.analytics`, `reports.standard`, `reports.builder`, `reports.scheduled`, `settings.organization`, `settings.users`, `settings.policies`, `settings.integrations`, `settings.security`, `settings.system`

---

## Error Responses

| Status | Description |
|--------|-------------|
| 400 | Invalid keys, protected items, validation error |
| 401 | Not authenticated |
| 403 | Not Admin / Super Admin |
| 404 | Menu key not found |
| 500 | Internal server error |
