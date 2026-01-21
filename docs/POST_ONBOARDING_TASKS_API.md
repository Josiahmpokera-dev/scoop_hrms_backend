# Post-Onboarding Tasks API Documentation

## Overview

The Post-Onboarding Tasks API allows you to track and manage tasks that need to be completed after an employee has been onboarded. These tasks include activities like sending welcome emails, creating email accounts, allocating equipment, conducting reviews, and more.

**Base URL:** `/api/v1/employees/post-onboarding`

**Authentication:** All endpoints require authentication via JWT token in the `Authorization` header.

---

## Available Task Types

The system supports 14 predefined task types:

1. **`send_welcome_email`** - Send Welcome Email
2. **`create_email_account`** - Create Email Account
3. **`allocate_laptop`** - Allocate Laptop
4. **`desk_allocate`** - Desk Allocation
5. **`issue_id_card`** - Issue ID Card
6. **`team_introduction`** - Team Introduction
7. **`policy_acknowledge`** - Policy Acknowledgment
8. **`induction_training`** - Induction Training
9. **`assign_buddy_mentor`** - Assign Buddy/Mentor
10. **`setup_dev_environment`** - Setup Development Environment
11. **`first_week_review`** - First Week Review
12. **`first_month_review`** - First Month Review
13. **`mid_probation_review`** - Mid-Probation Review
14. **`probation_confirmation`** - Probation Confirmation Assessment

### Task Status Values
- `pending` - Task is pending
- `in_progress` - Task is in progress
- `completed` - Task is completed
- `skipped` - Task was skipped

### Task Priority Values
- `low` - Low priority
- `medium` - Medium priority (default)
- `high` - High priority
- `critical` - Critical priority

---

## API Endpoints

### 1. Get Post-Onboarding Statistics

Get overall statistics for post-onboarding tasks including active onboarding, total tasks, completed tasks, and overdue tasks.

**Endpoint:** `GET /api/v1/employees/post-onboarding/statistics`

**Authentication:** Required

**Request:** No body or query parameters required

**Response:**
```json
{
  "success": true,
  "message": "Post-onboarding statistics retrieved successfully",
  "data": {
    "active_onboarding": 12,
    "total_tasks": 145,
    "completed_tasks": 98,
    "overdue_tasks": 5,
    "completion_percentage": 67.6,
    "tasks_by_status": {
      "pending": 25,
      "in_progress": 22,
      "completed": 98,
      "skipped": 0
    }
  }
}
```

