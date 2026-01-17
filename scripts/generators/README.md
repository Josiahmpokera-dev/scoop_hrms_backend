# GORM Gen Model Generators

This directory contains scripts to generate GORM models from your database tables using GORM Gen.

## Usage

### Option 1: Generic Generator (Recommended)

Generate models for any module and table:

```bash
go run scripts/generators/generate_models.go <module_name> <table_name>
```

**Examples:**
```bash
# Generate users table models
go run scripts/generators/generate_models.go users users

# Generate employees table models
go run scripts/generators/generate_models.go employees employees

# Generate departments table models
go run scripts/generators/generate_models.go departments departments
```

### Option 2: Module-Specific Generators

Use pre-configured generators for specific modules:

```bash
# Generate users models
go run scripts/generators/generate_users.go
```

## Output Structure

After running a generator, you'll get:

```
internal/modules/<module_name>/repositories/query/
├── gen.go              # Generated query methods
├── gen_test.go         # Generated tests
└── model.gen.go        # Generated model structs
```

## Creating Custom Generators

To create a generator for a new module, copy `generate_users.go` and modify:

1. Change the `OutPath` to your module path
2. Change the table name in `GenerateModel()`

**Example for employees module:**

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
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:       "internal/modules/employees/repositories/query",
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.UseDB(db)

	g.ApplyBasic(
		g.GenerateModel("employees"),
	)

	g.Execute()

	fmt.Println("✅ Successfully generated models for 'employees' table")
}
```

## Generated Code Usage

After generation, you can use the generated models in your repositories:

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

func (r *UserRepository) FindByID(id uint) (*query.User, error) {
	return r.query.User.Where(r.query.User.ID.Eq(id)).First()
}
```

## Regenerating Models

When your database schema changes:

1. Update your database schema
2. Run the generator again:
   ```bash
   go run scripts/generators/generate_models.go users users
   ```
3. The generated code will be updated automatically

## Notes

- Generated files should be committed to version control
- Regenerate after schema changes
- Customize generator config as needed per module
- The generator reads from your `.env` database configuration
