# Feature Access Control API

A comprehensive role-based feature access control system for the HRMS application. This API allows the frontend to dynamically control which menus, pages, and actions are visible to each user based on their assigned role(s).

---

## Overview

### How It Works

1. **Features** are system modules (e.g., "Employee Management", "Payroll", "Leave Management")
2. Each feature has **permissions** (e.g., `employee:read`, `employee:create`, `employee:delete`)
3. **Roles** (e.g., Admin, HR, Employee) are assigned a set of permissions
4. Users are assigned one or more roles
5. The frontend calls `GET /api/v1/features/my-access` on login to determine what to show/hide

### Default Roles

| Role | Code | Description |
|------|------|-------------|
| Administrator | `admin` | Full system access (all permissions) |
| HR Manager | `hr` | HR management, employees, leave, payroll, organization |
| IT Support | `it` | Helpdesk, assets, biometric configuration |
| Manager | `manager` | Team oversight, leave/attendance approval |
| Employee | `employee` | Self-service features only |

### Feature Categories

| Category | Features |
|----------|----------|
| Core | Dashboard |
| HR | Employees, Leave, Attendance, Shifts & Rosters |
| Organization | Departments, Teams, Positions, Organizations, Org Units, Locations, Cost Centers |
| Finance | Payroll |
| Operations | Assets |
| Support | Helpdesk |
| Administration | Users, Audit & Security, Reports, Settings |

---

## API Endpoints

### 1. Get My Feature Access

Returns the feature access map for the currently authenticated user. **Call this on login** to determine what to show in the UI.

```
GET /api/v1/features/my-access
Authorization: Bearer <token>
```

**Response:**

```json
{
  "success": true,
  "message": "Feature access retrieved successfully",
  "data": {
    "user_id": 1,
    "user_role": "hr",
    "accessible_features": [
      "dashboard",
      "employees",
      "leave",
      "departments",
      "payroll"
    ],
    "permissions": {
      "dashboard:view": true,
      "dashboard:manage_announcements": true,
      "employee:read": true,
      "employee:create": true,
      "employee:update": true,
      "employee:delete": false,
      "leave:read": true,
      "leave:approve": true,
      "payroll:read": true,
      "payroll:run": true,
      "audit:read": false,
      "settings:update": false
    },
    "features": [
      {
        "code": "dashboard",
        "name": "Dashboard",
        "description": "View dashboard statistics, announcements, and quick actions",
        "icon": "dashboard",
        "category": "core",
        "has_access": true,
        "permissions": [
          {
            "code": "dashboard:view",
            "name": "View Dashboard",
            "description": "Access the main dashboard",
            "granted": true
          },
          {
            "code": "dashboard:manage_announcements",
            "name": "Manage Announcements",
            "description": "Create, edit, and delete company announcements",
            "granted": true
          }
        ]
      },
      {
        "code": "employees",
        "name": "Employee Management",
        "description": "Manage employee records, onboarding, and offboarding",
        "icon": "people",
        "category": "hr",
        "has_access": true,
        "permissions": [
          {
            "code": "employee:read",
            "name": "View Employees",
            "description": "View employee profiles and records",
            "granted": true
          },
          {
            "code": "employee:create",
            "name": "Create Employees",
            "description": "Add new employee records",
            "granted": true
          },
          {
            "code": "employee:delete",
            "name": "Delete Employees",
            "description": "Remove employee records",
            "granted": false
          }
        ]
      }
    ]
  }
}
```

**Frontend Usage:**

```javascript
// On login or app initialization
const response = await fetch('/api/v1/features/my-access', {
  headers: { Authorization: `Bearer ${token}` }
});
const { data } = await response.json();

// Store in global state (e.g., Redux, Zustand, Context)
store.setPermissions(data.permissions);
store.setAccessibleFeatures(data.accessible_features);

// Check access in components
const canViewEmployees = data.permissions['employee:read'];
const canDeleteEmployee = data.permissions['employee:delete'];
const canRunPayroll = data.permissions['payroll:run'];

// Show/hide menu items
const showEmployeeMenu = data.accessible_features.includes('employees');
const showPayrollMenu = data.accessible_features.includes('payroll');
const showAuditMenu = data.accessible_features.includes('audit');
```

