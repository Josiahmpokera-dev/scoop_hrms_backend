# Performance Management - Implementation Documentation

## Overview

The Performance Management module has been fully implemented with **73 API endpoints** across 10 sections. All endpoints are under the base URL `/api/v1/performance` and require authentication via `AuthMiddleware`.

---

## Architecture

```
internal/modules/performance/
├── models/
│   ├── goal.go                  # Goal, KeyResult, GoalCheckIn + DTOs
│   ├── department_target.go     # DepartmentTarget, Milestone + DTOs
│   ├── employee_target.go       # EmployeeTarget + DTOs
│   ├── appraisal.go            # AppraisalCycle, WorkflowStep, Appraisal + DTOs
│   ├── feedback360.go          # Feedback360Campaign, RaterGroup + DTOs
│   └── talent_review.go        # TalentReview, CalibrationSession, SuccessionPlan + DTOs
├── repositories/
│   ├── goal_repository.go
│   ├── department_target_repository.go
│   ├── employee_target_repository.go
│   ├── appraisal_repository.go
│   ├── feedback360_repository.go
│   └── talent_review_repository.go
├── services/
│   ├── goal_service.go          # + Dashboard stats & upcoming actions
│   ├── department_target_service.go
│   ├── employee_target_service.go
│   ├── appraisal_service.go
│   ├── feedback360_service.go
│   ├── talent_review_service.go
│   └── report_service.go
└── handlers/
    ├── goal_handler.go          # + Dashboard endpoints
    ├── department_target_handler.go
    ├── employee_target_handler.go
    ├── appraisal_handler.go
    ├── feedback360_handler.go
    ├── talent_review_handler.go
    └── report_handler.go
```

### Database Tables (14 tables, auto-migrated)

| Table | Description |
|-------|-------------|
| `performance_goals` | Performance goals (OKR, KPI, Project, Development) |
| `performance_key_results` | Measurable key results for OKR goals |
| `performance_goal_check_ins` | Goal progress check-ins |
| `performance_department_targets` | Department-level KPIs and objectives |
| `performance_department_target_milestones` | Milestones for department targets |
| `performance_employee_targets` | Individual assigned targets |
| `performance_appraisal_cycles` | Appraisal review cycles |
| `performance_appraisal_workflow_steps` | Workflow steps within a cycle |
| `performance_appraisals` | Individual employee appraisals |
| `performance_feedback_360_campaigns` | 360-degree feedback campaigns |
| `performance_feedback_360_rater_groups` | Rater groups within campaigns |
| `performance_talent_reviews` | Talent review / 9-box entries |
| `performance_calibration_sessions` | HR calibration sessions |
| `performance_succession_plans` | Succession planning entries |

---

## API Endpoints - Quick Reference

### Section 1: Dashboard (2 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/dashboard/stats` | Any | Dashboard statistics (goals, appraisal, feedback, quick stats) |
| GET | `/performance/dashboard/upcoming-actions` | Any | Upcoming action items |

**Example: GET `/api/v1/performance/dashboard/stats`**
```json
{
  "success": true,
  "data": {
    "goals": { "total": 8, "completed": 3, "in_progress": 4, "at_risk": 1, "avg_progress": 57 },
    "appraisal": { "status": "Self Review", "self_rating": 4.2 },
    "feedback_360": { "total": 5, "available": 2 },
    "quick_stats": { "goals_on_track": 4, "review_completion_percent": 63 }
  }
}
```

---

