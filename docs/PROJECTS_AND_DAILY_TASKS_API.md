# Projects & Daily Tasks API Documentation

## Overview

The Projects & Daily Tasks module enables project management and daily work logging. Managers/HODs create projects, assign employees, and track progress. Employees log daily tasks linked to projects, and the system automatically calculates project completion percentage.

### Key Features
- **Projects** created by Managers or Head of Department
- **Multi-member** assignment (projects can have many employees)
- **Cross-department** support (any department category)
- **Daily Tasks** logged by employees against projects
- **Auto-progress** calculation from completed task contributions
- **Task categories** for work classification

---

## Authentication

All endpoints require a valid JWT token in the `Authorization` header:
```
Authorization: Bearer <token>
```

---

## 1. Projects API (HR/Admin — Management)

Base URL: `/api/v1/projects`

### 1.1 Create Project

```
POST /api/v1/projects
```

**Request Body:**
```json
{
  "name": "ERP System Migration",
  "description": "Migrate legacy ERP to new platform",
  "priority": "high",
  "department_id": 3,
  "category": "IT",
  "start_date": "2026-02-10T00:00:00Z",
  "due_date": "2026-06-30T00:00:00Z",
  "estimated_hours": 500,
  "budget": 50000,
  "notes": "Q2 priority project",
  "member_ids": ["EMP001", "EMP002", "EMP005"]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Project name (2-200 chars) |
| `description` | string | No | Project description |
| `priority` | string | No | `low`, `medium`, `high`, `critical` (default: `medium`) |
| `department_id` | uint | No | Link to a department |
| `category` | string | No | Custom category (e.g., "IT", "Marketing", "Finance") |
| `start_date` | datetime | No | Planned start date |
| `end_date` | datetime | No | Planned end date |
| `due_date` | datetime | No | Due/deadline date |
| `estimated_hours` | float | No | Estimated total hours |
| `budget` | float | No | Project budget |
| `notes` | string | No | Additional notes |
| `member_ids` | string[] | No | Employee IDs to assign initially |

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Project created successfully",
  "data": {
    "id": 1,
    "project_code": "PRJ-0001",
    "name": "ERP System Migration",
    "status": "planning",
    "priority": "high",
    "progress": 0,
    "members": [
      { "id": 1, "employee_id": "EMP001", "employee_name": "John Doe", "role": "member" },
      { "id": 2, "employee_id": "EMP002", "employee_name": "Jane Smith", "role": "member" }
    ]
  }
}
```

### 1.2 List Projects

```
GET /api/v1/projects?page=1&page_size=20&status=active&priority=high&department_id=3&category=IT&search=ERP
```

| Query Param | Type | Description |
|------------|------|-------------|
| `page` | int | Page number (default: 1) |
| `page_size` | int | Items per page (default: 20, max: 100) |
| `status` | string | Filter: `planning`, `active`, `on_hold`, `completed`, `cancelled` |
| `priority` | string | Filter: `low`, `medium`, `high`, `critical` |
| `department_id` | uint | Filter by department |
| `category` | string | Filter by category |
| `search` | string | Search by name, code, or description |

**Response:** `200 OK` with pagination meta.

### 1.3 Get Project Details

```
GET /api/v1/projects/:id
```

### 1.4 Update Project

```
PUT /api/v1/projects/:id
```

**Request Body** (all fields optional):
```json
{
  "name": "Updated Name",
  "status": "active",
  "priority": "critical",
  "progress": 25.5,
  "due_date": "2026-07-31T00:00:00Z"
}
```

| Field | Type | Description |
|-------|------|-------------|
| `status` | string | Change project status. Setting `completed` auto-sets progress to 100% |
| `progress` | float | Manual progress override (0-100) |

### 1.5 Delete Project

```
DELETE /api/v1/projects/:id
```

### 1.6 Get Project Statistics

```
GET /api/v1/projects/statistics
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_active": 12,
    "by_status": {
      "planning": 3,
      "active": 9,
      "completed": 15,
      "on_hold": 2
    }
  }
}
```

