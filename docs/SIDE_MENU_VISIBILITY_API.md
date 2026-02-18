# Side Menu Visibility API Documentation

## Overview

The Side Menu Visibility feature allows **Admin** and **Super Admin** users to control which navigation menu items are visible in the application's side navigation. This enables organizations to customize their HRMS experience by hiding modules they don't use, reducing clutter and improving focus for all users.

### Key Concepts

| Concept | Description |
|---------|-------------|
| **Menu Item** | A single entry in the side navigation (category, group, or page link) |
| **Menu Key** | Unique identifier for each menu item (e.g., `dashboard`, `performance.goals`) |
| **Visibility** | Boolean flag — `true` (shown) or `false` (hidden) |
| **Cascade Rule** | Hiding a parent hides all children; showing a child auto-shows its parent |
| **Scope** | Changes apply **globally** to the entire organization (all users see the same menu structure) |

### Menu Hierarchy

```
Category (title)                   ← e.g., "HR & People", "Payroll & Finance"
  └── Group (collapse)             ← e.g., "Employee Management", "Leave Management"
        └── Item (page link)       ← e.g., "Employee Directory", "Leave Calendar"
```

### How It Works

1. **Admin** opens System Settings → Side Menu Visibility tab
2. Admin toggles individual items, groups, or entire categories ON/OFF
3. On save, the visibility configuration is persisted to the backend
4. When **any user** loads the app, the frontend fetches the visibility config and filters the navigation accordingly
5. Hidden items are also hidden from search and breadcrumbs

---

## Database

### Table: `menu_visibility_settings`

| Column | Type | Nullable | Default | Description |
|--------|------|----------|---------|-------------|
| `id` | `bigint` | No | Auto | Primary key |
| `menu_key` | `varchar(100)` | No | — | Unique menu item key (e.g., `performance.goals`) |
| `visible` | `boolean` | No | `true` | Whether the item is visible |
| `updated_by` | `bigint` | Yes | — | FK → `users.id` (who last changed it) |
| `updated_at` | `timestamp` | No | `now()` | Last modification time |
| `created_at` | `timestamp` | No | `now()` | Row creation time |

**Indexes:**
- `UNIQUE INDEX` on `menu_key`
- `INDEX` on `visible` (for quick filtering)

**Notes:**
- Only items that are explicitly **hidden** (`visible = false`) need to be stored. Items not in this table are assumed visible by default.
- Alternatively, store all items for a complete audit trail.

---

## API Endpoints

**Base URL:** `/api/v1/settings`

**Auth:** All endpoints require `AuthMiddleware()` + `AdminMiddleware()` (only Admin/Super Admin)

| # | Method | Endpoint | Description |
|---|--------|----------|-------------|
| 1 | GET | `/settings/menu-visibility` | Get all menu visibility settings |
| 2 | PUT | `/settings/menu-visibility` | Bulk update menu visibility |
| 3 | PATCH | `/settings/menu-visibility/:menuKey` | Toggle single menu item |
| 4 | POST | `/settings/menu-visibility/reset` | Reset all to default (all visible) |
| 5 | GET | `/settings/menu-visibility/active` | Get only visible menu keys (for frontend filtering) |

**Total: 5 endpoints**

---

### 1. Get All Menu Visibility Settings

**`GET /api/v1/settings/menu-visibility`**

**Auth:** Admin, Super Admin

Returns the full list of menu items with their current visibility status.

**Response:**

