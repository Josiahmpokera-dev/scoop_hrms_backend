# Performance Management API Documentation

## Overview

The Performance Management module covers goal tracking, performance appraisals, 360° feedback, talent review, and reporting. All endpoints require authentication. Role-based access is enforced at the API level.

**Base URL:** `/api/v1/performance`

### Role Access Matrix

| API Category              | Employee | Manager | HR   | Admin |
|---------------------------|----------|---------|------|-------|
| Own goals & self review   | Yes      | Yes     | Yes  | Yes   |
| Team goals & reviews      | No       | Yes     | Yes  | Yes   |
| Appraisal cycles mgmt     | No       | No      | Yes  | Yes   |
| 360° Feedback mgmt        | No       | No      | Yes  | Yes   |
| Talent review / 9-Box     | No       | Yes     | Yes  | Yes   |
| Performance reports        | No       | No      | Yes  | Yes   |

---

## Section 1: Performance Dashboard

### 1.1 Get Dashboard Statistics

Aggregated stats for the current user's dashboard cards and quick stats.

**`GET /api/v1/performance/dashboard/stats`**

**Auth:** Any authenticated user

**Response:**

```json
{
  "success": true,
  "data": {
    "goals": {
      "total": 8,
      "completed": 3,
      "in_progress": 4,
      "at_risk": 1,
      "avg_progress": 57
    },
    "appraisal": {
      "status": "Self Review",
      "self_rating": 4.2,
      "last_rating": 4.3
    },
    "feedback_360": {
      "total": 5,
      "available": 2
    },
    "quick_stats": {
      "goals_on_track": 4,
      "review_completion_percent": 63,
      "last_overall_rating": 4.3
    }
  }
}
```

---

### 1.2 Get Upcoming Actions

Personalized action items (self reviews due, goal check-ins, 360 feedback requests).

**`GET /api/v1/performance/dashboard/upcoming-actions`**

**Auth:** Any authenticated user

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "type": "self_review",
      "title": "Complete Self Review",
      "description": "Annual performance self review",
      "due_date": "2026-02-28",
      "link": "/performance/appraisals",
      "priority": "high"
    },
    {
      "type": "goal_checkin",
      "title": "Q4 Goal Check-in",
      "description": "Update progress on Q4 goals",
      "due_date": "2026-02-20",
      "link": "/performance/goals",
      "priority": "medium"
    },
    {
      "type": "feedback_360",
      "title": "360° Feedback",
      "description": "Provide peer feedback for John Doe",
      "due_date": "2026-03-01",
      "link": "/performance/feedback-360",
      "priority": "medium"
    }
  ]
}
```

---

## Section 2: Goals & OKRs

### 2.1 List Goals

**`GET /api/v1/performance/goals`**

**Auth:** Employee (own goals), Manager (team goals), HR/Admin (all)

**Query Parameters:**

| Param        | Type     | Required | Description                                              |
|-------------|----------|----------|----------------------------------------------------------|
| `level`     | `string` | No       | Filter: `Individual`, `Team`, `Department`, `Company`    |
| `status`    | `string` | No       | Filter: `Not Started`, `In Progress`, `At Risk`, `Completed`, `Cancelled` |
| `owner_id`  | `uint`   | No       | Filter by owner employee ID                              |
| `department`| `string` | No       | Filter by department                                     |
| `page`      | `int`    | No       | Page number (default 1)                                  |
| `page_size` | `int`    | No       | Items per page (default 20)                              |
| `search`    | `string` | No       | Search by title/description                              |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "G-001",
      "title": "Increase Revenue by 30%",
      "description": "Grow annual revenue through new customer acquisition and upselling",
      "type": "OKR",
      "level": "Company",
      "owner": "Jane Smith",
      "owner_id": 3,
      "department": "Executive",
      "parent_goal_id": null,
      "weight": 100,
      "due_date": "2026-12-31",
      "status": "In Progress",
      "progress": 45,
      "visibility": "Public",
      "aligned_to": [],
      "tags": ["Strategic", "Revenue"],
      "key_results": [
        {
          "id": "KR-001",
          "title": "Acquire 50 new enterprise customers",
          "target_value": 50,
          "current_value": 22,
          "unit": "customers",
          "due_date": "2026-12-31",
          "status": "On Track"
        },
        {
          "id": "KR-002",
          "title": "Increase average deal size to $50K",
          "target_value": 50000,
          "current_value": 38000,
          "unit": "USD",
          "due_date": "2026-12-31",
          "status": "At Risk"
        }
      ],
      "check_ins": [],
      "created_at": "2026-01-01"
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 8,
    "total_pages": 1
  }
}
```

---

### 2.2 Get Goal Statistics

**`GET /api/v1/performance/goals/stats`**

**Auth:** Any authenticated user

**Response:**

```json
{
  "success": true,
  "data": {
    "total": 8,
    "in_progress": 4,
    "completed": 3,
    "at_risk": 1,
    "not_started": 0,
    "cancelled": 0,
    "avg_progress": 57
  }
}
```

---

### 2.3 Get Goal Detail

**`GET /api/v1/performance/goals/:id`**

