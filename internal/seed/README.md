# Seed package

This package holds all database seeding logic. It runs only when the app is in **development** (`SERVER_ENV=development`).

## How to call it from main.go

In `cmd/api/main.go`, after migrations and before starting the server, call:

```go
// Seed initial data (only in development)
seed.Run()
```

That single call runs all seeders in order: permissions → roles → admin user → non-employee users.

## Files

| File | Purpose |
|------|---------|
| `seed.go` | Entry point: `Run()` checks env and calls all seeders. |
| `permissions.go` | Seeds default permissions (RBAC). |
| `roles.go` | Seeds default roles (admin, hr, employee) and assigns permissions. |
| `users.go` | Seeds admin user and non-employee users (IT, HR samples). |

## Adding new seeders

1. Add a new file (e.g. `departments.go`) with a `RunXxx()` function.
2. Call it from `Run()` in `seed.go`:

```go
func Run() {
	// ...
	RunPermissions()
	RunRoles()
	RunAdminUser()
	RunNonEmployeeUsers()
	RunDepartments() // new
	// ...
}
```

## Credentials

- **Admin:** See logs on first run or `users.go` (admin@hrms.com / admin123).
- **Non-employee users (IT, HR):** See `docs/SEED_NON_EMPLOYEE_USERS.md`.
