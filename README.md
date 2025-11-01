# Citary Backend

A modern, production-ready Go backend API built with Clean Architecture principles.

> 💡 **New to the project?** Check out the [Quick Start Guide](QUICK_START.md) for common commands and tasks!

## 🚀 Features

- ✅ **Clean Architecture** - Separation of concerns with clear boundaries
- ✅ **Type-Safe** - Strongly typed throughout with no `interface{}` abuse
- ✅ **Context-Aware** - Proper context propagation for timeouts and cancellation
- ✅ **Production Ready** - Middleware for CORS, logging, and panic recovery
- ✅ **Dependency Injection** - No global variables, testable design
- ✅ **Graceful Shutdown** - Proper cleanup of resources
- ✅ **Well Documented** - Comprehensive GoDoc comments

## 📋 Prerequisites

- Go 1.24.3 or higher
- PostgreSQL database

## 🛠️ Installation

1. **Clone the repository**

2. **Set up environment variables**
   ```bash
   cp .env.example .env
   ```

3. **Configure `.env`**
   ```env
   DATABASE_URL=postgres://user:password@localhost:5432/citary?sslmode=disable
   PORT=3005
   ```

4. **Install dependencies**
   ```bash
   go mod download
   ```

5. **Build the application**
   ```bash
   go build -o citary-backend .
   ```

6. **Run the application**
   ```bash
   ./citary-backend
   ```

## 📡 API Endpoints

### Authentication

#### **User Signup**
```http
POST /auth/signup
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

**Success Response (201 Created)**
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "emailVerified": false,
    "createdDate": "2024-10-31T22:00:00Z"
  }
}
```

**Error Response (400 Bad Request)**
```json
{
  "success": false,
  "message": "Email format is invalid",
  "error": {
    "code": 400,
    "message": "Email format is invalid"
  }
}
```

**Password Requirements:**
- Minimum 8 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit
- At least one special character

### Health Check

#### **Health Status**
```http
GET /health
```

**Response (200 OK)**
```json
{
  "status": "ok"
}
```

## 🏗️ Project Structure

```
.
├── internal/                     # Private application code
│   ├── domain/                   # Business logic layer (core)
│   │   ├── entities/             # Domain entities (User, etc.)
│   │   ├── repositories/         # Repository interfaces (ports)
│   │   ├── usecases/             # Business use cases
│   │   │   └── auth/             # Authentication use cases
│   │   ├── dtos/                 # Data transfer objects
│   │   │   └── auth/             # Auth DTOs (SignupRequest, etc.)
│   │   └── errors/               # Domain-specific errors
│   │
│   └── infrastructure/           # External concerns (adapters)
│       ├── persistence/          # Database layer
│       │   └── postgres/
│       │       ├── entities/     # Database models
│       │       ├── mappers/      # Domain ↔ DB mappers
│       │       └── repositories/ # Repository implementations
│       │
│       ├── http/                 # HTTP layer
│       │   ├── handlers/         # Request handlers
│       │   │   └── auth/         # Auth handlers
│       │   ├── middleware/       # HTTP middleware
│       │   │   ├── cors.go       # CORS handling
│       │   │   ├── logging.go    # Request logging
│       │   │   └── recovery.go   # Panic recovery
│       │   ├── dto/              # API request/response DTOs
│       │   ├── response/         # Response helpers
│       │   └── router/           # Route configuration
│       │
│       ├── config/               # Configuration management
│       └── di/                   # Dependency injection container
│
├── pkg/                          # Public shared packages
│   └── constants/                # Application constants
│
├── main.go                       # Application entry point
├── go.mod                        # Go module definition
├── .env.example                  # Example environment variables
├── README.md                     # This file
└── ARCHITECTURE.md               # Detailed architecture documentation
```

## 🔧 Development

### Project Setup

#### 1. First Time Setup
```bash
# Clone the repository
git clone <repository-url>
cd citary-backend

# Install Go dependencies
go mod download

# Verify dependencies
go mod verify

# Copy environment template
cp .env.example .env

# Edit .env with your configuration
nano .env  # or use your preferred editor
```

#### 2. Install New Dependencies
```bash
# Add a new dependency
go get github.com/package/name

# Add a specific version
go get github.com/package/name@v1.2.3

# Update go.mod and go.sum
go mod tidy

# Example: Add JWT library
go get github.com/golang-jwt/jwt/v5
```