### 1.7 Get Project Progress

```
GET /api/v1/projects/:id/progress
```

**Response:**
```json
{
  "success": true,
  "data": {
    "project_id": 1,
    "project_code": "PRJ-0001",
    "project_name": "ERP System Migration",
    "status": "active",
    "progress": 35.5,
    "estimated_hours": 500,
    "actual_hours": 178.5,
    "hours_utilization": 35.7,
    "member_count": 5,
    "task_summary": {
      "total_tasks": 42,
      "total_hours": 178.5,
      "completed_tasks": 35,
      "pending_tasks": 3,
      "in_progress_tasks": 4,
      "blocked_tasks": 0,
      "by_category": [
        { "category": "development", "task_count": 25, "total_hours": 120 },
        { "category": "testing", "task_count": 10, "total_hours": 40 }
      ]
    }
  }
}
```

---

## 2. Project Members API (HR/Admin)

### 2.1 Add Members to Project

```
POST /api/v1/projects/:id/members
```

**Request Body:**
```json
{
  "members": [
    { "employee_id": "EMP003", "role": "lead" },
    { "employee_id": "EMP004", "role": "member" },
    { "employee_id": "EMP006" }
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `members[].employee_id` | string | Yes | Employee ID |
| `members[].role` | string | No | `lead`, `member`, `reviewer` (default: `member`) |

### 2.2 Remove Member

```
POST /api/v1/projects/:id/members/remove
```

**Request Body:**
```json
{
  "employee_id": "EMP003"
}
```

### 2.3 List Project Members

```
GET /api/v1/projects/:id/members
```

### 2.4 Assign Project to Employee(s)

A convenient endpoint to assign a project to one or more employees in a single call.

```
POST /api/v1/projects/assign
```

**Request Body:**
```json
{
  "project_id": 1,
  "employee_ids": ["EMP001", "EMP002", "EMP005"],
  "role": "member"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `project_id` | uint | Yes | Project to assign |
| `employee_ids` | string[] | Yes | One or more employee IDs to assign |
| `role` | string | No | Role for all assigned employees: `lead`, `member`, `reviewer` (default: `member`) |

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Project assigned to employees successfully",
  "data": {
    "project": {
      "id": 1,
      "project_code": "PRJ-0001",
      "name": "ERP System Migration",
      "members": [...]
    },
    "members_added": [
      { "id": 10, "employee_id": "EMP001", "employee_name": "John Doe", "role": "member" },
      { "id": 11, "employee_id": "EMP002", "employee_name": "Jane Smith", "role": "member" }
    ]
  }
}
```

> **Note:** Employees already assigned to the project are silently skipped (no duplicates).

### 2.5 Unassign Employee from Project

```
POST /api/v1/projects/unassign
```

**Request Body:**
```json
{
  "project_id": 1,
  "employee_id": "EMP003"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `project_id` | uint | Yes | Project ID |
| `employee_id` | string | Yes | Employee to remove from the project |

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Employee unassigned from project successfully"
}
```

---

## 3. Self-Service: My Project Details

### 3.1 Get Assigned Project Details

Returns project details and progress for a project the employee is assigned to.

```
GET /api/v1/self-service/projects/:id
```

**Response:** `200 OK`
```json
{
  "success": true,
  "message": "Project details retrieved successfully",
  "data": {
    "project": {
      "id": 1,
      "project_code": "PRJ-0001",
      "name": "ERP System Migration",
      "status": "active",
      "progress": 45.0,
      "members": [...]
    },
    "progress": {
      "progress": 45.0,
      "estimated_hours": 500,
      "actual_hours": 225,
      "hours_utilization": 45.0,
      "member_count": 5,
      "task_summary": { ... }
    }
  }
}
```

> Returns `403 Forbidden` if the employee is not assigned to the project.

---

## 4. Daily Tasks API — Admin/HR View

Base URL: `/api/v1/daily-tasks`

### 3.1 List All Daily Tasks (HR/Admin)

```
GET /api/v1/daily-tasks?page=1&page_size=20&project_id=1&employee_id=EMP001&status=completed&category=development&date_from=2026-02-01&date_to=2026-02-28
```

| Query Param | Type | Description |
|------------|------|-------------|
| `page` | int | Page number (default: 1) |
| `page_size` | int | Items per page (default: 20, max: 100) |
| `project_id` | uint | Filter by project |
| `employee_id` | string | Filter by employee |
| `status` | string | `pending`, `in_progress`, `completed`, `blocked` |
| `category` | string | Filter by work category |
| `date_from` | string | From date (YYYY-MM-DD) |
| `date_to` | string | To date (YYYY-MM-DD) |
| `search` | string | Search in title/description |

### 3.2 Get Task Categories

```
GET /api/v1/daily-tasks/categories
```

**Response:**
```json
{
  "success": true,
  "data": [
    { "value": "development", "label": "Development" },
    { "value": "testing", "label": "Testing / QA" },
    { "value": "design", "label": "Design / UI/UX" },
    { "value": "documentation", "label": "Documentation" },
    { "value": "meeting", "label": "Meeting" },
    { "value": "research", "label": "Research" },
    { "value": "planning", "label": "Planning" },
    { "value": "review", "label": "Code Review" },
    { "value": "support", "label": "Support / Bug Fix" },
    { "value": "other", "label": "Other" }
  ]
}
```

### 3.3 Get Project Task Summary

```
GET /api/v1/daily-tasks/project/:project_id/summary
```

### 3.4 Get Task Details

```
GET /api/v1/daily-tasks/:id
```

---

## 5. Daily Tasks API — Employee Self-Service

Base URL: `/api/v1/self-service/daily-tasks`

### 4.1 Log a Daily Task

```
POST /api/v1/self-service/daily-tasks
```

**Request Body:**
```json
{
  "project_id": 1,
  "task_date": "2026-02-06",
  "title": "Implemented user authentication module",
  "description": "Built login, registration, and JWT token refresh endpoints",
  "hours_spent": 6.5,
  "status": "completed",
  "category": "development",
  "progress_contribution": 5.0,
  "tags": "auth, jwt, security",
  "notes": "Needs code review before merge"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `project_id` | uint | Yes | Project to log task against |
| `task_date` | string | No | Date (YYYY-MM-DD), defaults to today |
| `title` | string | Yes | Task title (2-300 chars) |
| `description` | string | No | Detailed description |
| `hours_spent` | float | Yes | Hours spent (0.1-24) |
| `status` | string | No | `pending`, `in_progress`, `completed`, `blocked` (default: `completed`) |
| `category` | string | Yes | Work category (see categories list) |
| `progress_contribution` | float | No | How much this task contributes to project completion (0-100%). Auto-updates project progress |
| `tags` | string | No | Comma-separated tags |
| `notes` | string | No | Additional notes |

**Progress Contribution Logic:**
- Each completed task's `progress_contribution` is summed to calculate the project's overall progress
- The system prevents the total from exceeding 100%
- Example: If a project has 10 equal milestones, each task could contribute 10% (progress_contribution: 10)

**Response:** `201 Created`
```json
{
  "success": true,
  "message": "Daily task logged successfully",
  "data": {
    "id": 42,
    "project_id": 1,
    "project_name": "ERP System Migration",
    "employee_id": "EMP001",
    "employee_name": "admin",
    "task_date": "2026-02-06",
    "title": "Implemented user authentication module",
    "hours_spent": 6.5,
    "status": "completed",
    "category": "development",
    "progress_contribution": 5.0
  }
}
```

### 4.2 List My Daily Tasks

```
GET /api/v1/self-service/daily-tasks?page=1&page_size=20&project_id=1&date_from=2026-02-01&date_to=2026-02-28
```

### 4.3 Get Today's Tasks

```
GET /api/v1/self-service/daily-tasks/today
```

Returns all tasks logged for today.

### 4.4 Get My Task Summary

```
GET /api/v1/self-service/daily-tasks/summary?date_from=2026-02-01&date_to=2026-02-28
```

**Response:**
```json
{
  "success": true,
  "data": {
    "total_tasks": 18,
    "total_hours": 95.5,
    "completed_tasks": 15,
    "pending_tasks": 2,
    "in_progress_tasks": 1,
    "blocked_tasks": 0,
    "by_category": [
      { "category": "development", "task_count": 10, "total_hours": 60 },
      { "category": "meeting", "task_count": 5, "total_hours": 20 },
      { "category": "testing", "task_count": 3, "total_hours": 15.5 }
    ],
    "by_project": [
      { "project_id": 1, "project_name": "ERP Migration", "task_count": 12, "total_hours": 70 },
      { "project_id": 2, "project_name": "Website Redesign", "task_count": 6, "total_hours": 25.5 }
    ]
  }
}
```

### 4.5 Update My Task

```
PUT /api/v1/self-service/daily-tasks/:id
```

**Request Body** (all fields optional):
```json
{
  "title": "Updated title",
  "hours_spent": 7.0,
  "status": "completed",
  "progress_contribution": 7.5,
  "notes": "Updated notes"
}
```

> Only the task creator can update their own tasks.

### 4.6 Delete My Task

```
DELETE /api/v1/self-service/daily-tasks/:id
```

> Only the task creator can delete their own tasks. Project progress is recalculated after deletion.

### 4.7 Get Task Categories

```
GET /api/v1/self-service/daily-tasks/categories
```

### 4.8 Get Task Details

```
GET /api/v1/self-service/daily-tasks/:id
```

---

## 6. My Projects — Employee Self-Service

### 5.1 List My Projects

```
GET /api/v1/self-service/projects?page=1&page_size=20&status=active
```

Returns only projects where the current employee is an active member.

---

## Data Models

### Project Statuses
| Status | Description |
|--------|-------------|
| `planning` | Project is in planning phase (default) |
| `active` | Project is actively being worked on |
| `on_hold` | Project is temporarily paused |
| `completed` | Project is finished |
| `cancelled` | Project was cancelled |

### Project Priorities
| Priority | Description |
|----------|-------------|
| `low` | Low priority |
| `medium` | Medium priority (default) |
| `high` | High priority |
| `critical` | Critical / urgent |

### Task Statuses
| Status | Description |
|--------|-------------|
| `pending` | Task not started |
| `in_progress` | Task being worked on |
| `completed` | Task finished (default when logging) |
| `blocked` | Task is blocked |

### Task Categories
| Category | Label |
|----------|-------|
| `development` | Development |
| `testing` | Testing / QA |
| `design` | Design / UI/UX |
| `documentation` | Documentation |
| `meeting` | Meeting |
| `research` | Research |
| `planning` | Planning |
| `review` | Code Review |
| `support` | Support / Bug Fix |
| `other` | Other |

### Member Roles
| Role | Description |
|------|-------------|
| `lead` | Project lead |
| `member` | Regular team member (default) |
| `reviewer` | Reviewer / QA |

---

## Progress Calculation

The project progress is **automatically recalculated** whenever a daily task is created, updated, or deleted:

1. Each daily task has a `progress_contribution` field (0-100)
2. **Project progress** = Sum of `progress_contribution` from all **completed** tasks (capped at 100%)
3. The system prevents adding tasks that would push total progress beyond 100%
4. Managers can also manually override the progress via the Update Project endpoint

**Example flow:**
```
Project: ERP Migration (estimated 100% effort)

Day 1: Employee A logs "Database schema design"       → progress_contribution: 10%  → Project: 10%
Day 2: Employee B logs "API endpoint setup"            → progress_contribution: 15%  → Project: 25%
Day 3: Employee A logs "Authentication module"         → progress_contribution: 10%  → Project: 35%
Day 4: Employee C logs "Frontend dashboard"            → progress_contribution: 20%  → Project: 55%
...
```

---

## Error Responses

All errors follow the standard format:
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error info"
}
```

Common errors:
- `400` — Validation error, bad request
- `401` — Unauthorized (missing/invalid token)
- `403` — Forbidden (insufficient role)
- `404` — Resource not found
- `422` — Validation failed