**Auth:** Owner, Manager, HR, Admin

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "G-001",
    "title": "Increase Revenue by 30%",
    "description": "...",
    "type": "OKR",
    "level": "Company",
    "owner": "Jane Smith",
    "owner_id": 3,
    "department": "Executive",
    "parent_goal_id": null,
    "weight": 100,
    "due_date": "2026-12-31",
    "status": "In Progress",
    "progress": 45,
    "visibility": "Public",
    "aligned_to": [],
    "tags": ["Strategic"],
    "key_results": [],
    "check_ins": [
      {
        "id": "CI-001",
        "date": "2026-01-15",
        "progress": 45,
        "status": "On Track",
        "notes": "Good progress on Q1 targets",
        "evidence_url": null,
        "updated_by": "Jane Smith"
      }
    ],
    "created_at": "2026-01-01",
    "updated_at": "2026-01-15"
  }
}
```

---

### 2.4 Create Goal

**`POST /api/v1/performance/goals`**

**Auth:** Any authenticated user (own goals), Manager (team), HR/Admin (any)

**Request:**

```json
{
  "title": "Improve customer satisfaction score",
  "description": "Raise NPS from 45 to 60",
  "type": "OKR",
  "level": "Individual",
  "weight": 30,
  "due_date": "2026-12-31",
  "parent_goal_id": "G-002",
  "visibility": "Manager",
  "tags": ["Customer Success"]
}
```

| Field             | Type       | Required | Description                                          |
|-------------------|------------|----------|------------------------------------------------------|
| `title`           | `string`   | Yes      | Goal title                                           |
| `description`     | `string`   | No       | Detailed description                                 |
| `type`            | `string`   | Yes      | `OKR`, `KPI`, `Project`, `Development`               |
| `level`           | `string`   | Yes      | `Individual`, `Team`, `Department`, `Company`        |
| `weight`          | `number`   | Yes      | Weight percentage (0-100)                            |
| `due_date`        | `string`   | Yes      | ISO 8601 date                                        |
| `parent_goal_id`  | `string?`  | No       | ID of parent goal for alignment                      |
| `visibility`      | `string`   | No       | `Public`, `Manager`, `Private` (default: `Manager`)  |
| `tags`            | `string[]` | No       | Tags for categorization                              |

**Response:** `201 Created` with full goal object.

---

### 2.5 Update Goal

**`PATCH /api/v1/performance/goals/:id`**

**Auth:** Owner, Manager, HR, Admin

**Request (partial):**

```json
{
  "title": "Updated goal title",
  "status": "In Progress",
  "progress": 65
}
```

**Response:** Updated goal object.

---

### 2.6 Delete Goal

**`DELETE /api/v1/performance/goals/:id`**

**Auth:** Owner (if Not Started/Draft), HR, Admin

**Response:**

```json
{
  "success": true,
  "message": "Goal deleted successfully"
}
```

---

### 2.7 Create Goal Check-in

**`POST /api/v1/performance/goals/:goalId/check-ins`**

**Auth:** Goal owner, Manager

**Request:**

```json
{
  "progress": 75,
  "status": "On Track",
  "notes": "Completed Q1 milestones ahead of schedule"
}
```

| Field      | Type     | Required | Description                                 |
|------------|----------|----------|---------------------------------------------|
| `progress` | `number` | Yes      | Progress percentage (0-100)                 |
| `status`   | `string` | Yes      | `On Track`, `At Risk`, `Behind`             |
| `notes`    | `string` | Yes      | Check-in notes                              |

**Response:** `201 Created` with check-in object.

---

### 2.8 Manage Key Results

**Create Key Result:**

**`POST /api/v1/performance/goals/:goalId/key-results`**

**Request:**

```json
{
  "title": "Acquire 50 new enterprise customers",
  "target_value": 50,
  "current_value": 0,
  "unit": "customers",
  "due_date": "2026-12-31"
}
```

**Update Key Result:**

**`PATCH /api/v1/performance/goals/:goalId/key-results/:krId`**

**Delete Key Result:**

**`DELETE /api/v1/performance/goals/:goalId/key-results/:krId`**

---

### 2.9 Get Goal Alignment Map

**`GET /api/v1/performance/goals/alignment`**

**Auth:** Manager, HR, Admin

**Response:**

```json
{
  "success": true,
  "data": {
    "nodes": [
      { "id": "G-001", "title": "Increase Revenue", "level": "Company", "progress": 45, "owner": "CEO" },
      { "id": "G-002", "title": "Expand Sales Team", "level": "Department", "parent_id": "G-001", "progress": 65 }
    ],
    "edges": [
      { "from": "G-002", "to": "G-001" }
    ]
  }
}
```

---

### 2.10 Submit Goal for Approval

When an **Employee** creates or updates a goal, they must submit it for manager approval before it becomes active.

**`POST /api/v1/performance/goals/:id/submit-for-approval`**

**Auth:** Employee (goal owner)

**Request:**

```json
{
  "comments": "Please review my Q2 goal"
}
```

| Field      | Type     | Required | Description          |
|------------|----------|----------|----------------------|
| `comments` | `string` | No       | Optional note        |

**Response:**

```json
{
  "success": true,
  "message": "Goal submitted for approval",
  "data": {
    "id": "G-003",
    "approval_status": "Pending Approval",
    "submitted_date": "2026-02-01T10:00:00Z"
  }
}
```

---

### 2.11 Approve / Reject Goal

**Manager** approves or rejects a submitted goal.

**`PUT /api/v1/performance/goals/:id/approve`**

**Auth:** Manager (of the goal owner), HR, Admin

**Request:**

```json
{
  "action": "approve",
  "comments": "Approved. Well-defined targets."
}
```

| Field      | Type     | Required | Description                           |
|------------|----------|----------|---------------------------------------|
| `action`   | `string` | Yes      | `approve` or `reject`                 |
| `comments` | `string` | No       | Manager feedback                      |

**Response:**

```json
{
  "success": true,
  "message": "Goal approved successfully",
  "data": {
    "id": "G-003",
    "approval_status": "Approved",
    "approved_by": "John Mwita",
    "approved_date": "2026-02-02T09:00:00Z"
  }
}
```

---

### 2.12 Request Goal Completion

When an **Employee** believes they have achieved a goal, they request completion verification from their manager.

**`POST /api/v1/performance/goals/:id/request-completion`**

**Auth:** Employee (goal owner)

**Request:**

```json
{
  "evidence": "Completed all 3 key results. See attached Q4 report.",
  "evidence_url": "https://docs.example.com/q4-report",
  "final_progress": 100
}
```

| Field            | Type     | Required | Description                       |
|------------------|----------|----------|-----------------------------------|
| `evidence`       | `string` | Yes      | Summary of achievement            |
| `evidence_url`   | `string` | No       | Link to supporting evidence       |
| `final_progress` | `number` | Yes      | Final progress percentage (0-100) |

**Response:**

```json
{
  "success": true,
  "message": "Completion request submitted",
  "data": {
    "id": "G-003",
    "completion_status": "Pending Verification",
    "requested_date": "2026-03-15T10:00:00Z"
  }
}
```

---

### 2.13 Verify Goal Completion

**Manager** verifies or rejects a goal completion request.

**`PUT /api/v1/performance/goals/:id/verify-completion`**

**Auth:** Manager (of the goal owner), HR, Admin

**Request:**

```json
{
  "action": "verify",
  "final_rating": 4,
  "comments": "Well done. All key results achieved on time."
}
```

| Field          | Type     | Required | Description                               |
|----------------|----------|----------|-------------------------------------------|
| `action`       | `string` | Yes      | `verify` or `reject`                      |
| `final_rating` | `number` | No       | Manager rating for goal achievement (1-5) |
| `comments`     | `string` | No       | Manager feedback                          |

**Response (verify):**

```json
{
  "success": true,
  "message": "Goal completion verified",
  "data": {
    "id": "G-003",
    "status": "Completed",
    "completion_status": "Verified",
    "final_rating": 4,
    "verified_by": "John Mwita",
    "verified_date": "2026-03-16T09:00:00Z"
  }
}
```

**Response (reject):**

```json
{
  "success": true,
  "message": "Goal completion rejected",
  "data": {
    "id": "G-003",
    "status": "In Progress",
    "completion_status": "Rejected",
    "comments": "Key result KR-002 is not yet at target. Please continue."
  }
}
```

---

### 2.14 Assign Goal to Employee

**Manager** or **HR/Admin** assigns a goal to a specific employee.

**`POST /api/v1/performance/goals/assign`**

**Auth:** Manager (to their team), HR, Admin

**Request:**

```json
{
  "title": "Improve Code Quality",
  "description": "Reduce technical debt and improve code coverage",
  "type": "KPI",
  "level": "Individual",
  "assigned_to": 5,
  "weight": 30,
  "due_date": "2026-12-31",
  "parent_goal_id": "G-002",
  "key_results": [
    {
      "title": "Achieve 90% code coverage",
      "target_value": 90,
      "current_value": 75,
      "unit": "%",
      "due_date": "2026-12-31"
    }
  ],
  "tags": ["Technical", "Quality"]
}
```

| Field           | Type       | Required | Description                                   |
|-----------------|------------|----------|-----------------------------------------------|
| `title`         | `string`   | Yes      | Goal title                                    |
| `description`   | `string`   | No       | Goal description                              |
| `type`          | `string`   | Yes      | `OKR`, `KPI`, `Project`, `Development`        |
| `level`         | `string`   | Yes      | `Individual`, `Team`, `Department`, `Company` |
| `assigned_to`   | `uint`     | Yes      | Employee ID to assign goal to                 |
| `weight`        | `number`   | Yes      | Weight percentage                             |
| `due_date`      | `string`   | Yes      | ISO 8601 date                                 |
| `parent_goal_id`| `string?`  | No       | Parent goal ID for alignment                  |
| `key_results`   | `array`    | No       | Initial key results                           |
| `tags`          | `string[]` | No       | Tags                                          |

**Response:** `201 Created` with full goal object (approval_status = `Approved` since manager assigned it).

---

### Goal Status Workflow (Updated)

```
                     ┌──────────────┐
                     │  Not Started │
                     └──────┬───────┘
                            │ Employee creates goal
                     ┌──────▼───────────────┐
                     │  Pending Approval     │
                     └──────┬───────────────┘
                   ┌────────┴────────┐
                   │                 │
            ┌──────▼──────┐  ┌──────▼──────┐
            │  Approved   │  │  Rejected   │ → Employee revises and resubmits
            └──────┬──────┘  └─────────────┘
                   │
            ┌──────▼──────┐
            │ In Progress  │ ◄──── At Risk (recovery)
            └──────┬──────┘
                   │                 │
                   ▼                 ▼
        ┌─────────────────┐   ┌────────────┐
        │Pending Verify.  │   │  At Risk   │
        └────────┬────────┘   └────────────┘
          ┌──────┴──────┐
          │             │
   ┌──────▼──────┐ ┌───▼───────────┐
   │ Completed   │ │ Rejected      │ → Back to In Progress
   │ (Verified)  │ │ (needs more)  │
   └─────────────┘ └───────────────┘

