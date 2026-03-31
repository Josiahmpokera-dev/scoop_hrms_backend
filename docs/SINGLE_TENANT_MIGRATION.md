# Single-Tenant Migration

## Overview

The HRMS backend has been migrated from a multi-tenant architecture to a **single-tenant** architecture. This document describes what changed, what stayed the same, and how it affects existing deployments.

---

## What Changed

### 1. Models — `TenantID` Field Removed

The `TenantID *uint` field was removed from **all** model structs across the application:

| Module | Models Affected |
|--------|----------------|
| Users | `User` |
| Employees | `Employee`, `EmployeeOnboardingDraft` |
| Organizations | `Organization`, `Department`, `Team`, `JobPosition`, `Location`, `OrganizationUnit`, `CostCenter` |
| Payroll | `PayrollRun`, `SalaryStructure`, `SalaryComponent`, `Payslip`, `Loan`, `TaxSlab`, `StatutoryRule`, `NHIFSchedule`, `CompliancePayment` |
| Leave | `LeaveType`, `LeavePolicy`, `LeaveRequest`, `LeaveApproval`, `LeaveDocument`, `LeaveBalance`, `Holiday` |
| Helpdesk | `Ticket`, `Comment`, `Attachment`, `TicketCategory`, `RoutingRule`, `KnowledgeBaseArticle`, `KBArticleFeedback` |
| Assets | `Asset`, `AssetRequest`, `AssetIssue` |
| Shifts | `Shift`, `RosterAssignment`, `SwapRequest`, `RosterChangeRequest` |
| Biometric | `BioTimeTransaction`, `BioTimeConfig` | 
| Dashboard | `Announcement` |
| Audit | `AuditLog` |

**Database impact**: The `tenant_id` column remains in the database tables (GORM's `AutoMigrate` does not drop columns). Existing rows will retain their values. New rows will have `tenant_id = NULL`. This is fully backward-compatible.

### 2. Middleware

| File | Change | 
|------|--------|
| `internal/middleware/tenant.go` | Stripped down to a stub. `GetTenantID()` always returns `nil`. `TenantMiddleware()` and `RequireTenantMiddleware()` removed. |
| `internal/middleware/auth.go` | Removed code that set `tenant_id` in context from the user record. | 
| `internal/middleware/audit.go` | Removed code that read `tenant_id` from context for audit logging. | 

### 3. Auth / Onboarding

| File | Change |
|------|--------|
| `internal/modules/auth/services/onboarding_service.go` | Removed Tenant creation during onboarding. The flow now creates User + Organization only. |
| `internal/modules/auth/models/onboarding.go` | Removed `TenantInfo` struct and `Tenant` field from `OnboardingResponse`. |

**Before** (onboarding response):
```json
{
  "user": { ... },
  "organization": { ... },
  "tenant": { "id": 1, "name": "...", "domain": "..." },
  "access_token": "...",
  "refresh_token": "..."
}
```

**After** (onboarding response):
```json
{
  "user": { ... },
  "organization": { ... },
  "access_token": "...",
  "refresh_token": "..."
}
```

### 4. Repositories — Tenant Filtering Removed

All repository methods that previously filtered by `tenant_id` no longer apply that filter. Queries now return all records regardless of the (now unused) `tenant_id` column.

Affected repositories across all modules: Employees, Users, Dashboard, Payroll, Leave, Helpdesk, Shifts, Assets, Biometric, Departments, Locations, Positions, Teams, Organization Units, Cost Centers, Organizations, Roles, Audit.

### 5. Services — Tenant Ownership Checks Removed

All services that previously validated tenant ownership (e.g., checking if an asset belongs to the same tenant as the requesting user) no longer perform those checks. In single-tenant mode, all data belongs to one organization.

### 6. Handlers

All handlers that previously extracted `tenantID` from the request context via `middleware.GetTenantID(c)` still call that function, but it always returns `nil`. This means:
- No tenant filtering is applied to queries
- No tenant ID is set on newly created records
- The `X-Tenant-ID` header is no longer read or required

### 7. CORS Configuration

- Removed `X-Tenant-ID` from the default allowed CORS headers.
- The `X-Tenant-ID` header is no longer needed in API requests.

### 8. Database Migrations

- The `Tenant` model is no longer auto-migrated. The `tenants` table remains in the database but is no longer used by the application.

---

## What Did NOT Change

| Aspect | Status |
|--------|--------|
| **API endpoints** | All URL paths remain identical. No routes added or removed. |
| **Request body format** | All request payloads remain the same. `tenant_id` was never a required request field. |
| **Response body format** | Unchanged, except the onboarding response no longer includes `tenant`. The `tenant_id` field no longer appears in JSON responses (it was previously `null` for most users). |
| **Authentication** | JWT-based auth is unchanged. Tokens, login, refresh all work identically. |
| **Authorization** | Role-based access (Admin, HR, IT, User) is unchanged. |
| **Database schema** | No columns dropped. The `tenant_id` columns and `tenants` table remain but are unused. |
| **Function signatures** | Handler and service method signatures that accept `tenantID *uint` parameters are unchanged — they simply receive `nil`. |

---

## Migration Steps for Deployed Servers

1. **Pull the latest code** and rebuild:
   ```bash
   git pull
   docker compose down
   docker compose build --no-cache backend
   docker compose up -d
   ```

2. **No database migration needed**. GORM's `AutoMigrate` will not drop the existing `tenant_id` columns. They will simply be unused.

3. **No frontend changes required** for the `X-Tenant-ID` header — it was never enforced. If your frontend sends it, it will be silently ignored.

4. **(Optional) Clean up the database** if you want to remove the unused columns and table:
   ```sql
   -- Optional: Drop the tenants table (no longer used)
   DROP TABLE IF EXISTS tenants;
   
   -- Optional: Remove tenant_id columns from individual tables
   -- Only do this if you're certain you won't need multi-tenancy in the future
   -- ALTER TABLE users DROP COLUMN IF EXISTS tenant_id;
   -- ALTER TABLE employees DROP COLUMN IF EXISTS tenant_id;
   -- ... (repeat for all tables)
   ```

---

## Files Changed (Summary)

| Category | Files Modified |
|----------|---------------|
| Middleware | `tenant.go`, `auth.go`, `audit.go` |
| Config | `config.go` |
| Main | `main.go` |
| Auth | `onboarding_service.go`, `onboarding.go` |
| Models | 35+ model files across all modules |
| Repositories | 30+ repository files |
| Services | 20+ service files |
| Handlers | 15+ handler files |

**Total**: ~100+ files modified, zero API endpoint changes, zero breaking changes for existing clients.
