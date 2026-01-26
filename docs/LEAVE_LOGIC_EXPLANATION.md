# Leave Management Logic Explanation

## Current Status

**🚧 Leave Management Module is NOT YET IMPLEMENTED**

The leave management system is currently a **placeholder** for future development. However, there are some foundational elements in place that reference leave functionality.

---

## What Exists Currently

### 1. **Employee Status - "on_leave"**

The employee model supports an `on_leave` status, which can be set on employees:

**Location:** `internal/modules/employees/models/employee.go`

```go
type EmployeeStatus string

const (
    StatusActive     EmployeeStatus = "active"
    StatusInactive   EmployeeStatus = "inactive"
    StatusOnLeave    EmployeeStatus = "on_leave"  // ✅ Exists
    StatusSuspended  EmployeeStatus = "suspended"
    StatusTerminated EmployeeStatus = "terminated"
    StatusArchived   EmployeeStatus = "archived"
)
```

**What this means:**
- Employees can have their status set to `"on_leave"` 
- This is just a status flag - there's no leave request/approval workflow
- No leave balance tracking
- No leave type management

---

### 2. **Employee Policy - Leave Policy Reference**

During employee onboarding, there's a field to reference a leave policy:

**Location:** `internal/modules/employees/models/employee_policy.go`

```go
type EmployeePolicy struct {
    ID                uint           `json:"id"`
    EmployeeID       *uint          `json:"employee_id,omitempty"`
    LeavePolicyID    *uint          `json:"leave_policy_id,omitempty"`  // ⚠️ References leave_policies(id) - to be created
    AttendancePolicyID *uint        `json:"attendance_policy_id,omitempty"`
    WeeklyOffDays    *string        `json:"weekly_off_days,omitempty"`
    EffectiveDate    *time.Time     `json:"effective_date,omitempty"`
    // ...
}
```

**What this means:**
- The `employee_policies` table has a `leave_policy_id` column
- This is meant to reference a `leave_policies` table that **doesn't exist yet**
- During onboarding (Step 8), you can assign a `LeavePolicyID` to an employee
- But since the leave policies table doesn't exist, this is just storing a reference that can't be used

**Onboarding Step 8 - Policies:**
```go
type Step8PoliciesRequest struct {
    LeavePolicyID     *uint      `json:"leave_policy_id,omitempty"`  // Can be set, but no validation
    AttendancePolicyID *uint     `json:"attendance_policy_id,omitempty"`
    WeeklyOffDays     *string    `json:"weekly_off_days,omitempty"`
    EffectiveDate     *time.Time `json:"effective_date,omitempty"`
}
```

---

### 3. **Leave Encashment in Offboarding**

When an employee is offboarded, there's a field for leave encashment in the final settlement:

**Location:** `internal/modules/employees/models/offboarding_requests.go`

```go
type CalculateSettlementRequest struct {
    OutstandingSalary *float64 `json:"outstanding_salary,omitempty"`
    LeaveEncashment   *float64 `json:"leave_encashment,omitempty"`  // ✅ Exists
    Bonus             *float64 `json:"bonus,omitempty"`
    // ...
}
```

**What this means:**
- During offboarding, HR can manually enter a leave encashment amount
- This is a manual field - there's no automatic calculation based on leave balance
- No leave balance is tracked or deducted

---

## What's Missing (Planned Features)

According to `internal/modules/leave/README.md`, the following features are planned but **not implemented**:

1. ❌ **Leave Types Management** - No leave types (Annual, Sick, Casual, etc.)
2. ❌ **Leave Balance Tracking** - No tracking of leave balances per employee
3. ❌ **Leave Requests** - No API to request leave
4. ❌ **Leave Approvals** - No approval workflow
5. ❌ **Leave Calendar** - No calendar view of leaves
6. ❌ **Leave Reports** - No reporting functionality

---

## Current Workflow (What Actually Works)

### Setting Employee to "on_leave" Status

You can update an employee's status to `"on_leave"`:

```bash
PUT /api/v1/employees/:employee_id
{
  "status": "on_leave"
}
```

**Limitations:**
- This is just a status change - no leave request/approval
- No leave balance is checked or deducted
- No leave type is specified
- No start/end dates are tracked
- No automatic return to "active" status

### Assigning Leave Policy During Onboarding

During employee onboarding (Step 8), you can set a `LeavePolicyID`:

```bash
PUT /api/v1/employees/onboarding/drafts/:draft_id/steps/8
{
  "leave_policy_id": 1,  // ⚠️ This ID doesn't reference anything yet
  "attendance_policy_id": 1,
  "weekly_off_days": "Saturday,Sunday",
  "effective_date": "2026-01-01T00:00:00Z"
}
```

**Limitations:**
- The `leave_policy_id` is stored but not validated
- No leave policies table exists to reference
- This is just storing a placeholder value

### Leave Encashment During Offboarding

When calculating final settlement during offboarding:

```bash
POST /api/v1/employees/offboarding/:offboarding_id/calculate-settlement
{
  "leave_encashment": 500000.00,  // Manual entry
  "outstanding_salary": 1000000.00,
  // ...
}
```

**Limitations:**
- This is a manual field - no automatic calculation
- No leave balance is checked
- HR must manually calculate and enter the amount

---

## Database Schema

### Current Tables

1. **`employee_policies`** table exists with:
   - `leave_policy_id` (nullable, references non-existent `leave_policies` table)
   - `attendance_policy_id` (nullable)
   - `weekly_off_days` (string)
   - `effective_date` (timestamp)

2. **`employees`** table has:
   - `status` field that can be set to `"on_leave"`

### Missing Tables

The following tables need to be created for full leave management:

1. **`leave_policies`** - Define leave types and rules
2. **`leave_balances`** - Track leave balance per employee
3. **`leave_requests`** - Store leave requests
4. **`leave_approvals`** - Track approval workflow
5. **`leave_types`** - Define types (Annual, Sick, Casual, etc.)

---

## Summary

### ✅ What Works:
- Setting employee status to `"on_leave"` (just a status flag)
- Storing `LeavePolicyID` during onboarding (but no policies exist)
- Manual leave encashment entry during offboarding

### ❌ What Doesn't Work:
- Leave request submission
- Leave approval workflow
- Leave balance tracking
- Leave type management
- Automatic leave balance deduction
- Leave calendar
- Leave reports

### 🔧 To Implement Full Leave Management:

1. **Create Leave Policies Module:**
   - Define leave types (Annual, Sick, Casual, etc.)
   - Set accrual rules
   - Set carry-forward rules
   - Set maximum balances

2. **Create Leave Requests Module:**
   - Request leave (start date, end date, type, reason)
   - Approval workflow (Manager → HR)
   - Status tracking (pending, approved, rejected)

3. **Create Leave Balance Tracking:**
   - Track balance per employee per leave type
   - Automatic accrual
   - Deduction on approval
   - Carry-forward at year-end

4. **Create Leave Calendar:**
   - View all leaves in calendar format
   - Filter by employee, department, leave type

5. **Create Leave Reports:**
   - Leave utilization reports
   - Balance reports
   - Pending requests reports

---

## Recommendation

If you need leave management functionality, you should:

1. **Design the leave management module** following the same pattern as Shifts & Rosters
2. **Create the database models** for leave policies, requests, balances
3. **Implement the API endpoints** for leave requests and approvals
4. **Add leave balance tracking** with automatic accrual and deduction
5. **Integrate with employee status** to automatically set status to "on_leave" when leave is approved

The current system only provides the foundation (employee status and policy references) but no actual leave management functionality.