Any status → Cancelled
```

### Goal Approval Status Values

- `Draft` — Not yet submitted
- `Pending Approval` — Waiting for manager review
- `Approved` — Manager approved
- `Rejected` — Manager rejected (employee can revise)

### Goal Completion Status Values

- `null` — Not applicable yet
- `Pending Verification` — Employee requested completion
- `Verified` — Manager confirmed completion
- `Rejected` — Manager rejected (back to In Progress)

### Key Result Status Values

- `On Track`
- `At Risk`
- `Behind`
- `Completed`

---

## Section 2B: Department Targets

Department Targets are set by **Department Leaders** (managers/heads of department) to define measurable KPIs and objectives for the entire department. Individual employee goals should align to these targets.

### 2B.1 List Department Targets

**`GET /api/v1/performance/department-targets`**

**Auth:** Employee (own department), Manager (own department + team), HR/Admin (all)

**Query Parameters:**

| Param        | Type     | Required | Description                                              |
|-------------|----------|----------|----------------------------------------------------------|
| `department`| `string` | No       | Filter by department                                     |
| `category`  | `string` | No       | `Revenue`, `Efficiency`, `Quality`, `Growth`, `Customer`, `People`, `Innovation` |
| `status`    | `string` | No       | `Not Started`, `On Track`, `At Risk`, `Behind`, `Achieved`, `Missed` |
| `period`    | `string` | No       | Filter by period (e.g. `Q1 2026`, `H1 2026`, `FY 2026`) |
| `page`      | `int`    | No       | Page number                                              |
| `page_size` | `int`    | No       | Items per page                                           |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "DT-001",
      "title": "Achieve 95% Sprint Velocity",
      "description": "Maintain consistent sprint velocity across all engineering teams",
      "department": "Engineering",
      "category": "Efficiency",
      "metric": "Sprint Velocity Achievement",
      "target_value": 95,
      "current_value": 88,
      "unit": "%",
      "period": "Q1 2026",
      "start_date": "2026-01-01",
      "end_date": "2026-03-31",
      "status": "On Track",
      "priority": "High",
      "created_by": "John Mwita",
      "created_by_id": 2,
      "approved_by": "CEO",
      "linked_goal_ids": ["G-002", "G-003"],
      "progress": 88,
      "milestones": [
        {
          "id": "DM-001",
          "title": "Establish baseline velocity",
          "due_date": "2026-01-15",
          "status": "Completed",
          "completed_date": "2026-01-14"
        }
      ],
      "created_at": "2025-12-15"
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 4, "total_pages": 1 }
}
```

---

### 2B.2 Get Department Target Detail

**`GET /api/v1/performance/department-targets/:id`**

**Auth:** Department members, Manager, HR, Admin

**Response:** Full `DepartmentTarget` object.

---

### 2B.3 Create Department Target

**`POST /api/v1/performance/department-targets`**

**Auth:** Manager (own department), HR, Admin

**Request:**

```json
{
  "title": "Reduce Production Bugs by 40%",
  "description": "Decrease critical bugs through improved testing and code review",
  "department": "Engineering",
  "category": "Quality",
  "metric": "Bug Count Reduction",
  "target_value": 40,
  "unit": "% reduction",
  "period": "H1 2026",
  "start_date": "2026-01-01",
  "end_date": "2026-06-30",
  "priority": "Critical",
  "milestones": [
    { "title": "Implement mandatory code review policy", "due_date": "2026-01-31" },
    { "title": "Achieve 80% unit test coverage", "due_date": "2026-03-31" }
  ]
}
```

| Field          | Type       | Required | Description                                                     |
|----------------|------------|----------|-----------------------------------------------------------------|
| `title`        | `string`   | Yes      | Target title                                                    |
| `description`  | `string`   | No       | Detailed description                                            |
| `department`   | `string`   | Yes      | Department name                                                 |
| `category`     | `string`   | Yes      | `Revenue`, `Efficiency`, `Quality`, `Growth`, `Customer`, `People`, `Innovation` |
| `metric`       | `string`   | Yes      | What is being measured                                          |
| `target_value` | `number`   | Yes      | Target number                                                   |
| `unit`         | `string`   | Yes      | Unit of measurement                                             |
| `period`       | `string`   | Yes      | Period label (e.g. `Q1 2026`)                                   |
| `start_date`   | `string`   | Yes      | ISO 8601 start date                                             |
| `end_date`     | `string`   | Yes      | ISO 8601 end date                                               |
| `priority`     | `string`   | Yes      | `Critical`, `High`, `Medium`, `Low`                             |
| `milestones`   | `array`    | No       | Array of milestone objects with `title` and `due_date`          |

**Response:** `201 Created` with full target object.

---

### 2B.4 Update Department Target

**`PATCH /api/v1/performance/department-targets/:id`**

**Auth:** Creator, HR, Admin

**Request (partial):**

```json
{
  "current_value": 92,
  "status": "On Track"
}
```

**Response:** Updated target object.

---

### 2B.5 Update Target Progress (Check-in)

**`POST /api/v1/performance/department-targets/:id/progress`**

**Auth:** Creator, HR, Admin

**Request:**

```json
{
  "current_value": 92,
  "notes": "Sprint velocity improved after process changes",
  "status": "On Track"
}
```

**Response:** Updated target with progress history.

---

### 2B.6 Complete Milestone

**`PUT /api/v1/performance/department-targets/:id/milestones/:milestoneId/complete`**

**Auth:** Creator, HR, Admin

**Response:** Updated milestone with `status: "Completed"` and `completed_date`.

---

### 2B.7 Link Goal to Department Target

**`POST /api/v1/performance/department-targets/:id/link-goal`**

**Auth:** Manager, HR, Admin

**Request:**

```json
{
  "goal_id": "G-003"
}
```

**Response:** Updated target with linked goal added.

---

### 2B.8 Delete Department Target

**`DELETE /api/v1/performance/department-targets/:id`**

**Auth:** Creator (if Not Started), HR, Admin

**Response:**

```json
{
  "success": true,
  "message": "Department target deleted successfully"
}
```

---

### Department Target Status Values

- `Not Started` — Target period hasn't begun
- `On Track` — Progress is meeting expectations
- `At Risk` — Progress is below expectations but recoverable
- `Behind` — Significantly behind schedule
- `Achieved` — Target has been met or exceeded
- `Missed` — Target period ended without achievement

### Who Creates Department Targets?

| Role                   | Can Create For           |
|------------------------|--------------------------|
| **Department Leader**  | Their own department     |
| **HR**                 | Any department           |
| **Admin**              | Any department           |

### Who Can See Department Targets?

| Role                       | Visibility                                           |
|----------------------------|------------------------------------------------------|
| **Employee**               | Only targets belonging to **their own department**   |
| **Manager / Dept Leader**  | Only targets belonging to **their own department**   |
| **HR**                     | **All departments** (with department filter)         |
| **Admin / Super Admin**    | **All departments** (with department filter)         |

> **Important:** Every employee in a department can see the targets and projects set by their department leader. This ensures transparency — everyone knows what the department is aiming for and can align their individual goals accordingly.

### How It Works

1. **Department Leader** creates a target with measurable KPIs and milestones
2. The target is **automatically visible** to all members of that department
3. **Employees** can view their department's targets and align their individual goals to them
4. **Department Leader** tracks progress, updates milestones, and links employee goals
5. **Individual goal completion** contributes to department target progress
6. **HR/Admin** can view targets across all departments for oversight
7. At period end, target is marked `Achieved` or `Missed`

---

## Section 2C: Employee Targets (Assigned Targets)

Employee Targets are **individual measurable KPIs** assigned to specific employees by their **Manager, Project Manager, HR, or Admin**. They cascade down from Department Targets and can be linked to Projects, creating a complete chain: **Department Target → Employee Target → Project**.

### 2C.1 List Employee Targets

**`GET /api/v1/performance/employee-targets`**

**Auth:** Employee (own targets), Manager (team targets), HR/Admin (all)

**Query Parameters:**

