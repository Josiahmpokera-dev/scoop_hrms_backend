# GORM Gen Model Generation Guide

This project uses **GORM Gen** to generate type-safe models and query methods from your database schema.

## Quick Start

### 1. Generate Models for a Module

```bash
# Generic generator (recommended)
make generate-models MODULE=users TABLE=users

# Or directly
go run scripts/generators/generate_models.go users users
```

### 2. Generated Code Location

After generation, you'll find:
```
internal/modules/users/repositories/query/
├── gen.go          # Generated query methods
├── gen_test.go     # Generated tests
└── model.gen.go     # Generated model structs
```

## Usage Examples

### Example 1: Generate Users Models

```bash
make generate-models MODULE=users TABLE=users
```

### Example 2: Generate Employees Models

```bash
make generate-models MODULE=employees TABLE=employees
```

### Example 3: Generate Departments Models

```bash
make generate-models MODULE=departments TABLE=departments
```

## Using Generated Code

### In Your Repository

```go
package repositories

import (
    "github.com/Josiahmpokera-dev/hrms-backend/internal/database"
    "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories/query"
)

type UserRepository struct {
    query *query.Query
}

func NewUserRepository() *UserRepository {
    return &UserRepository{
        query: query.Use(database.GetDB()),
    }
}

// Find by email using generated code
func (r *UserRepository) FindByEmail(email string) (*query.User, error) {
    return r.query.User.Where(r.query.User.Email.Eq(email)).First()
}

// Find by ID
func (r *UserRepository) FindByID(id uint) (*query.User, error) {
    return r.query.User.Where(r.query.User.ID.Eq(id)).First()
}

// Create user
func (r *UserRepository) Create(user *query.User) error {
    return r.query.User.Create(user)
}
```

## Creating Custom Generators

To create a generator for a specific module, copy `scripts/generators/generate_users.go`:

1. Change `OutPath` to your module path
2. Change table name in `GenerateModel()`

**Template:**

```go
package main

import (
    "fmt"
    "github.com/Josiahmpokera-dev/hrms-backend/internal/config"
    "gorm.io/driver/postgres"
    "gorm.io/gen"
    "gorm.io/gorm"
)

func main() {
    cfg, err := config.LoadConfig()
    if err != nil {
        panic(fmt.Sprintf("Failed to load config: %v", err))
    }

    dsn := cfg.Database.GetDSN()
    db, err := gorm.Open(postgres.Open(dsn))
    if err != nil {
        panic(fmt.Sprintf("Failed to connect: %v", err))
    }

    g := gen.NewGenerator(gen.Config{
        OutPath:       "internal/modules/<module>/repositories/query",
        Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
        FieldNullable: true,
    })

    g.UseDB(db)
    g.ApplyBasic(g.GenerateModel("<table_name>"))
    g.Execute()

    fmt.Println("✅ Models generated successfully")
}
```

## Workflow

1. **Create/Update Database Schema**
   - Use GORM migrations or direct SQL
   - Ensure table exists in database

2. **Generate Models**
   ```bash
   make generate-models MODULE=<module> TABLE=<table>
   ```

3. **Use Generated Code**
   - Import generated query package
   - Use in your repositories

4. **Regenerate After Schema Changes**
   - Always regenerate when table structure changes
   - Generated code will be updated

## Generated Query Methods

GORM Gen generates type-safe query methods:

```go
// Find operations
queries.User.Where(queries.User.Email.Eq("test@example.com")).First()
queries.User.Where(queries.User.ID.In(1, 2, 3)).Find()

// Create operations
queries.User.Create(&user)

// Update operations
queries.User.Where(queries.User.ID.Eq(id)).Update(queries.User.Email, newEmail)

// Delete operations
queries.User.Where(queries.User.ID.Eq(id)).Delete()
```

## Best Practices

1. **Regenerate after schema changes**: Always run generator after modifying tables
2. **Commit generated code**: Include generated files in version control
3. **Use in repositories**: Keep generated code usage in repository layer
4. **Combine with GORM**: You can use both GORM and generated code
5. **Test generated code**: Write tests for your repositories

## Troubleshooting

### "Table does not exist" error
- Ensure table exists in database
- Check table name spelling
- Verify database connection

### Generated code not updating
- Delete old generated files
- Regenerate: `make generate-models MODULE=<name> TABLE=<name>`

### Import errors
- Ensure you've run the generator
- Check module path in generator script
- Verify Go module name matches

## Available Generators

- `scripts/generators/generate_models.go` - Generic generator (recommended)
- `scripts/generators/generate_users.go` - Users-specific generator

## Makefile Commands

```bash
# Generate models for any module
make generate-models MODULE=users TABLE=users

# Generate users models
make generate-users

# Build application
make build

# Run application
make run
```