```json
{
  "success": true,
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
              {
                "menu_key": "employees.directory",
                "title": "Employee Directory",
                "type": "item",
                "path": "/employees/directory",
                "visible": true
              },
              {
                "menu_key": "employees.organization",
                "title": "Organization Structure",
                "type": "item",
                "path": "/employees/organization",
                "visible": true
              },
              {
                "menu_key": "employees.assets",
                "title": "Assets & Equipment",
                "type": "item",
                "path": "/employees/assets",
                "visible": false
              }
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

**Auth:** Admin, Super Admin

Update visibility for multiple menu items at once. This is the primary endpoint used by the UI when the admin clicks "Save".

**Request:**

```json
{
  "items": [
    { "menu_key": "employees.assets", "visible": false },
    { "menu_key": "employees.documents", "visible": false },
    { "menu_key": "payrollManagement", "visible": false },
    { "menu_key": "category.payroll", "visible": false },
    { "menu_key": "recruitment.talentPool", "visible": false },
    { "menu_key": "training.certifications", "visible": true }
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `items` | `array` | Yes | Array of visibility changes |
| `items[].menu_key` | `string` | Yes | Menu item key |
| `items[].visible` | `boolean` | Yes | New visibility state |

**Response:**

```json
{
  "success": true,
  "message": "Menu visibility updated for 6 items",
  "data": {
    "updated_count": 6,
    "items": [
      { "menu_key": "employees.assets", "visible": false },
      { "menu_key": "employees.documents", "visible": false },
      { "menu_key": "payrollManagement", "visible": false },
      { "menu_key": "category.payroll", "visible": false },
      { "menu_key": "recruitment.talentPool", "visible": false },
      { "menu_key": "training.certifications", "visible": true }
    ]
  }
}
```

**Server-side cascade logic:**

The backend MUST enforce cascade rules:

```
When hiding a parent:
  → All children must also be set to hidden

When showing a child:
  → Parent must also be set to visible
  → Grandparent (if any) must also be set to visible
```

**Validation:**
- All `menu_key` values must match known menu item keys
- Unknown keys return `400` with details of invalid keys
- At least one item must be provided

**Error Response (invalid keys):**

```json
{
  "success": false,
  "message": "Invalid menu keys found",
  "errors": {
    "invalid_keys": ["foo.bar", "nonexistent.item"]
  }
}
```

---

### 3. Toggle Single Menu Item

**`PATCH /api/v1/settings/menu-visibility/:menuKey`**

**Auth:** Admin, Super Admin

Quick toggle for a single menu item. Useful for API integrations.

**URL Parameters:**

| Param | Type | Description |
|-------|------|-------------|
| `menuKey` | `string` | Menu item key (e.g., `performance.goals`) |

**Request:**

```json
{
  "visible": false
}
```

**Response:**

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

If hiding a parent (collapse/title), cascaded changes are included:

```json
{
  "success": true,
  "message": "Menu item 'payrollManagement' and 7 children are now hidden",
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

**Auth:** Admin, Super Admin

Resets all menu items to visible. Deletes all entries from the `menu_visibility_settings` table (or sets all to `visible = true`).

**Request:** No body required.

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

### 5. Get Active (Visible) Menu Keys

**`GET /api/v1/settings/menu-visibility/active`**

**Auth:** Any authenticated user (this is the endpoint the **frontend calls on app load**)

Returns only the list of **hidden** menu keys. The frontend uses this to filter out items from the navigation. This endpoint is lightweight and fast — called on every page load.

**Response (when some items are hidden):**

```json
{
  "success": true,
  "data": {
    "hidden_keys": [
      "employees.assets",
      "employees.documents",
      "payrollManagement",
      "payroll.dashboard",
      "payroll.run",
      "payroll.salaryStructure",
      "payroll.payslips",
      "payroll.loans",
      "payroll.compliance",
      "payroll.reports",
      "category.payroll",
      "recruitment.talentPool"
    ]
  }
}
```

**Response (when all items are visible):**

```json
{
  "success": true,
  "data": {
    "hidden_keys": []
  }
}
```

**Why return hidden_keys instead of visible_keys?**
- The hidden list is typically much smaller (most orgs hide only a few items)
- Smaller payload = faster load
- Frontend defaults everything to visible, then removes hidden items

---

## All Menu Item Keys (Reference)

The following keys are recognized by the system. Use these in API requests.

### Categories (type: `title`)

| Key | Title |
|-----|-------|
| `category.overview` | Overview |
| `category.hrPeople` | HR & People |
| `category.timeLeave` | Time & Leave |
| `category.payroll` | Payroll & Finance |
| `category.talent` | Talent Management |
| `category.mySpace` | My Space |
| `category.supportReports` | Support & Reports |
| `category.administration` | Administration |

### Groups (type: `collapse`)

| Key | Title | Parent Category |
|-----|-------|----------------|
| `employeeManagement` | Employee Management | HR & People |
| `timeAttendance` | Time & Attendance | Time & Leave |
| `leaveManagement` | Leave Management | Time & Leave |
| `payrollManagement` | Payroll | Payroll & Finance |
| `recruitmentManagement` | Recruitment | Talent Management |
| `performanceManagement` | Manage Performance | Talent Management |
| `trainingDevelopment` | Training & Development | Talent Management |
| `projectManagement` | Manage Projects | Talent Management |
| `selfService` | Self Service | My Space |
| `helpdeskSupport` | Helpdesk & Support | Support & Reports |
| `reportsAnalytics` | Reports & Analytics | Support & Reports |
| `settingsAdmin` | Settings & Admin | Administration |

### Items (type: `item`)

| Key | Title | Parent Group |
|-----|-------|-------------|
| `dashboard` | Dashboard | Overview |
| `employees.directory` | Employee Directory | Employee Management |
| `employees.organization` | Organization Structure | Employee Management |
| `employees.assets` | Assets & Equipment | Employee Management |
| `employees.documents` | Documents | Employee Management |
| `employees.activeOnboarding` | Active Onboarding | Employee Management |
| `employees.offboarding` | Offboarding | Employee Management |
| `attendance.daily` | Daily Attendance | Time & Attendance |
| `attendance.shifts` | Shift Management | Time & Attendance |
| `attendance.timesheets` | Timesheets | Time & Attendance |
| `attendance.overtime` | Overtime | Time & Attendance |
| `attendance.reports` | Attendance Reports | Time & Attendance |
| `leave.apply` | Apply Leave | Leave Management |
| `leave.requests` | Leave Requests | Leave Management |
| `leave.calendar` | Leave Calendar | Leave Management |
| `leave.policies` | Leave Policies | Leave Management |
| `leave.holidays` | Holidays | Leave Management |
| `leave.reports` | Leave Reports | Leave Management |
| `payroll.dashboard` | Dashboard | Payroll |
| `payroll.run` | Run Payroll | Payroll |
| `payroll.salaryStructure` | Salary Structure | Payroll |
| `payroll.payslips` | Payslips | Payroll |
| `payroll.loans` | Loans & Advances | Payroll |
| `payroll.compliance` | Tax & Compliance | Payroll |
| `payroll.reports` | Payroll Reports | Payroll |
| `recruitment.dashboard` | Dashboard | Recruitment |
| `recruitment.requisitions` | Requisitions | Recruitment |
| `recruitment.openings` | Job Openings | Recruitment |
| `recruitment.candidates` | Candidates | Recruitment |
| `recruitment.interviews` | Interviews | Recruitment |
| `recruitment.offers` | Offers | Recruitment |
| `recruitment.talentPool` | Talent Pool | Recruitment |
| `performance.dashboard` | Dashboard | Manage Performance |
| `performance.goals` | Goals & OKRs | Manage Performance |
| `performance.appraisals` | Appraisals | Manage Performance |
| `performance.feedback360` | 360° Feedback | Manage Performance |
| `performance.talentReview` | Talent Review | Manage Performance |
| `performance.reports` | Reports | Manage Performance |
| `training.dashboard` | Learning Dashboard | Training & Development |
| `training.catalog` | Course Catalog | Training & Development |
| `training.myLearning` | My Learning | Training & Development |
| `training.sessions` | Training Sessions | Training & Development |
| `training.certifications` | Certifications | Training & Development |
| `training.reports` | Training Reports | Training & Development |
| `projects.allProjects` | All Projects | Manage Projects |
| `projects.myProjects` | My Projects | Manage Projects |
| `projects.dailyTasks` | Daily Tasks | Manage Projects |
| `selfService.myProfile` | My Profile | Self Service |
| `selfService.myPayslips` | My Payslips | Self Service |
| `selfService.requests` | My Requests | Self Service |
| `selfService.directory` | People Directory | Self Service |
| `selfService.myAssets` | My Assets | Self Service |
| `selfService.requestLeave` | Request Leave | Self Service |
| `helpdesk.myTickets` | My Tickets | Helpdesk & Support |
| `helpdesk.knowledgeBase` | Knowledge Base | Helpdesk & Support |
| `helpdesk.dashboard` | Helpdesk Dashboard | Helpdesk & Support |
| `reports.analytics` | Analytics Dashboard | Reports & Analytics |
| `reports.standard` | Standard Reports | Reports & Analytics |
| `reports.builder` | Report Builder | Reports & Analytics |
| `reports.scheduled` | Scheduled Reports | Reports & Analytics |
| `settings.organization` | Organization | Settings & Admin |
| `settings.users` | Users & Roles | Settings & Admin |
| `settings.policies` | Policies | Settings & Admin |
| `settings.integrations` | Integrations | Settings & Admin |
| `settings.security` | Security | Settings & Admin |
| `settings.system` | System Settings | Settings & Admin |

---

## Cascade Rules (Important)

### Rule 1: Hiding a Parent Hides All Children

```
Admin hides "payrollManagement" (group)
  → Server automatically hides:
    - payroll.dashboard
    - payroll.run
    - payroll.salaryStructure
    - payroll.payslips
    - payroll.loans
    - payroll.compliance
    - payroll.reports
```

### Rule 2: Hiding a Category Hides All Groups and Items

```
Admin hides "category.talent" (category)
  → Server automatically hides:
    - recruitmentManagement + all 7 children
    - performanceManagement + all 6 children
    - trainingDevelopment + all 6 children
    - projectManagement + all 3 children
```

### Rule 3: Showing a Child Auto-Shows Its Parents

```
Admin shows "payroll.payslips" (item)
  → Server automatically shows:
    - payrollManagement (parent group)
    - category.payroll (grandparent category)
```

### Rule 4: Protected Items Cannot Be Hidden

The following items MUST always remain visible and cannot be hidden:

| Key | Reason |
|-----|--------|
| `category.overview` | Core navigation |
| `dashboard` | Every user needs a landing page |
| `category.administration` | Admins must always access settings |
| `settingsAdmin` | Admins must always access settings |
| `settings.system` | Prevents admin from locking themselves out |

If an API request attempts to hide a protected item, return:

```json
{
  "success": false,
  "message": "Protected menu items cannot be hidden",
  "errors": {
    "protected_keys": ["dashboard", "settings.system"]
  }
}
```

---

## Frontend Integration

### On App Load

The frontend should call `GET /api/v1/settings/menu-visibility/active` on app load (or cache it in session storage):

```typescript
// In a layout component or context provider
const response = await fetch("/api/v1/settings/menu-visibility/active");
const { data } = await response.json();
const hiddenKeys: string[] = data.hidden_keys;

// Filter navigation config
const filteredNav = navigationConfig.filter(
  (item) => !hiddenKeys.includes(item.key)
).map((category) => ({
  ...category,
  subMenu: category.subMenu.filter(
    (group) => !hiddenKeys.includes(group.key)
  ).map((group) => ({
    ...group,
    subMenu: group.subMenu.filter(
      (item) => !hiddenKeys.includes(item.key)
    ),
  })),
}));
```

### Caching Strategy

| Strategy | TTL | When to Invalidate |
|----------|-----|-------------------|
| Session Storage | Per session | On logout |
| React Context | Per page load | Never (refetches on load) |
| SWR / React Query | 5 minutes | On mutation (save visibility) |

---

## Error Responses

All endpoints follow the standard format:

```json
{
  "success": false,
  "message": "Error description",
  "errors": "optional detail"
}
```

| Status | Description |
|--------|-------------|
| 400 | Invalid request / unknown menu keys / protected items |
| 401 | Not authenticated |
| 403 | Not an Admin or Super Admin |
| 404 | Menu key not found (PATCH endpoint) |
| 500 | Internal server error |

---

## Role Access

| Endpoint | Employee | Manager | HR | Admin | Super Admin |
|----------|----------|---------|-----|-------|-------------|
| GET `/menu-visibility` | No | No | No | Yes | Yes |
| PUT `/menu-visibility` | No | No | No | Yes | Yes |
| PATCH `/menu-visibility/:key` | No | No | No | Yes | Yes |
| POST `/menu-visibility/reset` | No | No | No | Yes | Yes |
| GET `/menu-visibility/active` | **Yes** | **Yes** | **Yes** | **Yes** | **Yes** |

The `active` endpoint is the only one accessible by all authenticated users — it just returns the list of hidden keys for filtering.

---

## Complete Endpoint Summary

| # | Method | Endpoint | Auth | Description |
|---|--------|----------|------|-------------|
| 1 | GET | `/api/v1/settings/menu-visibility` | Admin | Full visibility tree |
| 2 | PUT | `/api/v1/settings/menu-visibility` | Admin | Bulk update visibility |
| 3 | PATCH | `/api/v1/settings/menu-visibility/:menuKey` | Admin | Toggle single item |
| 4 | POST | `/api/v1/settings/menu-visibility/reset` | Admin | Reset all to visible |
| 5 | GET | `/api/v1/settings/menu-visibility/active` | Any | Get hidden keys (frontend) |

**Total: 5 endpoints**