#### 3. Update Dependencies
```bash
# Update all dependencies to latest minor/patch versions
go get -u ./...

# Update a specific package
go get -u github.com/package/name

# Clean up unused dependencies
go mod tidy
```

### Build Commands

#### Development Build
```bash
# Build for current platform
go build -o citary-backend .

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o citary-backend.exe .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o citary-backend .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o citary-backend .
```

#### Production Build (Optimized)
```bash
# Build with optimizations (smaller binary, no debug info)
go build -ldflags="-s -w" -o citary-backend .

# Build with version information
VERSION=$(git describe --tags --always --dirty)
go build -ldflags="-X main.Version=$VERSION" -o citary-backend .
```

### Run Commands

#### Run Directly
```bash
# Run without building
go run main.go

# Run with race detector (development)
go run -race main.go

# Run with specific environment
DATABASE_URL="postgres://..." PORT=3005 go run main.go
```

#### Run Built Binary
```bash
# Run the compiled binary
./citary-backend

# Windows
citary-backend.exe

# Run in background (Linux/macOS)
nohup ./citary-backend > server.log 2>&1 &

# Run with custom port
PORT=8080 ./citary-backend
```

### Code Quality

#### Format Code
```bash
# Format all Go files
go fmt ./...

# Format specific package
go fmt ./internal/domain/...

# Use gofmt directly for more control
gofmt -w .
```

#### Lint and Vet
```bash
# Run Go vet (built-in static analysis)
go vet ./...

# Run vet on specific package
go vet ./internal/infrastructure/...

# Install and run golangci-lint (recommended)
# Install: https://golangci-lint.run/usage/install/
golangci-lint run

# Run with auto-fix
golangci-lint run --fix
```

#### Static Analysis
```bash
# Check for common mistakes
go vet ./...

# Check for shadowed variables
go vet -shadow ./...

# Install staticcheck
go install honnef.co/go/tools/cmd/staticcheck@latest

# Run staticcheck
staticcheck ./...
```

### Testing

#### Run Tests
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific test
go test -v -run TestSignupUserUseCase ./internal/domain/usecases/auth/

# Run tests with race detector
go test -race ./...

# Run benchmarks
go test -bench=. ./...
```

#### Test Coverage
```bash
# Generate coverage for all packages
go test -coverprofile=coverage.out ./...

# View coverage in terminal
go tool cover -func=coverage.out

# View coverage in browser
go tool cover -html=coverage.out

# Coverage with specific threshold
go test -cover ./... | grep -E "coverage: [0-9]+\.[0-9]+%"
```

### Clean Project

#### Clean Build Artifacts
```bash
# Remove compiled binaries
rm -f citary-backend citary-backend.exe

# Clean Go build cache
go clean

# Clean all cached test results
go clean -testcache

# Clean module cache (careful - will re-download)
go clean -modcache

# Full cleanup
go clean -i -r -cache -testcache -modcache
```

#### Reset Dependencies
```bash
# Remove go.sum and re-download
rm go.sum
go mod download

# Verify all dependencies
go mod verify

# Tidy up (remove unused, add missing)
go mod tidy
```

### Dependency Management

#### View Dependencies
```bash
# List all dependencies
go list -m all

# Show dependency graph
go mod graph

# Show why a package is needed
go mod why github.com/package/name

# Check for available updates
go list -u -m all
```

#### Vendor Dependencies (Optional)
```bash
# Create vendor directory with all dependencies
go mod vendor

# Build using vendored dependencies
go build -mod=vendor -o citary-backend .

# Remove vendor directory
rm -rf vendor/
```

### Debugging

#### Debug Build
```bash
# Build with debug symbols
go build -gcflags="all=-N -l" -o citary-backend .

# Run with delve debugger
# Install: go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug main.go
```

#### Profiling
```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./...
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./...
go tool pprof mem.prof

# Run with profiling server (add to main.go)
# import _ "net/http/pprof"
# go func() { http.ListenAndServe(":6060", nil) }()
```

### Development Workflow

#### Typical Development Cycle
```bash
# 1. Pull latest changes
git pull origin main

# 2. Update dependencies
go mod download
go mod tidy

# 3. Make code changes
# ... edit files ...

# 4. Format code
go fmt ./...

# 5. Run linter
go vet ./...

# 6. Run tests
go test ./...

# 7. Build
go build -o citary-backend .

# 8. Run locally
./citary-backend
```

#### Pre-commit Checklist
```bash
# Format all code
go fmt ./...

# Run static analysis
go vet ./...

# Run tests with race detector
go test -race ./...