| Param                 | Type     | Required | Description                                    |
|----------------------|----------|----------|------------------------------------------------|
| `employee_id`        | `string` | No       | Filter by specific employee                    |
| `department`         | `string` | No       | Filter by department                           |
| `department_target_id`| `string`| No       | Filter by parent department target             |
| `status`             | `string` | No       | `Not Started`, `In Progress`, `On Track`, etc  |
| `period`             | `string` | No       | Filter by period                               |
| `assigned_by`        | `string` | No       | Filter by who assigned                         |
| `page`               | `int`    | No       | Page number                                    |
| `page_size`          | `int`    | No       | Items per page                                 |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "ET-001",
      "title": "Deliver 3 API Modules by End of Q1",
      "description": "Complete and deploy 3 backend API modules...",
      "employee_id": "EMP-003",
      "employee_name": "Peter Ochieng",
      "department": "Engineering",
      "department_target_id": "DT-001",
      "department_target_title": "Achieve 95% Sprint Velocity",
      "linked_project_id": 1,
      "linked_project_code": "PRJ-0001",
      "linked_project_name": "ERP System Migration",
      "metric": "Modules Delivered",
      "target_value": 3,
      "current_value": 2,
      "unit": "modules",
      "weight": 30,
      "due_date": "2026-03-31",
      "period": "Q1 2026",
      "status": "On Track",
      "priority": "High",
      "assigned_by": "John Mwita",
      "assigned_date": "2026-01-05",
      "progress": 67
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 5, "total_pages": 1 }
}
```

---

### 2C.2 Get Employee Target Detail

**`GET /api/v1/performance/employee-targets/:id`**

**Auth:** Target employee, Manager, HR, Admin

---

### 2C.3 Assign Target to Single Employee

**`POST /api/v1/performance/employee-targets`**

**Auth:** Manager (own department), Project Manager, HR, Admin

**Request:**

```json
{
  "employee_id": "EMP-003",
  "title": "Deliver 3 API Modules by End of Q1",
  "description": "Complete and deploy 3 backend API modules with full test coverage",
  "department_target_id": "DT-001",
  "linked_project_id": 1,
  "metric": "Modules Delivered",
  "target_value": 3,
  "unit": "modules",
  "weight": 30,
  "due_date": "2026-03-31",
  "period": "Q1 2026",
  "priority": "High"
}
```

| Field                  | Type     | Required | Description                                           |
|-----------------------|----------|----------|-------------------------------------------------------|
| `employee_id`         | `string` | Yes      | Employee to assign the target to                      |
| `title`               | `string` | Yes      | Target title                                          |
| `description`         | `string` | No       | Detailed description                                  |
| `department_target_id`| `string` | No       | Link to a department target (cascading alignment)     |
| `linked_project_id`   | `number` | No       | Link to a project from Manage Project                 |
| `metric`              | `string` | Yes      | What is being measured                                |
| `target_value`        | `number` | Yes      | Target number to achieve                              |
| `unit`                | `string` | Yes      | Unit of measurement                                   |
| `weight`              | `number` | Yes      | Weight in overall performance (%)                     |
| `due_date`            | `string` | Yes      | ISO 8601 due date                                     |
| `period`              | `string` | Yes      | Period label (e.g. `Q1 2026`)                         |
| `priority`            | `string` | Yes      | `Critical`, `High`, `Medium`, `Low`                   |

**Response:** `201 Created` with full employee target object.

---

### 2C.4 Bulk Assign Targets to Department Employees

Assign the **same target template** to multiple employees in the same department at once. Each employee gets their own individual target with the same metric and due date.

**`POST /api/v1/performance/employee-targets/bulk-assign`**

**Auth:** Manager (own department), HR, Admin

**Request:**

```json
{
  "employee_ids": ["EMP-003", "EMP-004", "EMP-008"],
  "title": "Achieve 90% Code Review Compliance",
  "description": "Ensure all PRs go through code review with no exceptions",
  "department": "Engineering",
  "department_target_id": "DT-002",
  "linked_project_id": null,
  "metric": "Code Review Compliance",
  "target_value": 90,
  "unit": "%",
  "weight": 20,
  "due_date": "2026-06-30",
  "period": "H1 2026",
  "priority": "Medium"
}
```

| Field                  | Type       | Required | Description                                         |
|-----------------------|------------|----------|-----------------------------------------------------|
| `employee_ids`        | `string[]` | Yes      | Array of employee IDs to assign to (same department) |
| `title`               | `string`   | Yes      | Target title (same for all employees)                |
| `description`         | `string`   | No       | Detailed description                                 |
| `department`          | `string`   | Yes      | Department name (employees must belong to this dept) |
| `department_target_id`| `string`   | No       | Link to parent department target                     |
| `linked_project_id`   | `number`   | No       | Link to a project                                    |
| `metric`              | `string`   | Yes      | What is being measured                               |
| `target_value`        | `number`   | Yes      | Target number to achieve                             |
| `unit`                | `string`   | Yes      | Unit of measurement                                  |
| `weight`              | `number`   | Yes      | Weight in overall performance                        |
| `due_date`            | `string`   | Yes      | ISO 8601 due date                                    |
| `period`              | `string`   | Yes      | Period label                                         |
| `priority`            | `string`   | Yes      | `Critical`, `High`, `Medium`, `Low`                  |

**Response:**

```json
{
  "success": true,
  "message": "Target assigned to 3 employees successfully",
  "data": {
    "created_count": 3,
    "targets": [
      { "id": "ET-010", "employee_id": "EMP-003", "employee_name": "Peter Ochieng" },
      { "id": "ET-011", "employee_id": "EMP-004", "employee_name": "Mary Wanjiku" },
      { "id": "ET-012", "employee_id": "EMP-008", "employee_name": "David Kimani" }
    ]
  }
}
```

---

### 2C.5 Update Employee Target Progress

**`POST /api/v1/performance/employee-targets/:id/progress`**

**Auth:** Target employee, Manager, HR, Admin

**Request:**

```json
{
  "current_value": 2,
  "status": "On Track",
  "notes": "Second module deployed to staging"
}
```

---

### 2C.6 Update Employee Target

**`PATCH /api/v1/performance/employee-targets/:id`**

**Auth:** Assigner, Manager, HR, Admin

---

### 2C.7 Delete Employee Target

**`DELETE /api/v1/performance/employee-targets/:id`**

**Auth:** Assigner (if Not Started), HR, Admin

---

### Who Can Assign Employee Targets?

| Role                   | Can Assign To                               |
|------------------------|---------------------------------------------|
| **Manager / HOD**      | Employees in their department               |
| **Project Manager**    | Employees assigned to their project         |
| **HR**                 | Any employee                                |
| **Admin / Super Admin**| Any employee                                |

### Complete Alignment Chain

```
Company Goal (e.g., "Increase Revenue by 30%")
  └── Department Target (e.g., "Achieve 95% Sprint Velocity")   ← Created by Dept Leader
        ├── Employee Target (e.g., "Deliver 3 API Modules")     ← Assigned to Peter
        ├── Employee Target (e.g., "Complete Dashboard")         ← Assigned to Mary
        ├── Employee Target (e.g., "90% Code Review")           ← Bulk-assigned to all Eng
        └── Linked Project (e.g., "PRJ-0001 ERP Migration")     ← From Manage Project
```

### Employee Target vs Goal

| Feature           | Employee Target                          | Goal (existing)                          |
|-------------------|------------------------------------------|------------------------------------------|
| **Created by**    | Manager/HR assigns to employee           | Employee creates for themselves           |
| **Approval**      | Pre-approved (assigned by authority)     | Requires manager approval workflow        |
| **Metric**        | Always has a measurable metric + value   | May or may not have key results           |
| **Alignment**     | Always linked to department target       | Optionally aligned to parent goals        |
| **Use case**      | "You must achieve X by this date"        | "I want to improve X this quarter"        |

---

## Section 3: Performance Appraisals

### 3.1 List Appraisal Cycles

**`GET /api/v1/performance/appraisal-cycles`**

**Auth:** HR, Admin (full list), Manager/Employee (own cycles)

**Query Parameters:**

| Param    | Type     | Required | Description                                    |
|----------|----------|----------|------------------------------------------------|
| `status` | `string` | No       | `Draft`, `Active`, `Calibration`, `Completed`, `Closed` |
| `year`   | `int`    | No       | Filter by year                                 |
| `type`   | `string` | No       | `Annual`, `Mid-Year`, `Quarterly`, `Probation` |
| `page`   | `int`    | No       | Page number                                    |
| `page_size` | `int` | No       | Items per page                                 |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "AC-001",
      "name": "Annual Review 2026",
      "type": "Annual",
      "year": 2026,
      "start_date": "2026-01-01",
      "end_date": "2026-12-31",
      "status": "Active",
      "participants": 150,
      "completed": 45,
      "workflow": [
        {
          "step": 1,
          "name": "Self Assessment",
          "role": "Employee",
          "due_date": "2026-02-28",
          "status": "In Progress",
          "participants": 150,
          "completed": 45
        },
        {
          "step": 2,
          "name": "Manager Review",
          "role": "Manager",
          "due_date": "2026-03-15",
          "status": "Pending",
          "participants": 150,
          "completed": 0
        },
        {
          "step": 3,
          "name": "HR Calibration",
          "role": "HR",
          "due_date": "2026-03-31",
          "status": "Pending",
          "participants": 150,
          "completed": 0
        }
      ],
      "rating_scale": { "min": 1, "max": 5, "labels": ["Needs Improvement", "Partially Meets", "Meets", "Exceeds", "Outstanding"] },
      "created_at": "2025-12-01"
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 3, "total_pages": 1 }
}
```

---

### 3.2 Get Active Appraisal Cycle

**`GET /api/v1/performance/appraisal-cycles/active`**

**Auth:** Any authenticated user

**Response:** Single `AppraisalCycle` object (same shape as list item).

---

### 3.3 Get Appraisal Cycle Detail

**`GET /api/v1/performance/appraisal-cycles/:id`**

**Auth:** HR, Admin

**Response:** Full `AppraisalCycle` with workflow steps and participant counts.