**Response Fields:**
- `active_onboarding` - Number of employees with active (incomplete) onboarding tasks
- `total_tasks` - Total number of post-onboarding tasks
- `completed_tasks` - Number of completed tasks
- `overdue_tasks` - Number of overdue tasks (pending/in_progress tasks with due_date < today)
- `completion_percentage` - Overall completion percentage (rounded to 1 decimal place)
- `tasks_by_status` - Breakdown of tasks by status (pending, in_progress, completed, skipped)

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/statistics
Authorization: Bearer {{TOKEN}}
```

**Use Cases:**
- Dashboard overview of post-onboarding progress
- Track overall completion rates
- Identify overdue tasks that need attention
- Monitor active onboarding employees

---

### 2. List All Employees with Task Completion

Get a paginated list of all employees with their post-onboarding task completion percentages.

**Endpoint:** `GET /api/v1/employees/post-onboarding/employees`

**Authentication:** Required

**Query Parameters:**
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)

**Response:**
```json
{
  "success": true,
  "message": "Employees with task completion retrieved successfully",
  "data": [
    {
      "employee_id": "EMP567073",
      "first_name": "John",
      "last_name": "Doe",
      "email": "john.doe@company.com",
      "department_id": 5,
      "position_id": 10,
      "status": "active",
      "hire_date": "2024-01-15T00:00:00Z",
      "total_tasks": 14,
      "completed_tasks": 10,
      "pending_tasks": 3,
      "in_progress_tasks": 1,
      "skipped_tasks": 0,
      "completion_percentage": 71.43,
      "is_completed": false
    },
    {
      "employee_id": "EMP567074",
      "first_name": "Jane",
      "last_name": "Smith",
      "email": "jane.smith@company.com",
      "department_id": 5,
      "position_id": 11,
      "status": "active",
      "hire_date": "2024-01-20T00:00:00Z",
      "total_tasks": 14,
      "completed_tasks": 14,
      "pending_tasks": 0,
      "in_progress_tasks": 0,
      "skipped_tasks": 0,
      "completion_percentage": 100.0,
      "is_completed": true
    },
    {
      "employee_id": "EMP567075",
      "first_name": "Bob",
      "last_name": "Johnson",
      "email": "bob.johnson@company.com",
      "department_id": 6,
      "position_id": 12,
      "status": "active",
      "hire_date": "2024-02-01T00:00:00Z",
      "total_tasks": 0,
      "completed_tasks": 0,
      "pending_tasks": 0,
      "in_progress_tasks": 0,
      "skipped_tasks": 0,
      "completion_percentage": 0.0,
      "is_completed": false
    }
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/employees?page=1&page_size=20
Authorization: Bearer {{TOKEN}}
```

**Response Fields:**
- `employee_id` - Employee ID (e.g., EMP567073)
- `first_name` - Employee's first name
- `last_name` - Employee's last name
- `email` - Employee's work email
- `department_id` - Department ID
- `position_id` - Position ID
- `status` - Employee status (active, inactive, etc.)
- `hire_date` - Date when employee was hired
- `total_tasks` - Total number of post-onboarding tasks assigned
- `completed_tasks` - Number of completed tasks
- `pending_tasks` - Number of pending tasks
- `in_progress_tasks` - Number of in-progress tasks
- `skipped_tasks` - Number of skipped tasks
- `completion_percentage` - Percentage of tasks completed (0-100)
- `is_completed` - Boolean indicating if all tasks are completed (100%)

**Use Cases:**
- View all employees and their post-onboarding progress at a glance
- Identify employees who haven't completed all tasks
- Track completion rates across the organization
- Filter employees by completion status in the frontend

---

### 3. Get Available Task Types

Get a list of all available task types with their descriptions.

**Endpoint:** `GET /api/v1/employees/post-onboarding/types`

**Authentication:** Required

**Request:** No body required

**Response:**
```json
{
  "success": true,
  "message": "Task types retrieved successfully",
  "data": [
    {
      "type": "send_welcome_email",
      "name": "Send Welcome Email",
      "description": "Send welcome email to new employee",
      "default_assigned_to": "HR",
      "is_conditional": false
    },
    {
      "type": "create_email_account",
      "name": "Create Email Account",
      "description": "Create company email account for employee",
      "default_assigned_to": "HR",
      "is_conditional": false
    },
    {
      "type": "allocate_laptop",
      "name": "Allocate Laptop",
      "description": "Allocate laptop/computer to employee",
      "default_assigned_to": "IT",
      "is_conditional": false
    }
    // ... more task types
  ]
}
```

---

### 4. Create a Single Task

Create a new post-onboarding task for an employee.

**Endpoint:** `POST /api/v1/employees/post-onboarding/tasks`

**Authentication:** Required

**Request Body:**
```json
{
  "employee_id": "EMP567073",
  "task_type": "send_welcome_email",
  "title": "Send Welcome Email to John Doe",  // Optional - auto-generated if not provided
  "description": "Send welcome email with company information",  // Optional
  "priority": "high",  // Optional: low, medium, high, critical (default: medium)
  "assigned_to": "HR",  // Optional: department/role (default based on task type)
  "assigned_to_user_id": 5,  // Optional: specific user ID
  "due_date": "2024-02-15",  // Optional: YYYY-MM-DD format
  "notes": "Priority task - send within first day",  // Optional
  "metadata": "{\"email_template\": \"welcome_v1\"}"  // Optional: JSON string for task-specific data
}
```

**Response:**
```json
{
  "success": true,
  "message": "Post-onboarding task created successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP567073",
    "task_type": "send_welcome_email",
    "title": "Send Welcome Email to John Doe",
    "description": "Send welcome email with company information",
    "status": "pending",
    "priority": "high",
    "assigned_to": "HR",
    "assigned_to_user_id": 5,
    "assigned_to_user_name": "Jane Smith",
    "completed_by_user_id": null,
    "completed_by_user_name": null,
    "completed_at": null,
    "due_date": "2024-02-15T00:00:00Z",
    "notes": "Priority task - send within first day",
    "metadata": "{\"email_template\": \"welcome_v1\"}",
    "created_by": 1,
    "updated_by": 1,
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-20T10:00:00Z"
  }
}
```

**Example: Create "Send Welcome Email" Task**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "employee_id": "EMP567073",
  "task_type": "send_welcome_email",
  "priority": "high",
  "assigned_to": "HR",
  "due_date": "2024-02-15"
}
```

**Example: Create "Create Email Account" Task**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "employee_id": "EMP567073",
  "task_type": "create_email_account",
  "priority": "critical",
  "assigned_to": "IT",
  "assigned_to_user_id": 3,
  "due_date": "2024-02-14",
  "notes": "Create email account with format: firstname.lastname@company.com"
}
```

---

### 5. Bulk Create Tasks

Create multiple tasks for an employee at once.

**Endpoint:** `POST /api/v1/employees/post-onboarding/tasks/bulk`

**Authentication:** Required

**Request Body:**
```json
{
  "employee_id": "EMP567073",
  "task_types": [
    "send_welcome_email",
    "create_email_account",
    "allocate_laptop",
    "issue_id_card"
  ],
  "priority": "medium",  // Optional: default priority for all tasks
  "due_date": "2024-02-20"  // Optional: default due date for all tasks
}
```

**Response:**
```json
{
  "success": true,
  "message": "Post-onboarding tasks created successfully",
  "data": [
    {
      "id": 1,
      "employee_id": "EMP567073",
      "task_type": "send_welcome_email",
      "title": "Send Welcome Email",
      "status": "pending",
      "priority": "medium",
      "assigned_to": "HR",
      "due_date": "2024-02-20T00:00:00Z",
      "created_at": "2024-01-20T10:00:00Z"
    },
    {
      "id": 2,
      "employee_id": "EMP567073",
      "task_type": "create_email_account",
      "title": "Create Email Account",
      "status": "pending",
      "priority": "medium",
      "assigned_to": "HR",
      "due_date": "2024-02-20T00:00:00Z",
      "created_at": "2024-01-20T10:00:00Z"
    }
    // ... more tasks
  ]
}
```

**Example:**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/bulk
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "employee_id": "EMP567073",
  "task_types": [
    "send_welcome_email",
    "create_email_account",
    "allocate_laptop",
    "desk_allocate",
    "issue_id_card"
  ],
  "priority": "high",
  "due_date": "2024-02-20"
}
```