### Section 2: Goals & OKRs (18 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/goals` | Any | List goals (filtered by role) |
| GET | `/performance/goals/stats` | Any | Goal statistics |
| GET | `/performance/goals/alignment` | Manager+ | Goal alignment map (tree) |
| GET | `/performance/goals/:id` | Owner/Manager/HR | Get goal detail |
| POST | `/performance/goals` | Any | Create goal |
| PATCH | `/performance/goals/:id` | Owner/Manager/HR | Update goal |
| DELETE | `/performance/goals/:id` | Owner/HR/Admin | Delete goal |
| POST | `/performance/goals/:id/submit-for-approval` | Employee | Submit goal for approval |
| PUT | `/performance/goals/:id/approve` | Manager/HR/Admin | Approve or reject goal |
| POST | `/performance/goals/:id/request-completion` | Employee | Request completion verification |
| PUT | `/performance/goals/:id/verify-completion` | Manager/HR/Admin | Verify or reject completion |
| POST | `/performance/goals/assign` | Manager/HR/Admin | Assign goal to employee |
| POST | `/performance/goals/:id/link-project` | Owner/Manager/HR | Link goal to project |
| DELETE | `/performance/goals/:id/unlink-project` | Owner/Manager/HR | Unlink goal from project |
| POST | `/performance/goals/:goalId/key-results` | Owner/Manager | Create key result |
| PATCH | `/performance/goals/:goalId/key-results/:krId` | Owner/Manager | Update key result |
| DELETE | `/performance/goals/:goalId/key-results/:krId` | Owner/Manager | Delete key result |
| POST | `/performance/goals/:goalId/check-ins` | Owner/Manager | Create check-in |

**Example: POST `/api/v1/performance/goals`**
```json
{
  "title": "Improve customer satisfaction score",
  "description": "Raise NPS from 45 to 60",
  "type": "OKR",
  "level": "Individual",
  "weight": 30,
  "due_date": "2026-12-31",
  "visibility": "Manager",
  "tags": ["Customer Success"]
}
```

**Example: POST `/api/v1/performance/goals/:id/submit-for-approval`**
```json
{ "comments": "Please review my Q2 goal" }
```

**Example: PUT `/api/v1/performance/goals/:id/approve`**
```json
{ "action": "approve", "comments": "Approved. Well-defined targets." }
```

**Example: PUT `/api/v1/performance/goals/:id/verify-completion`**
```json
{ "action": "verify", "final_rating": 4, "comments": "Well done." }
```

---

### Section 2B: Department Targets (10 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/department-targets` | Any | List targets |
| GET | `/performance/department-targets/:id` | Any | Get target detail |
| POST | `/performance/department-targets` | Manager+ | Create target |
| PATCH | `/performance/department-targets/:id` | Manager+ | Update target |
| DELETE | `/performance/department-targets/:id` | Manager+ | Delete target |
| POST | `/performance/department-targets/:id/progress` | Manager+ | Update progress |
| PUT | `/performance/department-targets/:id/milestones/:milestoneId/complete` | Manager+ | Complete milestone |
| POST | `/performance/department-targets/:id/link-goal` | Manager+ | Link goal to target |
| POST | `/performance/department-targets/:id/link-project` | Manager+ | Link project to target |
| DELETE | `/performance/department-targets/:id/unlink-project` | Manager+ | Unlink project |

**Example: POST `/api/v1/performance/department-targets`**
```json
{
  "title": "Reduce Production Bugs by 40%",
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
    { "title": "Implement code review policy", "due_date": "2026-01-31" }
  ]
}
```

---

### Section 2C: Employee Targets (7 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/employee-targets` | Any | List targets |
| GET | `/performance/employee-targets/:id` | Any | Get target detail |
| POST | `/performance/employee-targets` | Manager+ | Assign target |
| POST | `/performance/employee-targets/bulk-assign` | Manager+ | Bulk assign targets |
| PATCH | `/performance/employee-targets/:id` | Manager+ | Update target |
| DELETE | `/performance/employee-targets/:id` | Manager+ | Delete target |
| POST | `/performance/employee-targets/:id/progress` | Any | Update progress |

**Example: POST `/api/v1/performance/employee-targets`**
```json
{
  "employee_id": "EMP-003",
  "title": "Deliver 3 API Modules by End of Q1",
  "metric": "Modules Delivered",
  "target_value": 3,
  "unit": "modules",
  "weight": 30,
  "due_date": "2026-03-31",
  "period": "Q1 2026",
  "priority": "High"
}
```