---

### 2. Get All Features (Feature Catalog)

Returns the complete list of all system features and their permissions. Used by admins when configuring role permissions.

```
GET /api/v1/features
Authorization: Bearer <token>
```

**Requires:** HR or Admin role

**Response:**

```json
{
  "success": true,
  "message": "Feature catalog retrieved successfully",
  "data": {
    "total": 18,
    "features": [
      {
        "code": "dashboard",
        "name": "Dashboard",
        "description": "View dashboard statistics, announcements, and quick actions",
        "icon": "dashboard",
        "category": "core",
        "permissions": [
          {
            "code": "dashboard:view",
            "name": "View Dashboard",
            "description": "Access the main dashboard"
          },
          {
            "code": "dashboard:manage_announcements",
            "name": "Manage Announcements",
            "description": "Create, edit, and delete company announcements"
          }
        ]
      }
    ],
    "categories": {
      "core": [...],
      "hr": [...],
      "organization": [...],
      "finance": [...],
      "operations": [...],
      "support": [...],
      "administration": [...]
    }
  }
}
```

---

### 3. Get Role Feature Access

Returns the feature access map for a specific role. Used by admins to see what a role can do.

```
GET /api/v1/features/roles/:id
Authorization: Bearer <token>
```

**Requires:** HR or Admin role

**Parameters:**

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| id | integer | path | Role ID |

**Response:**

```json
{
  "success": true,
  "message": "Role feature access retrieved successfully",
  "data": {
    "role": {
      "id": 2,
      "code": "hr",
      "name": "HR Manager",
      "description": "HR management with employee, leave, attendance, and organizational access"
    },
    "accessible_features": [
      "dashboard",
      "employees",
      "leave",
      "attendance",
      "shifts",
      "departments",
      "teams",
      "positions",
      "payroll",
      "assets",
      "helpdesk"
    ],
    "permissions": {
      "employee:read": true,
      "employee:create": true,
      "employee:update": true,
      "employee:delete": true,
      "audit:read": false,
      "settings:update": false
    },
    "features": [...]
  }
}
```

---

### 4. Update Role Permissions

Updates the permissions assigned to a role. Replaces all existing permissions with the provided list.

```
PUT /api/v1/features/roles/:id
Authorization: Bearer <token>
Content-Type: application/json
```

**Requires:** Admin role only

**Parameters:**

| Parameter | Type | Location | Description |
|-----------|------|----------|-------------|
| id | integer | path | Role ID |

**Request Body:**

```json
{
  "permissions": [
    "dashboard:view",
    "employee:read",
    "employee:create",
    "employee:update",
    "leave:read",
    "leave:create",
    "leave:approve",
    "department:read",
    "helpdesk:read",
    "helpdesk:create"
  ]
}
```

**Response:**

```json
{
  "success": true,
  "message": "Role permissions updated successfully",
  "data": {
    "permissions_count": 10,
    "features": [...]
  }
}
```

---

## Complete Permission Reference

### Dashboard
| Permission | Description |
|------------|-------------|
| `dashboard:view` | Access the main dashboard |
| `dashboard:manage_announcements` | Create, edit, and delete company announcements |

### Employee Management
| Permission | Description |
|------------|-------------|
| `employee:read` | View employee profiles and records |
| `employee:create` | Add new employee records |
| `employee:update` | Edit employee information |
| `employee:delete` | Remove employee records |
| `employee:onboard` | Run the employee onboarding process |
| `employee:offboard` | Run the employee offboarding process |

