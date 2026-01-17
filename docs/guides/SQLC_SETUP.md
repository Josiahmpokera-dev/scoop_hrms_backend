# SQLC Setup Guide

This project uses **SQLC** alongside **GORM** for type-safe database queries. SQLC generates Go code from SQL queries, providing compile-time safety and better performance.

## Architecture

- **GORM**: Used for migrations and existing repositories (continues to work)
- **SQLC**: Used for type-safe, optimized queries (new queries can use this)

Both can coexist - you can gradually migrate to SQLC or use both as needed.

## Installation

### 1. Install SQLC

```bash
# macOS
brew install sqlc

# Linux
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Or download from: https://github.com/sqlc-dev/sqlc/releases
```

### 2. Verify Installation

```bash
sqlc version
```

## Project Structure

```
internal/database/sqlc/
├── schema/          # SQL schema definitions
│   └── users.sql   # Users table schema
├── queries/         # SQL queries
│   └── users.sql   # User-related queries
└── (generated)      # Generated Go code (after running sqlc generate)
    ├── db.go       # Database connection
    ├── models.go    # Generated models
    └── users.sql.go # Generated query functions
```

## Usage

### 1. Generate Code

After writing SQL queries, generate Go code:

```bash
sqlc generate
```

This will create type-safe Go code in `internal/database/sqlc/`.

### 2. Using SQLC in Your Code

```go
import (
    "context"
    "github.com/Josiahmpokera-dev/hrms-backend/internal/database"
    db "github.com/Josiahmpokera-dev/hrms-backend/internal/database/sqlc"
)

// Get SQLC database connection
sqlcDB := database.GetSQLCDB()

// Create a new queries instance
queries := db.New(sqlcDB)

// Use generated queries
ctx := context.Background()

// Get user by email
user, err := queries.GetUserByEmail(ctx, "user@example.com")
if err != nil {
    // handle error
}

// Create a new user
newUser, err := queries.CreateUser(ctx, db.CreateUserParams{
    Username:     "johndoe",
    Email:        "john@example.com",
    PasswordHash: "hashed_password",
    FirstName:    "John",
    LastName:     "Doe",
    UserType:     "user",
    IsActive:     true,
})
```

### 3. Example: SQLC Repository

```go
package repositories

import (
    "context"
    "github.com/Josiahmpokera-dev/hrms-backend/internal/database"
    db "github.com/Josiahmpokera-dev/hrms-backend/internal/database/sqlc"
)

type UserSQLCRepository struct {
    queries *db.Queries
}

func NewUserSQLCRepository() *UserSQLCRepository {
    return &UserSQLCRepository{
        queries: db.New(database.GetSQLCDB()),
    }
}

func (r *UserSQLCRepository) GetByEmail(ctx context.Context, email string) (*db.User, error) {
    return r.queries.GetUserByEmail(ctx, email)
}

func (r *UserSQLCRepository) Create(ctx context.Context, params db.CreateUserParams) (*db.User, error) {
    return r.queries.CreateUser(ctx, params)
}
```

## Adding New Queries

### 1. Add SQL Query

Edit `internal/database/sqlc/queries/users.sql`:

```sql
-- name: GetActiveUsers :many
SELECT * FROM users
WHERE is_active = true AND deleted_at IS NULL
ORDER BY created_at DESC;
```

### 2. Regenerate Code

```bash
sqlc generate
```

### 3. Use in Code

```go
users, err := queries.GetActiveUsers(ctx)
```

## SQLC vs GORM

### Use SQLC when:
- You need maximum performance
- You want compile-time query safety
- You prefer writing raw SQL
- You want explicit control over queries

### Use GORM when:
- You need migrations (AutoMigrate)
- You prefer ORM-style queries
- You want rapid prototyping
- You need complex relationships

## Configuration

The SQLC configuration is in `sqlc.yaml`:

```yaml
version: "2"
sql:
  - engine: "postgresql"
    queries: "internal/database/sqlc/queries"
    schema: "internal/database/sqlc/schema"
    gen:
      go:
        package: "db"
        out: "internal/database/sqlc"
        sql_package: "pgx/v5"
```

## Generated Code

After running `sqlc generate`, you'll get:

- **Models**: Type-safe structs matching your database schema
- **Queries**: Functions for each SQL query
- **Interfaces**: Query interfaces for testing

## Best Practices

1. **Keep schemas in sync**: Update `schema/` files when database changes
2. **Version control generated code**: Commit generated files to git
3. **Use transactions**: SQLC supports transactions via `db.Queries.WithTx()`
4. **Test queries**: Write tests for your SQL queries
5. **Document queries**: Add comments in SQL files explaining complex queries

## Example: Transaction

```go
tx, err := sqlcDB.Begin(ctx)
if err != nil {
    return err
}
defer tx.Rollback(ctx)

qtx := queries.WithTx(tx)

// Use qtx for queries in transaction
user, err := qtx.CreateUser(ctx, params)
if err != nil {
    return err
}

err = tx.Commit(ctx)
return err
```

## Troubleshooting

### SQLC not generating code
- Check `sqlc.yaml` configuration
- Verify SQL syntax is correct
- Run `sqlc generate -v` for verbose output

### Type mismatches
- Ensure schema matches actual database
- Regenerate after schema changes: `sqlc generate`

### Connection issues
- Verify database credentials in `.env`
- Check `internal/database/sqlc_connection.go`

## Next Steps

1. Run `sqlc generate` to generate initial code
2. Create a sample repository using SQLC
3. Gradually migrate queries from GORM to SQLC
4. Add more queries as needed

## Resources

- [SQLC Documentation](https://docs.sqlc.dev/)
- [SQLC GitHub](https://github.com/sqlc-dev/sqlc)
- [PostgreSQL SQLC Examples](https://docs.sqlc.dev/en/latest/tutorials/getting-started-postgresql.html)