**Example: POST `/api/v1/performance/employee-targets/bulk-assign`**
```json
{
  "employee_ids": ["EMP-003", "EMP-004", "EMP-008"],
  "title": "Achieve 90% Code Review Compliance",
  "target_value": 90,
  "unit": "%",
  "weight": 20,
  "due_date": "2026-06-30",
  "period": "H1 2026",
  "priority": "Medium"
}
```

---

### Section 3: Appraisal Cycles (7 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/appraisal-cycles` | Any | List cycles |
| GET | `/performance/appraisal-cycles/active` | Any | Get active cycle |
| GET | `/performance/appraisal-cycles/:id` | Any | Get cycle detail |
| POST | `/performance/appraisal-cycles` | HR/Admin | Create cycle |
| PUT | `/performance/appraisal-cycles/:id` | HR/Admin | Update cycle |
| PATCH | `/performance/appraisal-cycles/:id/status` | HR/Admin | Change status |
| DELETE | `/performance/appraisal-cycles/:id` | Admin | Delete cycle (Draft only) |

**Example: POST `/api/v1/performance/appraisal-cycles`**
```json
{
  "name": "Annual Review 2026",
  "cycle_type": "Annual",
  "year": 2026,
  "start_date": "2026-01-01",
  "end_date": "2026-12-31",
  "steps": [
    { "step_order": 1, "name": "Self Assessment", "due_date": "2026-02-28" },
    { "step_order": 2, "name": "Manager Review", "due_date": "2026-03-15" },
    { "step_order": 3, "name": "HR Calibration", "due_date": "2026-03-31" }
  ]
}
```

**Status transitions:** `Draft -> Active -> Calibration -> Completed -> Closed`

---

### Section 3B: Appraisals (8 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/appraisals` | Any | List appraisals |
| GET | `/performance/appraisals/summary` | Any | Appraisal summary |
| GET | `/performance/appraisals/:id` | Owner/Manager/HR | Get detail |
| POST | `/performance/appraisals/:id/self-review` | Employee | Submit self review |
| POST | `/performance/appraisals/:id/manager-review` | Manager | Submit manager review |
| POST | `/performance/appraisals/:id/calibration` | HR/Admin | Submit calibration |
| PUT | `/performance/appraisals/:id/send-back` | Manager/HR | Send back to previous step |
| PUT | `/performance/appraisals/:id/finalize` | HR/Admin | Finalize with outcomes |

**Example: POST `/api/v1/performance/appraisals/:id/self-review`**
```json
{
  "self_rating": 4.2,
  "strengths": ["Technical leadership", "Problem solving"],
  "areas_for_improvement": ["Delegation"],
  "comments": "Great year of growth."
}
```

**Appraisal Flow:** `Not Started -> Self Review -> Manager Review -> Calibration -> Completed`

---

### Section 4: 360-Degree Feedback (3 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/feedback-360` | Any | List campaigns |
| GET | `/performance/feedback-360/:id` | Any | Get campaign detail |
| POST | `/performance/feedback-360` | HR/Admin | Launch campaign |

**Example: POST `/api/v1/performance/feedback-360`**
```json
{
  "employee_id": 5,
  "due_date": "2026-03-01",
  "anonymous": true,
  "min_raters": 3,
  "rater_groups": [
    { "group": "Self", "required": true },
    { "group": "Manager", "required": true },
    { "group": "Peers", "rater_ids": "10,11,12", "required": true }
  ],
  "competencies": "Communication,Leadership,Technical Skills"
}
```

---

### Section 5: Talent Review / 9-Box (5 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/talent-review` | Manager+ | List talent reviews |
| GET | `/performance/talent-review/employees/:id` | Manager+ | Get review detail |
| PATCH | `/performance/talent-review/employees/:id` | HR/Admin | Update review |
| POST | `/performance/talent-review/calibration-sessions` | HR/Admin | Create calibration session |
| POST | `/performance/talent-review/employees/:id/succession` | HR/Admin | Add succession plan |