### Leave Management
| Permission | Description |
|------------|-------------|
| `leave:read` | View leave requests and balances |
| `leave:create` | Submit leave applications |
| `leave:approve` | Approve or reject leave requests |
| `leave:manage_types` | Create and configure leave types |
| `leave:manage_policies` | Create and configure leave policies |
| `leave:manage_holidays` | Create and manage public holidays |

### Attendance & Biometric
| Permission | Description |
|------------|-------------|
| `attendance:read` | View attendance records and reports |
| `attendance:approve` | Approve attendance corrections |
| `attendance:manage_config` | Configure biometric devices and settings |

### Shifts & Rosters
| Permission | Description |
|------------|-------------|
| `shift:read` | View shifts and roster assignments |
| `shift:create` | Create new shift definitions |
| `shift:update` | Modify shift details |
| `shift:delete` | Remove shift definitions |
| `shift:manage_rosters` | Create and manage roster assignments |
| `shift:approve_swaps` | Approve or reject shift swap requests |

### Departments
| Permission | Description |
|------------|-------------|
| `department:read` | View department information |
| `department:create` | Add new departments |
| `department:update` | Edit department details |
| `department:delete` | Remove departments |

### Teams
| Permission | Description |
|------------|-------------|
| `team:read` | View team information |
| `team:create` | Add new teams |
| `team:update` | Edit team details |
| `team:delete` | Remove teams |

### Job Positions
| Permission | Description |
|------------|-------------|
| `position:read` | View job positions |
| `position:create` | Add new job positions |
| `position:update` | Edit job position details |
| `position:delete` | Remove job positions |

### Organizations
| Permission | Description |
|------------|-------------|
| `organization:read` | View organization information |
| `organization:create` | Add new organizations |
| `organization:update` | Edit organization details |
| `organization:delete` | Remove organizations |

### Organization Units
| Permission | Description |
|------------|-------------|
| `organization_unit:read` | View organization units |
| `organization_unit:create` | Add new organization units |
| `organization_unit:update` | Edit organization unit details |
| `organization_unit:delete` | Remove organization units |

### Locations
| Permission | Description |
|------------|-------------|
| `location:read` | View location information |
| `location:create` | Add new locations |
| `location:update` | Edit location details |
| `location:delete` | Remove locations |

### Cost Centers
| Permission | Description |
|------------|-------------|
| `cost_center:read` | View cost center information |
| `cost_center:create` | Add new cost centers |
| `cost_center:update` | Edit cost center details |
| `cost_center:delete` | Remove cost centers |

### Payroll
| Permission | Description |
|------------|-------------|
| `payroll:read` | View payroll runs, payslips, and reports |
| `payroll:run` | Execute payroll calculations and processing |
| `payroll:manage_structures` | Create and edit salary structures and components |
| `payroll:manage_loans` | Approve and manage employee loans |
| `payroll:manage_compliance` | Configure tax slabs and statutory rules |
| `payroll:manage_reports` | Generate and export payroll reports |

### Asset Management
| Permission | Description |
|------------|-------------|
| `asset:read` | View asset inventory |
| `asset:create` | Add new assets to inventory |
| `asset:update` | Edit asset information |
| `asset:delete` | Remove assets from inventory |
| `asset:assign` | Assign and reassign assets to employees |

### Helpdesk
| Permission | Description |
|------------|-------------|
| `helpdesk:read` | View helpdesk tickets |
| `helpdesk:create` | Submit support tickets |
| `helpdesk:manage` | Assign, resolve, and close tickets |
| `helpdesk:manage_kb` | Create and manage knowledge base articles |

### User Management
| Permission | Description |
|------------|-------------|
| `user:read` | View user accounts |
| `user:create` | Add new user accounts |
| `user:update` | Edit user account details |
| `user:delete` | Remove user accounts |
| `user:manage_roles` | Assign and revoke roles |

### Audit & Security
| Permission | Description |
|------------|-------------|
| `audit:read` | View system audit trails |
| `audit:export` | Export audit log data |