---

### 3.4 Create Appraisal Cycle

**`POST /api/v1/performance/appraisal-cycles`**

**Auth:** HR, Admin

**Request:**

```json
{
  "name": "Annual Review 2026",
  "type": "Annual",
  "year": 2026,
  "start_date": "2026-01-01",
  "end_date": "2026-12-31",
  "workflow": [
    { "step": 1, "name": "Self Assessment", "role": "Employee", "due_date": "2026-02-28" },
    { "step": 2, "name": "Manager Review", "role": "Manager", "due_date": "2026-03-15" },
    { "step": 3, "name": "HR Calibration", "role": "HR", "due_date": "2026-03-31" }
  ],
  "rating_scale": { "min": 1, "max": 5 },
  "sections": ["Goals", "Competencies", "Values"]
}
```

**Response:** `201 Created` with full cycle object.

---

### 3.5 Update Appraisal Cycle

**`PUT /api/v1/performance/appraisal-cycles/:id`**

**Auth:** HR, Admin

---

### 3.6 Change Cycle Status

**`PATCH /api/v1/performance/appraisal-cycles/:id/status`**

**Auth:** HR, Admin

**Request:**

```json
{
  "status": "Active"
}
```

**Valid transitions:** `Draft → Active → Calibration → Completed → Closed`

---

### 3.7 Delete Appraisal Cycle

**`DELETE /api/v1/performance/appraisal-cycles/:id`**

**Auth:** Admin (Draft cycles only)

---

### 3.8 List Appraisals

**`GET /api/v1/performance/appraisals`**

**Auth:** Employee (own), Manager (team), HR/Admin (all)

**Query Parameters:**

| Param         | Type     | Required | Description                                           |
|---------------|----------|----------|-------------------------------------------------------|
| `mine`        | `bool`   | No       | Only my appraisals (default false)                    |
| `employee_id` | `uint`   | No       | Filter by employee                                    |
| `cycle_id`    | `string` | No       | Filter by cycle                                       |
| `status`      | `string` | No       | `Not Started`, `Self Review`, `Manager Review`, `Calibration`, `Completed` |
| `department`  | `string` | No       | Filter by department                                  |
| `page`        | `int`    | No       | Page number                                           |
| `page_size`   | `int`    | No       | Items per page                                        |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "APR-001",
      "cycle_id": "AC-001",
      "cycle_name": "Annual Review 2026",
      "employee_id": 5,
      "employee_name": "John Doe",
      "designation": "Software Engineer",
      "department": "Engineering",
      "manager": "Jane Smith",
      "review_period": { "start": "2026-01-01", "end": "2026-12-31" },
      "current_step": "Self Review",
      "status": "Self Review",
      "overall_rating": null,
      "final_rating": null,
      "self_review": null,
      "manager_review": null,
      "submitted_date": null
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 45, "total_pages": 3 }
}
```

---

### 3.9 Get Appraisal Summary

**`GET /api/v1/performance/appraisals/summary`**

**Auth:** Any authenticated user

**Response:**

```json
{
  "success": true,
  "data": {
    "my_total": 3,
    "pending_count": 1,
    "completed_count": 2,
    "last_rating": 4.3
  }
}
```

---

### 3.10 Get Appraisal Detail

**`GET /api/v1/performance/appraisals/:id`**

**Auth:** Employee (own), Manager (team), HR/Admin

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "APR-001",
    "cycle_id": "AC-001",
    "cycle_name": "Annual Review 2026",
    "employee_id": 5,
    "employee_name": "John Doe",
    "designation": "Software Engineer",
    "department": "Engineering",
    "manager": "Jane Smith",
    "review_period": { "start": "2026-01-01", "end": "2026-12-31" },
    "status": "Completed",
    "self_review": {
      "goals_rating": 4,
      "competencies_rating": 4.5,
      "values_rating": 4,
      "overall_rating": 4.2,
      "strengths": ["Technical leadership", "Problem solving"],
      "areas_for_improvement": ["Delegation", "Documentation"],
      "comments": "Great year of growth.",
      "submitted_date": "2026-02-15",
      "submitted_by": "John Doe"
    },
    "manager_review": {
      "goals_rating": 4.5,
      "competencies_rating": 4,
      "values_rating": 4.5,
      "overall_rating": 4.3,
      "strengths": ["Consistent performer", "Team player"],
      "areas_for_improvement": ["Cross-team collaboration"],
      "comments": "Exceeded expectations in Q3-Q4.",
      "submitted_date": "2026-03-01",
      "submitted_by": "Jane Smith"
    },
    "final_rating": 4.3,
    "created_at": "2026-01-01"
  }
}
```

---

### 3.11 Submit Self Review

**`POST /api/v1/performance/appraisals/:id/self-review`**

**Auth:** Employee (own appraisal only)

**Request:**

```json
{
  "goals_rating": 4,
  "competencies_rating": 4.5,
  "values_rating": 4,
  "overall_rating": 4.2,
  "strengths": ["Technical leadership", "Problem solving"],
  "areas_for_improvement": ["Delegation", "Documentation"],
  "comments": "I've grown significantly this year."
}
```

| Field                    | Type       | Required | Description                 |
|--------------------------|------------|----------|-----------------------------|
| `goals_rating`           | `number`   | Yes      | Rating for goals (1-5)      |
| `competencies_rating`    | `number`   | Yes      | Rating for competencies     |
| `values_rating`          | `number`   | Yes      | Rating for values           |
| `overall_rating`         | `number`   | Yes      | Overall self rating         |
| `strengths`              | `string[]` | Yes      | List of strengths           |
| `areas_for_improvement`  | `string[]` | Yes      | List of improvement areas   |
| `comments`               | `string`   | No       | Free-text comments          |

**Response:** Updated appraisal (status advances to `Manager Review`).

---

### 3.12 Submit Manager Review

**`POST /api/v1/performance/appraisals/:id/manager-review`**

**Auth:** Manager (of the employee)

**Request:** Same shape as self review.

**Response:** Updated appraisal (status advances to `Calibration` or `Completed`).

---

### 3.13 Submit Calibration

**`POST /api/v1/performance/appraisals/:id/calibration`**

**Auth:** HR, Admin

**Request:**

```json
{
  "final_rating": 4.3,
  "comments": "Calibrated based on peer comparison."
}
```

**Response:** Updated appraisal (status = `Completed`).

---

### 3.14 Reopen / Send Back Appraisal

**Manager** or **HR** can send an appraisal back to a previous step.

**`PUT /api/v1/performance/appraisals/:id/send-back`**

**Auth:** Manager (own team), HR, Admin

**Request:**

```json
{
  "target_step": "Self Review",
  "reason": "Please revise your self-assessment for the Goals section."
}
```

| Field         | Type     | Required | Description                                |
|---------------|----------|----------|--------------------------------------------|
| `target_step` | `string` | Yes      | `Self Review`, `Manager Review`            |
| `reason`      | `string` | Yes      | Reason for sending back                    |

**Response:** Updated appraisal with status reverted.

---

### 3.15 Finalize Appraisal (HR/Admin)

After calibration, HR finalizes the appraisal with outcomes.

**`PUT /api/v1/performance/appraisals/:id/finalize`**

**Auth:** HR, Admin

**Request:**

```json
{
  "final_rating": 4.3,
  "rating_label": "Exceeds Expectations",
  "merit_increase": 8,
  "bonus_recommendation": 15,
  "promotion_flag": false,
  "development_actions": ["Leadership training", "Cross-team project lead"],
  "next_review_date": "2027-01-01"
}
```

| Field                  | Type       | Required | Description                         |
|------------------------|------------|----------|-------------------------------------|
| `final_rating`         | `number`   | Yes      | Calibrated final rating             |
| `rating_label`         | `string`   | Yes      | Label matching rating scale         |
| `merit_increase`       | `number`   | No       | Merit increase percentage           |
| `bonus_recommendation` | `number`   | No       | Bonus recommendation percentage     |
| `promotion_flag`       | `boolean`  | No       | Whether promotion is recommended    |
| `development_actions`  | `string[]` | No       | Development action items            |
| `next_review_date`     | `string`   | No       | ISO 8601 date for next review       |

**Response:** Updated appraisal (status = `Completed`) with full outcomes.

---

### Appraisal Status Workflow