---

### 6. Get Tasks by Employee ID

Get all tasks for a specific employee.

**Endpoint:** `GET /api/v1/employees/post-onboarding/tasks/employee/:employee_id`

**Authentication:** Required

**URL Parameters:**
- `employee_id` (required) - Employee ID (e.g., EMP567073)

**Response:**
```json
{
  "success": true,
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": 1,
      "employee_id": "EMP567073",
      "task_type": "send_welcome_email",
      "title": "Send Welcome Email",
      "description": "Send welcome email to new employee",
      "status": "completed",
      "priority": "high",
      "assigned_to": "HR",
      "assigned_to_user_id": 5,
      "assigned_to_user_name": "Jane Smith",
      "completed_by_user_id": 5,
      "completed_by_user_name": "Jane Smith",
      "completed_at": "2024-01-21T14:30:00Z",
      "due_date": "2024-02-15T00:00:00Z",
      "notes": "Email sent successfully",
      "metadata": null,
      "created_at": "2024-01-20T10:00:00Z",
      "updated_at": "2024-01-21T14:30:00Z"
    },
    {
      "id": 2,
      "employee_id": "EMP567073",
      "task_type": "create_email_account",
      "title": "Create Email Account",
      "status": "pending",
      "priority": "critical",
      "assigned_to": "IT",
      "assigned_to_user_id": 3,
      "assigned_to_user_name": "John Admin",
      "completed_by_user_id": null,
      "completed_by_user_name": null,
      "completed_at": null,
      "due_date": "2024-02-14T00:00:00Z",
      "notes": "Create email account with format: firstname.lastname@company.com",
      "created_at": "2024-01-20T10:00:00Z",
      "updated_at": "2024-01-20T10:00:00Z"
    }
  ]
}
```