**9-Box Grid:**
```
              Low Performance    Mid Performance    High Performance
High Pot   |  Box 7: Enigma   | Box 8: Growth Gem | Box 9: Future Star  |
Med  Pot   |  Box 4: Dilemma  | Box 5: Core Player| Box 6: High Potential|
Low  Pot   |  Box 1: Action   | Box 2: Up/Out     | Box 3: Solid Prof.  |
```

Box is auto-computed from `performance_rating` and `potential_rating`:
- Low: 1.0 - 2.5
- Mid: 2.5 - 3.5
- High: 3.5 - 5.0

---

### Section 6: Reports (3 endpoints)

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| GET | `/performance/reports/summary` | HR/Admin | Summary with key metrics |
| GET | `/performance/reports/rating-distribution` | HR/Admin | Rating distribution chart data |
| GET | `/performance/reports/department-summary` | HR/Admin | Department-level summary |

All report endpoints support `?department=` filter.

---

## Role Access Matrix

| API Category | Employee | Manager | HR | Admin |
|---|---|---|---|---|
| Own goals & self review | Yes | Yes | Yes | Yes |
| Team goals & reviews | No | Yes | Yes | Yes |
| Department targets (create/edit) | No | Yes | Yes | Yes |
| Employee targets (assign) | No | Yes | Yes | Yes |
| Appraisal cycles mgmt | No | No | Yes | Yes |
| 360 Feedback mgmt | No | No | Yes | Yes |
| Talent review / 9-Box | No | Yes | Yes | Yes |
| Performance reports | No | No | Yes | Yes |

---

## Middleware Applied

| Middleware | Used For |
|---|---|
| `AuthMiddleware()` | All performance endpoints (base group) |
| `ManagerMiddleware()` | Department targets CRUD, Employee targets assign, Talent review list/detail |
| `HRMiddleware()` | Appraisal cycles create/update/status, Calibration, Finalize, 360 launch, Talent review update, Reports |
| `AdminMiddleware()` | Delete appraisal cycle (Draft only) |

---

## Error Responses

All endpoints follow the standard format:

```json
{
  "success": false,
  "message": "Error description",
  "errors": "optional details"
}
```

| Status | Description |
|--------|-------------|
| 400 | Invalid request / validation error |
| 401 | Not authenticated |
| 403 | Insufficient permissions |
| 404 | Resource not found |
| 500 | Internal server error |

---

## Goal Workflow

```
Not Started -> Pending Approval -> Approved -> In Progress -> Pending Verification -> Completed (Verified)
                                -> Rejected (revise & resubmit)
                                                              -> Rejected (back to In Progress)
Any status -> Cancelled
```

## Appraisal Workflow

```
Not Started -> Self Review (Employee) -> Manager Review (Manager) -> Calibration (HR) -> Completed
                                      <- Send-back                <- Send-back
```

## Cycle Status Workflow

```
Draft -> Active -> Calibration -> Completed -> Closed
```

---

## Auto-Generated Codes

| Entity | Format | Example |
|--------|--------|---------|
| Goal | G-XXX | G-001, G-002 |
| Key Result | KR-XXX | KR-001, KR-002 |
| Department Target | DT-XXX | DT-001, DT-002 |
| Employee Target | ET-XXX | ET-001, ET-002 |
| Appraisal Cycle | AC-XXX | AC-001, AC-002 |
| Appraisal | APR-XXX | APR-001, APR-002 |
| 360 Feedback | F360-XXX | F360-001, F360-002 |

---

## Alignment Chain

```
Company Goal (e.g., "Increase Revenue by 30%")
  └── Department Target (e.g., "Achieve 95% Sprint Velocity")
        ├── Employee Target (e.g., "Deliver 3 API Modules")
        ├── Employee Target (e.g., "Complete Dashboard")
        └── Linked Project (e.g., "PRJ-0001 ERP Migration")
              └── Individual Goals linked to same project
```