```
                    ┌─────────────────┐
                    │   Not Started   │
                    └────────┬────────┘
                             │ Cycle activated
                    ┌────────▼────────┐
                    │   Self Review   │ ◄──── Send-back from Manager
                    │  (Employee)     │
                    └────────┬────────┘
                             │ Employee submits self-review
                    ┌────────▼────────┐
                    │ Manager Review  │ ◄──── Send-back from HR
                    │  (Manager)      │
                    └────────┬────────┘
                             │ Manager submits review
                    ┌────────▼────────┐
                    │  Calibration    │
                    │  (HR/Admin)     │
                    └────────┬────────┘
                             │ HR finalizes
                    ┌────────▼────────┐
                    │   Completed     │
                    └─────────────────┘
```

### Cycle Status Workflow

```
Draft → Active → Calibration → Completed → Closed
```

---

## Section 4: 360° Feedback

### 4.1 List 360° Feedback Campaigns

**`GET /api/v1/performance/feedback-360`**

**Auth:** HR/Admin (all), Manager (team), Employee (own)

**Query Parameters:**

| Param        | Type     | Required | Description                                  |
|-------------|----------|----------|----------------------------------------------|
| `status`    | `string` | No       | `Setup`, `In Progress`, `Completed`, `Shared`|
| `department`| `string` | No       | Filter by department                         |
| `page`      | `int`    | No       | Page number                                  |
| `page_size` | `int`    | No       | Items per page                               |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": "F360-001",
      "employee_name": "John Doe",
      "employee_id": 5,
      "designation": "Software Engineer",
      "department": "Engineering",
      "status": "Completed",
      "anonymous": true,
      "launch_date": "2026-01-05",
      "due_date": "2026-02-05",
      "total_raters": 12,
      "completed_responses": 10,
      "rater_groups": [
        { "group": "Self", "selected": 1, "completed": 1, "required": true },
        { "group": "Manager", "selected": 1, "completed": 1, "required": true },
        { "group": "Peers", "selected": 5, "completed": 4, "required": true },
        { "group": "Direct Reports", "selected": 3, "completed": 3, "required": false },
        { "group": "Others", "selected": 2, "completed": 1, "required": false }
      ],
      "overall_score": 4.1
    }
  ],
  "meta": { "page": 1, "per_page": 20, "total": 5, "total_pages": 1 }
}
```

---

### 4.2 Get 360° Feedback Detail

**`GET /api/v1/performance/feedback-360/:id`**

**Auth:** Employee (own, if Shared), Manager, HR, Admin

**Response:**

```json
{
  "success": true,
  "data": {
    "id": "F360-001",
    "employee_name": "John Doe",
    "employee_id": 5,
    "designation": "Software Engineer",
    "department": "Engineering",
    "status": "Completed",
    "anonymous": true,
    "launch_date": "2026-01-05",
    "due_date": "2026-02-05",
    "total_raters": 12,
    "completed_responses": 10,
    "overall_score": 4.1,
    "results": {
      "score_breakdown": {
        "self": 4.0,
        "manager": 4.5,
        "peers": 3.9,
        "direct_reports": 4.2
      },
      "competency_analysis": [
        {
          "competency": "Communication",
          "self_score": 4.0,
          "others_avg": 3.8,
          "gap": -0.2
        },
        {
          "competency": "Leadership",
          "self_score": 3.5,
          "others_avg": 4.2,
          "gap": 0.7
        }
      ],
      "key_strengths": ["Strategic thinking", "Mentoring"],
      "development_areas": ["Time management", "Delegation"],
      "verbatim_comments": [
        "Great at explaining complex topics",
        "Could improve on meeting deadlines"
      ]
    }
  }
}
```

---

### 4.3 Launch 360° Feedback

**`POST /api/v1/performance/feedback-360`**

**Auth:** HR, Admin

**Request:**

```json
{
  "employee_id": 5,
  "due_date": "2026-03-01",
  "anonymous": true,
  "min_raters": 3,
  "rater_groups": [
    { "group": "Self", "required": true },
    { "group": "Manager", "required": true },
    { "group": "Peers", "rater_ids": [10, 11, 12, 13, 14], "required": true },
    { "group": "Direct Reports", "rater_ids": [20, 21, 22], "required": false }
  ],
  "competencies": ["Communication", "Leadership", "Technical Skills", "Teamwork"]
}
```

**Response:** `201 Created` with campaign object.

---

### 4.4 Download 360° Report

**`GET /api/v1/performance/feedback-360/:id/report`**

**Auth:** HR, Admin, Employee (if Shared)

**Query:** `?format=pdf`

**Response:** Binary file download.

---

### 360° Feedback Status Workflow

```
Setup → In Progress → Completed → Shared
```

---

## Section 5: Talent Review (9-Box Grid)

### 5.1 List Talent Review Employees

**`GET /api/v1/performance/talent-review`**

**Auth:** Manager (team), HR, Admin

**Query Parameters:**

| Param           | Type     | Required | Description                          |
|-----------------|----------|----------|--------------------------------------|
| `department`    | `string` | No       | Filter by department                 |
| `box`           | `int`    | No       | Filter by 9-box position (1-9)       |
| `critical_role` | `bool`   | No       | Filter critical role employees       |
| `risk_of_loss`  | `string` | No       | `Low`, `Medium`, `High`             |

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "id": 5,
      "name": "John Doe",
      "designation": "Software Engineer",
      "department": "Engineering",
      "photo": "/avatars/john.jpg",
      "performance_rating": 4.2,
      "potential_rating": 4.5,
      "box": 2,
      "box_label": "High Potential",
      "critical_role": false,
      "tenure": "3 years",
      "risk_of_loss": "Low",
      "readiness": "6-12 Months",
      "development_priority": "High",
      "successor_for": ["Tech Lead"],
      "notes": "Strong candidate for promotion"
    }
  ]
}
```

### 9-Box Grid Layout

```
              Low Performance    Mid Performance    High Performance
High Pot   |  Box 7: Enigma   | Box 8: Growth Gem | Box 9: Future Star  |
Med  Pot   |  Box 4: Dilemma  | Box 5: Core Player| Box 6: High Potential|
Low  Pot   |  Box 1: Action   | Box 2: Up/Out     | Box 3: Solid Prof.  |
```

---

### 5.2 Get Talent Review Employee Detail

**`GET /api/v1/performance/talent-review/employees/:id`**

**Auth:** Manager, HR, Admin

**Response:** Full `TalentReviewEmployee` object.

---

### 5.3 Update Talent Review (Calibration)

**`PATCH /api/v1/performance/talent-review/employees/:id`**

**Auth:** HR, Admin

**Request:**

```json
{
  "performance_rating": 4.5,
  "potential_rating": 4.0,
  "box": 6,
  "readiness": "Now",
  "risk_of_loss": "Medium",
  "development_priority": "High",
  "critical_role": true,
  "successor_for": ["Tech Lead", "Engineering Manager"],
  "notes": "Updated after calibration session"
}
```

---

### 5.4 Start Calibration Session

**`POST /api/v1/performance/talent-review/calibration-sessions`**

**Auth:** HR, Admin

**Request:**

```json
{
  "name": "Q1 2026 Calibration",
  "department": "Engineering",
  "participant_ids": [3, 8, 15]
}
```

---

### 5.5 Add to Succession Plan

**`POST /api/v1/performance/talent-review/employees/:id/succession`**

**Auth:** HR, Admin

**Request:**

```json
{
  "role": "Engineering Manager",
  "readiness": "6-12 Months",
  "development_plan": "Leadership training + mentorship program"
}
```

---

### 5.6 Export Talent Review

**`GET /api/v1/performance/talent-review/export`**

**Auth:** HR, Admin

**Query:** `?format=xlsx&department=Engineering`

**Response:** Binary file download.

---

## Section 6: Performance Reports

### 6.1 Get Report Summary / Key Metrics

**`GET /api/v1/performance/reports/summary`**

**Auth:** HR, Admin

**Query Parameters:**

| Param        | Type     | Required | Description              |
|-------------|----------|----------|--------------------------|
| `department`| `string` | No       | Filter by department     |

**Response:**

```json
{
  "success": true,
  "data": {
    "avg_goal_completion": 78,
    "avg_performance_rating": 3.8,
    "avg_potential_rating": 3.5,
    "high_potential_count": 24,
    "total_employees": 150,
    "goal_completion_trend": "+5%"
  }
}
```

---

### 6.2 Get Rating Distribution

**`GET /api/v1/performance/reports/rating-distribution`**

**Auth:** HR, Admin

**Query:** `?department=Engineering`

**Response:**

```json
{
  "success": true,
  "data": [
    { "rating": 5, "label": "Outstanding", "count": 15, "percentage": 10 },
    { "rating": 4, "label": "Exceeds Expectations", "count": 45, "percentage": 30 },
    { "rating": 3, "label": "Meets Expectations", "count": 60, "percentage": 40 },
    { "rating": 2, "label": "Partially Meets", "count": 22, "percentage": 15 },
    { "rating": 1, "label": "Needs Improvement", "count": 8, "percentage": 5 }
  ]
}
```