**Example:**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/employee/EMP567073
Authorization: Bearer {{TOKEN}}
```

---

### 7. Get Task Summary for Employee

Get a summary of tasks (completed, pending, in progress, etc.) for an employee.

**Endpoint:** `GET /api/v1/employees/post-onboarding/tasks/employee/:employee_id/summary`

**Authentication:** Required

**URL Parameters:**
- `employee_id` (required) - Employee ID (e.g., EMP567073)

**Response:**
```json
{
  "success": true,
  "message": "Task summary retrieved successfully",
  "data": {
    "employee_id": "EMP567073",
    "total_tasks": 14,
    "completed_tasks": 5,
    "pending_tasks": 7,
    "in_progress_tasks": 2,
    "skipped_tasks": 0,
    "completion_percentage": 35.71,
    "tasks_by_status": {
      "completed": 5,
      "pending": 7,
      "in_progress": 2,
      "skipped": 0
    },
    "tasks_by_priority": {
      "critical": 2,
      "high": 4,
      "medium": 6,
      "low": 2
    },
    "overdue_tasks": 1,
    "upcoming_due_tasks": 3
  }
}
```

**Example:**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/employee/EMP567073/summary
Authorization: Bearer {{TOKEN}}
```

---

### 8. List All Tasks (with Filters and Pagination)

Get a paginated list of all tasks with optional filters.

**Endpoint:** `GET /api/v1/employees/post-onboarding/tasks`

**Authentication:** Required

**Query Parameters:**
- `employee_id` (optional) - Filter by employee ID
- `status` (optional) - Filter by status (pending, in_progress, completed, skipped)
- `task_type` (optional) - Filter by task type
- `priority` (optional) - Filter by priority (low, medium, high, critical)
- `assigned_to` (optional) - Filter by assigned department/role
- `page` (optional) - Page number (default: 1)
- `page_size` (optional) - Items per page (default: 20, max: 100)

**Example Request:**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks?status=completed&priority=high&page=1&page_size=20
Authorization: Bearer {{TOKEN}}
```

**Response:**
```json
{
  "success": true,
  "message": "Tasks retrieved successfully",
  "data": [
    {
      "id": 1,
      "employee_id": "EMP567073",
      "task_type": "send_welcome_email",
      "title": "Send Welcome Email",
      "status": "completed",
      "priority": "high",
      "assigned_to": "HR",
      "completed_at": "2024-01-21T14:30:00Z",
      "created_at": "2024-01-20T10:00:00Z"
    }
    // ... more tasks
  ],
  "meta": {
    "page": 1,
    "per_page": 20,
    "total": 45,
    "total_pages": 3
  }
}
```

**Example: Get All Pending Tasks**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks?status=pending&page=1&page_size=50
Authorization: Bearer {{TOKEN}}
```

**Example: Get All Tasks for Specific Employee**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks?employee_id=EMP567073
Authorization: Bearer {{TOKEN}}
```

---

### 9. Get Single Task by ID

Get details of a specific task by its ID.

**Endpoint:** `POST /api/v1/employees/post-onboarding/tasks/get`

**Authentication:** Required

**Request Body:**
```json
{
  "id": 1
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task retrieved successfully",
  "data": {
    "id": 1,
    "employee_id": "EMP567073",
    "task_type": "send_welcome_email",
    "title": "Send Welcome Email",
    "description": "Send welcome email to new employee",
    "status": "completed",
    "priority": "high",
    "assigned_to": "HR",
    "assigned_to_user_id": 5,
    "assigned_to_user_name": "Jane Smith",
    "completed_by_user_id": 5,
    "completed_by_user_name": "Jane Smith",
    "completed_at": "2024-01-21T14:30:00Z",
    "due_date": "2024-02-15T00:00:00Z",
    "notes": "Email sent successfully",
    "metadata": null,
    "created_at": "2024-01-20T10:00:00Z",
    "updated_at": "2024-01-21T14:30:00Z"
  }
}
```

**Example:**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/get
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 1
}
```