### Reports
| Permission | Description |
|------------|-------------|
| `report:view` | View system reports and analytics |
| `report:export` | Export report data |

### Settings
| Permission | Description |
|------------|-------------|
| `settings:read` | View system settings |
| `settings:update` | Modify system configuration |

---

## Frontend Integration Guide

### 1. Fetch Permissions on Login

```javascript
// After successful login
async function initializePermissions(token) {
  const res = await fetch('/api/v1/features/my-access', {
    headers: { Authorization: `Bearer ${token}` }
  });
  const { data } = await res.json();
  
  // Store globally
  localStorage.setItem('permissions', JSON.stringify(data.permissions));
  localStorage.setItem('features', JSON.stringify(data.accessible_features));
}
```

### 2. Permission Check Helper

```javascript
// utils/permissions.js
export function hasPermission(permissionCode) {
  const permissions = JSON.parse(localStorage.getItem('permissions') || '{}');
  return permissions[permissionCode] === true;
}

export function hasFeatureAccess(featureCode) {
  const features = JSON.parse(localStorage.getItem('features') || '[]');
  return features.includes(featureCode);
}

export function canPerformAny(...permissionCodes) {
  return permissionCodes.some(code => hasPermission(code));
}
```

### 3. Route Guard (React Example)

```jsx
// components/ProtectedRoute.jsx
function ProtectedRoute({ permission, feature, children }) {
  if (permission && !hasPermission(permission)) {
    return <AccessDenied />;
  }
  if (feature && !hasFeatureAccess(feature)) {
    return <AccessDenied />;
  }
  return children;
}

// Usage
<ProtectedRoute feature="payroll">
  <PayrollPage />
</ProtectedRoute>

<ProtectedRoute permission="employee:delete">
  <DeleteButton />
</ProtectedRoute>
```

### 4. Menu Configuration

```javascript
const menuItems = [
  { label: 'Dashboard', icon: 'dashboard', path: '/dashboard', feature: 'dashboard' },
  { label: 'Employees', icon: 'people', path: '/employees', feature: 'employees' },
  { label: 'Leave', icon: 'event_busy', path: '/leave', feature: 'leave' },
  { label: 'Payroll', icon: 'payments', path: '/payroll', feature: 'payroll' },
  { label: 'Helpdesk', icon: 'support_agent', path: '/helpdesk', feature: 'helpdesk' },
  { label: 'Settings', icon: 'settings', path: '/settings', feature: 'settings' },
];

// Filter based on access
const visibleMenuItems = menuItems.filter(item => hasFeatureAccess(item.feature));
```

### 5. Conditional UI Elements

```jsx
// Show/hide buttons based on permissions
{hasPermission('employee:create') && (
  <Button onClick={addEmployee}>Add Employee</Button>
)}

{hasPermission('employee:delete') && (
  <Button onClick={deleteEmployee} color="danger">Delete</Button>
)}

{hasPermission('leave:approve') && (
  <Button onClick={approveLeave}>Approve</Button>
)}
```

---

## Admin: Managing Role Permissions

### View a Role's Permissions

```
GET /api/v1/features/roles/2
```

### Update a Role's Permissions

```
PUT /api/v1/features/roles/2
{
  "permissions": [
    "dashboard:view",
    "employee:read",
    "employee:create",
    "employee:update",
    "leave:read",
    "leave:create",
    "leave:approve"
  ]
}
```

### Create a Custom Role

1. First, create the role via `POST /api/v1/roles` (if not already available)
2. Then, assign permissions via `PUT /api/v1/features/roles/:id`

---

## Error Responses

| Status | Description |
|--------|-------------|
| 401 | Authentication required (no/invalid token) |
| 403 | Insufficient permissions (role doesn't have access) |
| 404 | Role not found |
| 422 | Validation error (missing/invalid permissions) |
| 500 | Internal server error |
