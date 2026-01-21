# Employee User Account Creation on Onboarding

## Overview

When an employee completes onboarding, the system automatically creates a user account for them, allowing them to log in to the system. The employee receives their login credentials (email and generated password) in the onboarding completion response.

**Feature:** Automatic user account creation upon onboarding completion

**Default Role:** `user` (not admin)

**Multiple Roles Support:** Users can have multiple roles assigned (e.g., both "employee" and "system admin")

---

## How It Works

### 1. Onboarding Completion Process

When an employee completes onboarding:

1. **Employee Record Created:** The employee record is created in the `employees` table
2. **User Account Created:** A user account is automatically created in the `users` table
3. **Role Assignment:** The user is assigned:
   - Default role: `user` (from UserRole enum - not admin)
   - Additional role: `employee` (from roles table, if exists)
4. **Credentials Generated:** A secure random password is generated
5. **Employee-User Link:** The employee record is linked to the user account via `user_id`
6. **Credentials Returned:** Email and password are returned in the response

### 2. User Account Details

- **Email:** Uses the employee's work email (`official_email` from Step 2)
- **Username:** Generated from email (part before @, sanitized)
- **Password:** Secure random password (12 characters with complexity)
- **Role:** Default `user` role (not admin)
- **Status:** Active
- **Email Verified:** False (employee should verify on first login)

### 3. Role Assignment

The system assigns roles in this order:

1. **Default Role:** `user` (from UserRole enum) - Always assigned
2. **Additional Role:** `employee` (from roles table) - Assigned if role exists with code "employee"

**Note:** Users can have multiple roles. An admin can have both "employee" and "system admin" roles.

---

## API Endpoints

### Complete Onboarding (with User Creation)

**Endpoint:** `POST /api/v1/employees/onboarding/complete`

**Request Body:**
```json
{
  "draft_id": 123
}
```

**Response:**
```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "employee": {
      "id": 1,
      "employee_id": "EMP567840",
      "first_name": "John",
      "last_name": "Doe",
      "work_email": "john.doe@company.com",
      "user_id": 45,
      ...
    },
    "credentials": {
      "email": "john.doe@company.com",
      "username": "johndoe",
      "password": "aB3xY9mK2!A1"
    }
  }
}
```

**Example Request:**
```bash
POST {{BASE_URL}}/api/v1/employees/onboarding/complete
Authorization: Bearer {{TOKEN}}
Content-Type: application/json

{
  "draft_id": 123
}
```

---

### Complete Onboarding by Employee ID (with User Creation)

**Endpoint:** `POST /api/v1/employees/onboarding/{employee_id}/complete`

**URL Parameters:**
- `employee_id` (required) - Employee ID (e.g., "EMP567840")

**Response:**
```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "employee": {
      "id": 1,
      "employee_id": "EMP567840",
      "first_name": "John",
      "last_name": "Doe",
      "work_email": "john.doe@company.com",
      "user_id": 45,
      ...
    },
    "credentials": {
      "email": "john.doe@company.com",
      "username": "johndoe",
      "password": "aB3xY9mK2!A1"
    }
  }
}
```

**Example Request:**
```bash
POST {{BASE_URL}}/api/v1/employees/onboarding/EMP567840/complete
Authorization: Bearer {{TOKEN}}
```

---

## Response Fields

### Employee Object
- Standard employee fields (id, employee_id, first_name, last_name, etc.)
- **`user_id`:** ID of the created user account (linked)

### Credentials Object

| Field | Type | Description |
|-------|------|-------------|
| `email` | string | Employee's work email (login email) |
| `username` | string | Generated username (from email) |
| `password` | string | Generated secure password (employee should change on first login) |

---

## Password Generation

The system generates a secure random password with the following characteristics:

- **Length:** 12 characters
- **Complexity:** Includes uppercase, lowercase, numbers, and special characters
- **Format:** 8 random characters + 1 uppercase + 1 number + 1 special character
- **Security:** Generated using cryptographically secure random number generator

**Example Password:** `aB3xY9mK2!A1`

**Important:** The employee should change this password on their first login.

---

## Role Assignment Details

### Default Role: `user`

- Assigned to all new employee users
- Not an admin role
- Provides basic user access

### Additional Role: `employee`

- Assigned if a role with code "employee" exists in the `roles` table
- Allows users to have multiple roles
- Can be combined with other roles (e.g., "system admin")

### Multiple Roles Support

Users can have multiple roles assigned:

- **Via UserRole enum:** One role (user, hr, admin)
- **Via roles table:** Multiple roles via `user_roles` table

**Example:** An admin user can have:
- UserRole: `admin`
- Roles: `employee`, `system_admin`, `hr_manager`

---

## Error Handling

### User Already Exists

If a user with the employee's email already exists:

- **Behavior:** User creation is skipped
- **Response:** `credentials` field will be `null`
- **Employee:** Still created and linked to existing user (if `user_id` matches)

**Response Example:**
```json
{
  "success": true,
  "message": "Employee onboarding completed successfully",
  "data": {
    "employee": { ... },
    "credentials": null
  }
}
```

### Missing Work Email

If the employee doesn't have a work email:

- **Error:** "employee work email is required to create user account"
- **Behavior:** User creation fails, but onboarding continues
- **Response:** `credentials` field will be `null`

### Role Assignment Failure

If role assignment fails:

- **Behavior:** User is still created with default `user` role
- **Error:** Logged but doesn't fail the onboarding process

---

## Security Considerations

1. **Password Security:**
   - Passwords are hashed using bcrypt before storage
   - Plain password is only returned in the response (one-time)
   - Employee should change password on first login

2. **Email Verification:**
   - New users have `email_verified: false`
   - Employee should verify email on first login

3. **Default Role:**
   - New users are assigned `user` role (not admin)
   - Admin roles must be assigned manually

4. **Tenant Isolation:**
   - Users are created with the same `tenant_id` as the employee
   - Ensures proper multi-tenant isolation

---

## Use Cases

### 1. Standard Employee Onboarding

1. HR completes employee onboarding
2. System creates employee record
3. System creates user account automatically
4. HR receives credentials in response
5. HR shares credentials with employee
6. Employee logs in and changes password

### 2. Employee with Existing Account

1. Employee already has a user account
2. Onboarding completes
3. System detects existing user
4. Employee record is linked to existing user
5. No new credentials generated

### 3. Admin Employee

1. Employee completes onboarding
2. System creates user with `user` role
3. Admin manually assigns additional roles (e.g., "system admin")
4. Employee can have multiple roles

---

## JavaScript Examples

### Complete Onboarding and Get Credentials

```javascript
const completeOnboarding = async (draftId) => {
  const response = await fetch(`${BASE_URL}/api/v1/employees/onboarding/complete`, {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      draft_id: draftId
    })
  });
  
  const data = await response.json();
  
  if (data.success && data.data.credentials) {
    console.log('Employee Credentials:');
    console.log('Email:', data.data.credentials.email);
    console.log('Username:', data.data.credentials.username);
    console.log('Password:', data.data.credentials.password);
    
    // Share credentials with employee
    // Employee should change password on first login
  }
  
  return data;
};
```

### Handle Credentials Display

```javascript
const handleOnboardingComplete = (response) => {
  if (response.data.credentials) {
    // Show credentials to HR/admin
    alert(`Employee Login Credentials:\n\nEmail: ${response.data.credentials.email}\nPassword: ${response.data.credentials.password}\n\nPlease share these with the employee. They should change their password on first login.`);
  } else {
    // User already exists or email missing
    alert('Employee created successfully. User account could not be created (may already exist or email missing).');
  }
};
```

---

## Notes

1. **Password Change:** Employees should change their password on first login for security.

2. **Email Verification:** New users should verify their email address.

3. **Role Management:** Additional roles can be assigned via the roles API after user creation.

4. **Existing Users:** If a user with the email already exists, no new user is created, but the employee is still linked.

5. **Work Email Required:** User creation requires the employee to have a work email (from Step 2 of onboarding).

6. **Non-Blocking:** User creation failure doesn't block onboarding completion. Employee can be created even if user creation fails.

7. **Multiple Roles:** Users can have multiple roles. The default `user` role is always assigned, and additional roles from the roles table can be added.

---

## Summary

The automatic user creation feature:

- ✅ Creates user account automatically on onboarding completion
- ✅ Generates secure random password
- ✅ Assigns default `user` role (not admin)
- ✅ Supports multiple roles (employee + system admin, etc.)
- ✅ Links employee to user account
- ✅ Returns credentials in response
- ✅ Handles existing users gracefully
- ✅ Non-blocking (onboarding continues even if user creation fails)
- ✅ Secure password hashing
- ✅ Tenant isolation

This ensures that every onboarded employee can immediately access the system with their credentials.