---

### 10. Update Task

Update an existing task (status, priority, assignment, etc.).

**Endpoint:** `POST /api/v1/employees/post-onboarding/tasks/update`

**Authentication:** Required

**Request Body:**
```json
{
  "id": 2,
  "title": "Create Email Account - Updated",  // Optional
  "description": "Updated description",  // Optional
  "status": "in_progress",  // Optional: pending, in_progress, completed, skipped
  "priority": "critical",  // Optional: low, medium, high, critical
  "assigned_to": "IT",  // Optional
  "assigned_to_user_id": 3,  // Optional
  "due_date": "2024-02-16",  // Optional: YYYY-MM-DD format
  "notes": "Updated notes",  // Optional
  "metadata": "{\"updated\": true}"  // Optional: JSON string
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task updated successfully",
  "data": {
    "id": 2,
    "employee_id": "EMP567073",
    "task_type": "create_email_account",
    "title": "Create Email Account - Updated",
    "status": "in_progress",
    "priority": "critical",
    "assigned_to": "IT",
    "assigned_to_user_id": 3,
    "due_date": "2024-02-16T00:00:00Z",
    "notes": "Updated notes",
    "updated_at": "2024-01-22T09:00:00Z"
  }
}
```

**Example: Mark Task as In Progress**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/update
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 2,
  "status": "in_progress",
  "notes": "Started working on email account creation"
}
```

**Example: Change Task Priority**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/update
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 2,
  "priority": "critical"
}
```

---

### 11. Complete Task

Mark a task as completed. This automatically sets the status to "completed", records who completed it, and sets the completion timestamp.

**Endpoint:** `POST /api/v1/employees/post-onboarding/tasks/complete`

**Authentication:** Required

**Request Body:**
```json
{
  "id": 2,
  "notes": "Email account created successfully. Username: john.doe@company.com"  // Optional: completion notes
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task completed successfully",
  "data": {
    "id": 2,
    "employee_id": "EMP567073",
    "task_type": "create_email_account",
    "title": "Create Email Account",
    "status": "completed",
    "priority": "critical",
    "assigned_to": "IT",
    "assigned_to_user_id": 3,
    "assigned_to_user_name": "John Admin",
    "completed_by_user_id": 3,
    "completed_by_user_name": "John Admin",
    "completed_at": "2024-01-22T10:30:00Z",
    "notes": "Email account created successfully. Username: john.doe@company.com",
    "updated_at": "2024-01-22T10:30:00Z"
  }
}
```

**Example: Complete "Send Welcome Email" Task**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/complete
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 1,
  "notes": "Welcome email sent to john.doe@company.com at 2:30 PM"
}
```

**Example: Complete "Create Email Account" Task**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/complete
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 2,
  "notes": "Email account created: john.doe@company.com. Password sent via secure channel."
}
```

---

### 12. Delete Task

Delete a task (soft delete).

**Endpoint:** `POST /api/v1/employees/post-onboarding/tasks/delete`

**Authentication:** Required

**Request Body:**
```json
{
  "id": 3
}
```

**Response:**
```json
{
  "success": true,
  "message": "Task deleted successfully",
  "data": null
}
```

**Example:**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/delete
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 3
}
```

---

## Common Use Cases

### Use Case 1: Check if Welcome Email was Sent

**Step 1:** Get all tasks for the employee
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/employee/EMP567073
Authorization: Bearer {{TOKEN}}
```

**Step 2:** Look for task with `task_type: "send_welcome_email"` and check `status: "completed"`

**Alternative:** Filter by task type
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks?employee_id=EMP567073&task_type=send_welcome_email
Authorization: Bearer {{TOKEN}}
```

### Use Case 2: Mark Email Account as Created

**Step 1:** Find the task (if it exists)
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/employee/EMP567073
Authorization: Bearer {{TOKEN}}
```