---

### 6.3 Get Department Performance Summary

**`GET /api/v1/performance/reports/department-summary`**

**Auth:** HR, Admin

**Response:**

```json
{
  "success": true,
  "data": [
    {
      "department": "Engineering",
      "employee_count": 45,
      "avg_performance": 4.1,
      "avg_potential": 3.8,
      "high_potential_count": 12,
      "at_risk_count": 3
    },
    {
      "department": "Sales",
      "employee_count": 30,
      "avg_performance": 3.7,
      "avg_potential": 3.5,
      "high_potential_count": 8,
      "at_risk_count": 5
    }
  ]
}
```

---

### 6.4 Export Reports

**`GET /api/v1/performance/reports/export`**

**Auth:** HR, Admin

**Query Parameters:**

| Param        | Type     | Required | Description                                                              |
|-------------|----------|----------|--------------------------------------------------------------------------|
| `type`      | `string` | Yes      | `goal_completion`, `appraisal_compliance`, `high_potential`, `promotion_readiness`, `9box_distribution` |
| `department`| `string` | No       | Filter by department                                                     |
| `format`    | `string` | No       | `xlsx`, `pdf`, `csv` (default `xlsx`)                                    |

**Response:** Binary file download.

---

### 6.5 Individual Report Endpoints

| Method | Endpoint                                            | Description                    |
|--------|-----------------------------------------------------|--------------------------------|
| GET    | `/api/v1/performance/reports/goal-completion`       | Goal completion report data    |
| GET    | `/api/v1/performance/reports/appraisal-compliance`  | Appraisal compliance report    |
| GET    | `/api/v1/performance/reports/high-potential`         | High potential employee list   |
| GET    | `/api/v1/performance/reports/promotion-readiness`    | Promotion readiness report     |
| GET    | `/api/v1/performance/reports/9box-distribution`      | 9-Box distribution aggregates  |

All support `?department=` filter.

---

## Complete API Endpoint Summary

| #  | Method | Endpoint                                               | Page                  | Critical? |
|----|--------|--------------------------------------------------------|-----------------------|-----------|
| 1  | GET    | `/performance/dashboard/stats`                         | Dashboard             |           |
| 2  | GET    | `/performance/dashboard/upcoming-actions`              | Dashboard             |           |
| 3  | GET    | `/performance/goals`                                   | Goals & OKRs          |           |
| 4  | GET    | `/performance/goals/stats`                             | Goals & OKRs          |           |
| 5  | GET    | `/performance/goals/:id`                               | Goals & OKRs          |           |
| 6  | POST   | `/performance/goals`                                   | Goals & OKRs          |           |
| 7  | PATCH  | `/performance/goals/:id`                               | Goals & OKRs          |           |
| 8  | DELETE | `/performance/goals/:id`                               | Goals & OKRs          |           |
| 9  | POST   | `/performance/goals/:goalId/check-ins`                 | Goals & OKRs          |           |
| 10 | POST   | `/performance/goals/:goalId/key-results`               | Goals & OKRs          |           |
| 11 | PATCH  | `/performance/goals/:goalId/key-results/:krId`         | Goals & OKRs          |           |
| 12 | DELETE | `/performance/goals/:goalId/key-results/:krId`         | Goals & OKRs          |           |
| 13 | GET    | `/performance/goals/alignment`                         | Goals & OKRs          |           |
| 14 | POST   | `/performance/goals/:id/submit-for-approval`           | Goals & OKRs          | **NEW**   |
| 15 | PUT    | `/performance/goals/:id/approve`                       | Goals & OKRs          | **NEW**   |
| 16 | POST   | `/performance/goals/:id/request-completion`            | Goals & OKRs          | **NEW**   |
| 17 | PUT    | `/performance/goals/:id/verify-completion`             | Goals & OKRs          | **NEW**   |
| 18 | POST   | `/performance/goals/assign`                            | Goals & OKRs          | **NEW**   |
| 19 | GET    | `/performance/department-targets`                      | Department Targets    | **NEW**   |
| 20 | GET    | `/performance/department-targets/:id`                  | Department Targets    | **NEW**   |
| 21 | POST   | `/performance/department-targets`                      | Department Targets    | **NEW**   |
| 22 | PATCH  | `/performance/department-targets/:id`                  | Department Targets    | **NEW**   |
| 23 | POST   | `/performance/department-targets/:id/progress`         | Department Targets    | **NEW**   |
| 24 | PUT    | `/performance/department-targets/:id/milestones/:mid/complete` | Department Targets | **NEW** |
| 25 | POST   | `/performance/department-targets/:id/link-goal`        | Department Targets    | **NEW**   |
| 26 | DELETE | `/performance/department-targets/:id`                  | Department Targets    | **NEW**   |
| 27 | GET    | `/performance/appraisal-cycles`                        | Appraisals            |           |
| 28 | GET    | `/performance/appraisal-cycles/active`                 | Appraisals            |           |
| 29 | GET    | `/performance/appraisal-cycles/:id`                    | Appraisals            |           |
| 30 | POST   | `/performance/appraisal-cycles`                        | Appraisals            |           |
| 31 | PUT    | `/performance/appraisal-cycles/:id`                    | Appraisals            |           |
| 32 | PATCH  | `/performance/appraisal-cycles/:id/status`             | Appraisals            |           |
| 33 | DELETE | `/performance/appraisal-cycles/:id`                    | Appraisals            |           |
| 34 | GET    | `/performance/appraisals`                              | Appraisals            |           |
| 35 | GET    | `/performance/appraisals/summary`                      | Appraisals            |           |
| 36 | GET    | `/performance/appraisals/:id`                          | Appraisals            |           |
| 37 | POST   | `/performance/appraisals/:id/self-review`              | Appraisals            |           |
| 38 | POST   | `/performance/appraisals/:id/manager-review`           | Appraisals            |           |
| 39 | POST   | `/performance/appraisals/:id/calibration`              | Appraisals            |           |
| 40 | PUT    | `/performance/appraisals/:id/send-back`                | Appraisals            | **NEW**   |
| 41 | PUT    | `/performance/appraisals/:id/finalize`                 | Appraisals            | **NEW**   |
| 42 | GET    | `/performance/feedback-360`                            | 360° Feedback         |           |
| 43 | GET    | `/performance/feedback-360/:id`                        | 360° Feedback         |           |
| 44 | POST   | `/performance/feedback-360`                            | 360° Feedback         |           |
| 45 | GET    | `/performance/feedback-360/:id/report`                 | 360° Feedback         |           |
| 46 | GET    | `/performance/talent-review`                           | Talent Review         |           |
| 47 | GET    | `/performance/talent-review/employees/:id`             | Talent Review         |           |
| 48 | PATCH  | `/performance/talent-review/employees/:id`             | Talent Review         |           |
| 49 | POST   | `/performance/talent-review/calibration-sessions`      | Talent Review         |           |
| 50 | POST   | `/performance/talent-review/employees/:id/succession`  | Talent Review         |           |
| 51 | GET    | `/performance/talent-review/export`                    | Talent Review         |           |
| 52 | GET    | `/performance/reports/summary`                         | Reports               |           |
| 53 | GET    | `/performance/reports/rating-distribution`             | Reports               |           |
| 54 | GET    | `/performance/reports/department-summary`              | Reports               |           |
| 55 | GET    | `/performance/reports/export`                          | Reports               |           |
| 56 | GET    | `/performance/reports/goal-completion`                 | Reports               |           |
| 57 | GET    | `/performance/reports/appraisal-compliance`            | Reports               |           |
| 58 | GET    | `/performance/reports/high-potential`                  | Reports               |           |
| 59 | GET    | `/performance/reports/promotion-readiness`             | Reports               |           |
| 60 | GET    | `/performance/reports/9box-distribution`               | Reports               |           |
| 61 | POST   | `/performance/goals/:id/link-project`                  | Project Alignment     | **NEW**   |
| 62 | DELETE | `/performance/goals/:id/unlink-project`                | Project Alignment     | **NEW**   |
| 63 | POST   | `/performance/department-targets/:id/link-project`     | Project Alignment     | **NEW**   |
| 64 | DELETE | `/performance/department-targets/:id/unlink-project`   | Project Alignment     | **NEW**   |
| 65 | GET    | `/performance/project-alignment/:projectId`            | Project Alignment     | **NEW**   |
| 66 | GET    | `/performance/project-alignment/summary`               | Project Alignment     | **NEW**   |
| 67 | GET    | `/performance/employee-targets`                        | Employee Targets      | **NEW**   |
| 68 | GET    | `/performance/employee-targets/:id`                    | Employee Targets      | **NEW**   |
| 69 | POST   | `/performance/employee-targets`                        | Employee Targets      | **NEW**   |
| 70 | POST   | `/performance/employee-targets/bulk-assign`            | Employee Targets      | **NEW**   |
| 71 | POST   | `/performance/employee-targets/:id/progress`           | Employee Targets      | **NEW**   |
| 72 | PATCH  | `/performance/employee-targets/:id`                    | Employee Targets      | **NEW**   |
| 73 | DELETE | `/performance/employee-targets/:id`                    | Employee Targets      | **NEW**   |