# Ensure no unused dependencies
go mod tidy

# Verify build
go build -o citary-backend .
```

### Useful Go Commands

```bash
# Show Go environment
go env

# Show Go version
go version

# Download module to cache
go mod download

# Verify dependencies
go mod verify

# List available packages
go list ./...

# Show package documentation
go doc package/name

# Generate code (if using go:generate)
go generate ./...
```

## 🏛️ Architecture

This project follows **Clean Architecture** (Hexagonal Architecture) principles:

### Layers

1. **Domain Layer** (`internal/domain/`)
   - Contains business entities, use cases, and repository interfaces
   - No external dependencies
   - Pure business logic

2. **Infrastructure Layer** (`internal/infrastructure/`)
   - Implements domain interfaces
   - Handles external concerns (database, HTTP, config)
   - Depends on domain layer

3. **Shared Layer** (`pkg/`)
   - Constants and utilities
   - Can be imported by any layer

### Dependency Direction
```
HTTP Handlers → Use Cases → Repository Interfaces
                                    ↑
                                    |
                        Repository Implementations
```

For detailed architecture documentation, see [ARCHITECTURE.md](ARCHITECTURE.md)

## 🔒 Security

- **Password Hashing** - bcrypt with default cost (10)
- **SQL Injection Protection** - Parameterized queries with context
- **CORS** - Configurable via middleware
- **Input Validation** - Comprehensive request validation
- **Panic Recovery** - Prevents server crashes and information leakage

## ⚙️ Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DATABASE_URL` | PostgreSQL connection string | - | ✅ Yes |
| `PORT` | HTTP server port | `3005` | ❌ No |

### Database Configuration

```go
MaxOpenConnections: 25
MaxIdleConnections: 5
```

### HTTP Server Timeouts

```go
ReadTimeout:  15 seconds
WriteTimeout: 15 seconds
IdleTimeout:  60 seconds
```

## 📊 Database Schema

### `data.data_user` Table

```sql
CREATE TABLE data.data_user (
    use_id                SERIAL PRIMARY KEY,
    use_email             VARCHAR(100) UNIQUE NOT NULL,
    use_password_hash     TEXT NOT NULL,
    use_email_verified    BOOLEAN DEFAULT FALSE,
    use_phone_verified    BOOLEAN DEFAULT FALSE,
    use_two_factor_enabled BOOLEAN DEFAULT FALSE,
    use_two_factor_secret TEXT,
    use_last_login        TIMESTAMP,
    use_login_attempts    INTEGER DEFAULT 0,
    use_locked_until      TIMESTAMP,
    use_terms_accepted_at TIMESTAMP,
    use_privacy_accepted_at TIMESTAMP,
    use_created_date      TIMESTAMP DEFAULT NOW(),
    use_record_status     VARCHAR(1) DEFAULT '0'
);
```

## 🎯 Roadmap

Future enhancements planned:

- [ ] Unit tests with mocks
- [ ] Integration tests
- [ ] User login endpoint
- [ ] JWT authentication & authorization
- [ ] Email verification flow
- [ ] Password reset functionality
- [ ] Refresh token mechanism
- [ ] Rate limiting middleware
- [ ] API versioning (v1, v2)
- [ ] OpenAPI/Swagger documentation
- [ ] Docker & Docker Compose support
- [ ] CI/CD pipeline (GitHub Actions)
- [ ] Monitoring & observability (Prometheus, Grafana)
- [ ] Distributed tracing (OpenTelemetry)

## 🧪 Testing

Testing structure (to be implemented):

```
internal/
├── domain/
│   └── usecases/
│       └── auth/
│           ├── signup_user.go
│           └── signup_user_test.go
└── infrastructure/
    └── http/
        └── handlers/
            └── auth/
                ├── auth_handler.go
                └── auth_handler_test.go
```

## 🤝 Contributing

1. Follow Go best practices and conventions
2. Maintain Clean Architecture principles
3. Add tests for new features
4. Update documentation as needed
5. Use conventional commit messages

## 📝 Code Style

This project follows:
- [Effective Go](https://golang.org/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber's Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

## 📚 Learn More

- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) - Robert C. Martin
- [Hexagonal Architecture](https://alistair.cockburn.us/hexagonal-architecture/) - Alistair Cockburn
- [Go Project Layout](https://github.com/golang-standards/project-layout) - Standard Go project structure

## 📄 License

[Add your license here]

## 👤 Author

[Add author information]

---

**Built with ❤️ using Go 1.24.3**