**Step 2:** Complete the task
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/complete
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 2,
  "notes": "Email account created: john.doe@company.com"
}
```

### Use Case 3: Create All Initial Tasks for New Employee

After an employee completes onboarding, create all standard tasks:

```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/bulk
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "employee_id": "EMP567073",
  "task_types": [
    "send_welcome_email",
    "create_email_account",
    "allocate_laptop",
    "desk_allocate",
    "issue_id_card",
    "team_introduction",
    "policy_acknowledge",
    "induction_training",
    "assign_buddy_mentor",
    "first_week_review"
  ],
  "priority": "high",
  "due_date": "2024-02-20"
}
```

### Use Case 4: Get Task Completion Status

Check how many tasks are completed for an employee:

```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/employee/EMP567073/summary
Authorization: Bearer {{TOKEN}}
```

This returns:
- Total tasks
- Completed tasks count
- Pending tasks count
- Completion percentage
- Tasks by status and priority

### Use Case 5: Track Task Progress

**Get all pending tasks:**
```bash
GET {{BASE_URL}}/api/v1/employees/post-onboarding/tasks?status=pending&page=1&page_size=50
Authorization: Bearer {{TOKEN}}
```

**Update task status to in_progress:**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/update
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 2,
  "status": "in_progress"
}
```

**Complete the task:**
```bash
POST {{BASE_URL}}/api/v1/employees/post-onboarding/tasks/complete
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "id": 2,
  "notes": "Task completed successfully"
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": "employee_id is required"
}
```

### 401 Unauthorized
```json
{
  "success": false,
  "message": "Unauthorized",
  "data": null
}
```

### 404 Not Found
```json
{
  "success": false,
  "message": "Task not found",
  "data": null
}
```

### 500 Internal Server Error
```json
{
  "success": false,
  "message": "Internal server error",
  "data": null
}
```

---

## Notes

1. **Task Types:** All task types are predefined. Use the `GET /types` endpoint to see all available types.

2. **Employee ID:** The `employee_id` must be a valid employee ID (e.g., EMP567073) that exists in the system.

3. **Status Updates:** When you complete a task using the `/complete` endpoint, the system automatically:
   - Sets status to "completed"
   - Records the user who completed it (from JWT token)
   - Sets the completion timestamp
   - Updates the task

4. **Metadata Field:** The `metadata` field accepts a JSON string, allowing you to store task-specific data (e.g., email addresses, account details, etc.).

5. **Due Dates:** Due dates should be in `YYYY-MM-DD` format. The system will convert them to full timestamps.

6. **Pagination:** When listing tasks, use `page` and `page_size` query parameters. Default page size is 20, maximum is 100.

7. **Filtering:** You can combine multiple filters when listing tasks (e.g., `?employee_id=EMP567073&status=completed&priority=high`).

---

## Integration Example

Here's a complete flow for tracking post-onboarding tasks:

```javascript
// 1. After employee onboarding is completed, create initial tasks
const createTasks = async (employeeId) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/post-onboarding/tasks/bulk`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      employee_id: employeeId,
      task_types: [
        'send_welcome_email',
        'create_email_account',
        'allocate_laptop',
        'issue_id_card'
      ],
      priority: 'high',
      due_date: '2024-02-20'
    })
  });
  return response.json();
};

// 2. Check task status
const getTaskStatus = async (employeeId) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/post-onboarding/tasks/employee/${employeeId}/summary`, {
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  return response.json();
};

// 3. Mark task as completed
const completeTask = async (taskId, notes) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/post-onboarding/tasks/complete`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      id: taskId,
      notes: notes
    })
  });
  return response.json();
};
```

---

## Summary

The Post-Onboarding Tasks API provides comprehensive tracking for all tasks that need to be completed after an employee is onboarded. You can:

- ✅ Create single or bulk tasks
- ✅ Track task status (pending, in_progress, completed, skipped)
- ✅ Assign tasks to departments or specific users
- ✅ Set priorities and due dates
- ✅ Get task summaries and completion statistics
- ✅ Filter and paginate tasks
- ✅ Update and complete tasks
- ✅ Track who completed each task and when

This system ensures that all post-onboarding activities (like sending welcome emails, creating accounts, allocating equipment, etc.) are properly tracked and completed.