**Total: 73 endpoints** (7 Employee Target endpoints + 6 Project Alignment + 8 Department Target + 7 critical workflow endpoints)

---

## Section 8: Project-Performance Alignment

This section defines how **Manage Performance** integrates with **Manage Project**. Projects and performance goals/targets are cross-linked so that project work directly contributes to employee performance measurement.

### Why Alignment Matters

- **Employees** can link their individual goals to projects they're working on, so project progress counts toward goal achievement
- **Department Leaders** can link department targets to projects, creating a clear chain: *Department Target → Project → Employee Goals*
- **Managers** can see which projects relate to which performance goals during appraisals
- **HR/Admin** can report on how projects contribute to overall organizational performance

### 8.1 Link Goal to Project

**`POST /api/v1/performance/goals/:id/link-project`**

**Auth:** Goal owner, Manager, HR, Admin

**Request:**

```json
{
  "project_id": 1,
  "project_code": "PRJ-0001",
  "project_name": "ERP System Migration",
  "contribution_notes": "Database module development contributes to this goal"
}
```

| Field                | Type     | Required | Description                                         |
|---------------------|----------|----------|-----------------------------------------------------|
| `project_id`        | `number` | Yes      | ID of the project from Manage Project module        |
| `project_code`      | `string` | No       | Project code (auto-fetched if not provided)         |
| `project_name`      | `string` | No       | Project name (auto-fetched if not provided)         |
| `contribution_notes`| `string` | No       | How this project contributes to the goal            |

**Response:**

```json
{
  "success": true,
  "message": "Project linked to goal successfully",
  "data": {
    "goal_id": "G-003",
    "linked_project": {
      "project_id": 1,
      "project_code": "PRJ-0001",
      "project_name": "ERP System Migration",
      "project_status": "active",
      "project_progress": 55
    }
  }
}
```

---

### 8.2 Unlink Goal from Project

**`DELETE /api/v1/performance/goals/:id/unlink-project`**

**Auth:** Goal owner, Manager, HR, Admin

**Request:**

```json
{
  "project_id": 1
}
```

**Response:**

```json
{
  "success": true,
  "message": "Project unlinked from goal"
}
```

---

### 8.3 Link Project to Department Target

**`POST /api/v1/performance/department-targets/:id/link-project`**

**Auth:** Target creator, Manager, HR, Admin

**Request:**

```json
{
  "project_id": 1,
  "project_code": "PRJ-0001",
  "project_name": "ERP System Migration"
}
```

**Response:**

```json
{
  "success": true,
  "message": "Project linked to department target",
  "data": {
    "target_id": "DT-001",
    "linked_project": {
      "project_id": 1,
      "project_code": "PRJ-0001",
      "project_name": "ERP System Migration",
      "project_status": "active",
      "project_progress": 55
    }
  }
}
```

---

### 8.4 Unlink Project from Department Target

**`DELETE /api/v1/performance/department-targets/:id/unlink-project`**

**Auth:** Target creator, Manager, HR, Admin

**Request:**

```json
{
  "project_id": 1
}
```

---

### 8.5 Get Performance Data for a Project

Returns all goals and department targets linked to a specific project.

**`GET /api/v1/performance/project-alignment/:projectId`**

**Auth:** Project member, Manager, HR, Admin

**Response:**

```json
{
  "success": true,
  "data": {
    "project": {
      "id": 1,
      "project_code": "PRJ-0001",
      "name": "ERP System Migration",
      "status": "active",
      "progress": 55
    },
    "linked_goals": [
      {
        "id": "G-002",
        "title": "Launch New Product Line",
        "owner": "John Mwita",
        "status": "In Progress",
        "progress": 65,
        "type": "OKR"
      },
      {
        "id": "G-003",
        "title": "Improve Code Quality",
        "owner": "Peter Ochieng",
        "status": "In Progress",
        "progress": 60,
        "type": "KPI"
      }
    ],
    "linked_targets": [
      {
        "id": "DT-001",
        "title": "Achieve 95% Sprint Velocity",
        "department": "Engineering",
        "status": "On Track",
        "progress": 88
      }
    ]
  }
}
```

---

### 8.6 Project Alignment Summary

Overview of all project-performance linkages for reporting.

**`GET /api/v1/performance/project-alignment/summary`**

**Auth:** HR, Admin

**Query Parameters:**

| Param        | Type     | Required | Description                        |
|-------------|----------|----------|------------------------------------|
| `department`| `string` | No       | Filter by department               |
| `period`    | `string` | No       | Filter by period                   |

**Response:**

```json
{
  "success": true,
  "data": {
    "total_linked_projects": 3,
    "total_linked_goals": 8,
    "total_linked_targets": 4,
    "projects": [
      {
        "project_id": 1,
        "project_code": "PRJ-0001",
        "project_name": "ERP System Migration",
        "linked_goals_count": 2,
        "linked_targets_count": 2,
        "avg_goal_progress": 62.5,
        "project_progress": 55
      }
    ]
  }
}
```

---

### Cross-Module Data Flow

```
┌─────────────────────────────────────────────────────┐
│                   Manage Project                     │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐ │
│  │ PRJ-0001    │  │ PRJ-0002    │  │ PRJ-0003    │ │
│  │ ERP System  │  │ Mobile App  │  │ Website     │ │
│  │ Progress:55%│  │ Progress:35%│  │ Progress:80%│ │
│  └──────┬──────┘  └──────┬──────┘  └─────────────┘ │
└─────────┼────────────────┼──────────────────────────┘
          │                │
          ▼                ▼
┌─────────────────────────────────────────────────────┐
│              Manage Performance                      │
│                                                      │
│  Department Targets:                                 │
│  ┌──────────────────────────────────────────────┐   │
│  │ DT-001: Achieve 95% Sprint Velocity          │   │
│  │ Linked Projects: PRJ-0001                    │   │
│  │ Linked Goals: G-002, G-003                   │   │
│  └──────────────────────────────────────────────┘   │
│  ┌──────────────────────────────────────────────┐   │
│  │ DT-002: Reduce Production Bugs by 40%        │   │
│  │ Linked Projects: PRJ-0001, PRJ-0002          │   │
│  │ Linked Goals: G-003                          │   │
│  └──────────────────────────────────────────────┘   │
│                                                      │
│  Individual Goals:                                   │
│  ┌──────────────────────────────────────────────┐   │
│  │ G-002: Launch New Product Line → PRJ-0001    │   │
│  │ G-003: Improve Code Quality   → PRJ-0001    │   │
│  └──────────────────────────────────────────────┘   │
│                                                      │
│  Appraisals:                                         │
│  ┌──────────────────────────────────────────────┐   │
│  │ Manager sees: "Employee worked on PRJ-0001"  │   │
│  │ Project progress visible in review context   │   │
│  └──────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────┘
```

### Integration Rules

| Scenario                                    | Behavior                                                    |
|--------------------------------------------|-------------------------------------------------------------|
| Goal linked to Project                     | Project progress shown alongside goal progress               |
| Department Target linked to Project        | Project status & progress visible in target detail           |
| Employee on a Project with linked Goal     | Daily task contributions feed into project → goal chain      |
| During Appraisal Review                    | Manager can see which projects the employee's goals link to  |
| Project completed                          | Linked goals flagged for completion review                   |
| Project cancelled                          | Linked goals flagged with warning status                     |

---

## Error Responses

All endpoints follow the standard error format:

```json
{
  "success": false,
  "message": "Error description",
  "errors": []
}
```

| Status | Description                          |
|--------|--------------------------------------|
| 400    | Invalid request / validation error   |
| 401    | Not authenticated                    |
| 403    | Insufficient permissions             |
| 404    | Resource not found                   |
| 409    | Conflict (e.g. duplicate)            |
| 500    | Internal server error                |
