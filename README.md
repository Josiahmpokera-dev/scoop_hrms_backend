# HRMS Backend

Human Resource Management System (HRMS) Backend API built with Go and Gin framework.

## Project Structure

```
hrms-backend/
├── cmd/
│   └── api/                    # Application entry point
│       └── main.go
│
├── internal/                   # Private application code
│   ├── app/                   # Application initialization
│   ├── config/                # Configuration management
│   ├── database/              # Database connection & setup
│   │   └── optional/         # Optional tools (SQLC - not actively used)
│   ├── middleware/            # HTTP middleware (auth, RBAC, tenant)
│   ├── router/                # HTTP routing
│   ├── types/                 # Type definitions (API request/response)
│   ├── utils/                 # Utility functions
│   │   ├── handlers/         # Generic action handlers
│   │   └── response/         # HTTP response helpers
│   └── modules/               # Business logic modules
│       ├── auth/             # ✅ Authentication & Authorization
│       ├── users/            # ✅ User Management
│       ├── tenants/          # ✅ Multi-tenancy
│       ├── roles/            # ✅ RBAC (Roles & Permissions)
│       ├── organizations/    # ✅ Organization Management
│       ├── organization_units/ # ✅ Business Units/Divisions
│       ├── departments/      # ✅ Department Management
│       ├── teams/            # ✅ Team Management
│       ├── positions/        # ✅ Job Positions
│       ├── locations/        # ✅ Location Management
│       ├── cost_centers/      # ✅ Cost Center Management
│       ├── employees/        # ✅ Employee Management
│       ├── attendance/       # 🚧 Placeholder (future)
│       ├── leave/            # 🚧 Placeholder (future)
│       ├── payroll/          # 🚧 Placeholder (future)
│       ├── recruitment/      # 🚧 Placeholder (future)
│       └── performance/      # 🚧 Placeholder (future)
│
├── scripts/                   # Build and utility scripts
│   └── generators/           # GORM Gen model generators
│
├── docs/                      # Documentation
│   ├── guides/              # Setup guides (GORM Gen, SQLC)
│   └── api/                # API documentation
├── go.mod                     # Go module definition
├── Makefile                   # Build commands
└── README.md                  # This file
```

**Each module follows standard structure:**
```
internal/modules/<module>/
├── handlers/          # HTTP handlers (controllers)
├── models/            # Data models (GORM structs + request DTOs)
├── repositories/      # Data access layer
└── services/          # Business logic layer
```

**Note:** Some modules (attendance, leave, payroll, recruitment, performance) are placeholders for future implementation.

📖 **See [docs/api/API_DOCUMENTATION.md](docs/api/API_DOCUMENTATION.md) for API documentation.**

## Architecture Pattern

This project follows a **Layered Architecture** pattern:

1. **Handlers Layer** (`handlers/`): HTTP request/response handling
2. **Services Layer** (`services/`): Business logic and validation
3. **Repositories Layer** (`repositories/`): Data access and database operations
4. **Models Layer** (`models/`): Data structures and entities

## Module Structure

Each module follows the same structure:

```
module-name/
├── handlers/          # HTTP handlers (controllers)
├── services/          # Business logic
├── repositories/      # Data access
│   └── query/        # Generated GORM Gen code (auto-generated)
└── models/            # Data models (GORM structs)
```

### Model Generation

This project uses **GORM Gen** to generate type-safe models from your database schema.

**Generate models for a module:**
```bash
# Generic generator (recommended)
make generate-models MODULE=users TABLE=users
make generate-models MODULE=employees TABLE=employees

# Or use module-specific generators
make generate-users
```

Generated code will be in: `internal/modules/<module>/repositories/query/`

See `scripts/generators/README.md` for generator usage details.

## Key Dependencies

- **Gin**: HTTP web framework
- **GORM**: ORM for database operations and migrations
- **GORM Gen**: Code generation for type-safe database queries
- **PostgreSQL**: Database driver
- **JWT**: Authentication tokens
- **Validator**: Request validation
- **CORS**: Cross-origin resource sharing
- **Godotenv**: Environment variable management
- **Bcrypt**: Password hashing

## Getting Started

1. Install dependencies:
   ```bash
   go mod download
   ```

2. Set up environment variables:
   ```bash
   cp env.example .env
   # Edit .env with your database credentials
   ```

3. Generate GORM models (for each module):
   ```bash
   # Generate users models
   make generate-users
   # Or: go run scripts/generators/generate_users.go
   
   # Generate models for any module
   make generate-models MODULE=employees TABLE=employees
   # Or: go run scripts/generators/generate_models.go employees employees
   ```

4. Start the server:
   ```bash
   go run cmd/api/main.go
   # Or: make run
   ```

## Environment Variables

Create a `.env` file with the following variables:

```
# Server
PORT=8080
ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=hrms_db

# JWT
JWT_SECRET=your_jwt_secret_key
JWT_EXPIRY=24h
```

## Development Guidelines

1. Follow Go naming conventions
2. Keep handlers thin - delegate to services
3. Services contain business logic
4. Repositories handle all database operations
5. Use models for data structures
6. Generate models after schema changes: `make generate-models MODULE=<name> TABLE=<name>`
7. Add proper error handling
8. Write tests for services and repositories

## Model Generation Workflow

1. **Create/Update database schema** (via migrations or directly)
2. **Generate models** using GORM Gen:
   ```bash
   make generate-models MODULE=users TABLE=users
   ```
3. **Use generated code** in your repositories:
   ```go
   import "github.com/Josiahmpokera-dev/hrms-backend/internal/modules/users/repositories/query"
   
   queries := query.Use(database.GetDB())
   user, err := queries.User.Where(queries.User.Email.Eq(email)).First()
   ```
4. **Regenerate** whenever schema changes

## License

[Your License Here]
