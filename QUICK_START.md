# Quick Start Guide

A quick reference for common tasks in Citary Backend.

---

## ⚡ Quick Commands

### First Time Setup
```bash
git clone <repo-url>
cd citary-backend
cp .env.example .env
# Edit .env with your DATABASE_URL
go mod download
go build -o citary-backend .
./citary-backend
```

### Daily Development
```bash
# Pull latest + update dependencies
git pull && go mod tidy

# Run without building
go run main.go

# Build and run
go build -o citary-backend . && ./citary-backend

# Format + Lint + Test + Build
go fmt ./... && go vet ./... && go test ./... && go build .
```

### Common Tasks

| Task | Command |
|------|---------|
| **Run** | `go run main.go` |
| **Build** | `go build -o citary-backend .` |
| **Test** | `go test ./...` |
| **Format** | `go fmt ./...` |
| **Lint** | `go vet ./...` |
| **Clean** | `go clean && rm -f citary-backend*` |
| **Add dependency** | `go get github.com/package/name` |
| **Update deps** | `go get -u ./... && go mod tidy` |
| **Coverage** | `go test -cover ./...` |

---

## 📁 Project Layout Quick Reference

```
internal/
├── domain/              # Business logic (no external dependencies)
│   ├── entities/        # Business entities
│   ├── repositories/    # Interfaces (ports)
│   ├── usecases/        # Business logic
│   ├── dtos/           # Request/response data
│   └── errors/         # Domain errors
└── infrastructure/     # External adapters
    ├── persistence/    # Database
    ├── http/          # Web server
    ├── config/        # Configuration
    └── di/            # Dependency injection

pkg/constants/          # Shared constants
```

---

## 🎯 Adding a New Feature

### Example: Add Login Endpoint

1. **Create DTO** (`internal/domain/dtos/auth/login_request.go`)
   ```go
   package auth

   type LoginRequest struct {
       Email    string `json:"email"`
       Password string `json:"password"`
   }
   ```

2. **Create Use Case** (`internal/domain/usecases/auth/login_user.go`)
   ```go
   package auth

   func (uc *LoginUserUseCase) Execute(ctx context.Context, dto LoginRequest) (*LoginResponse, error) {
       // Business logic
   }
   ```

3. **Create Handler** (`internal/infrastructure/http/handlers/auth/auth_handler.go`)
   ```go
   func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
       // HTTP handling
   }
   ```

4. **Add Route** (`internal/infrastructure/http/router/router.go`)
   ```go
   mux.HandleFunc("/auth/login", rt.authHandler.LoginUser)
   ```

5. **Wire in DI** (`internal/infrastructure/di/container.go`)
   ```go
   loginUseCase := auth.NewLoginUserUseCase(userRepository)
   ```

---

## 🧪 Testing Quick Reference

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Specific package
go test ./internal/domain/usecases/auth/

# With race detector
go test -race ./...

# Generate HTML coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

---

## 🐛 Debugging

```bash
# Build with debug info
go build -gcflags="all=-N -l" -o citary-backend .

# Run with race detector
go run -race main.go

# Use delve debugger
dlv debug main.go
```

---

## 🚀 Build for Production

```bash
# Optimized build (smaller binary)
go build -ldflags="-s -w" -o citary-backend .

# With version info
VERSION=$(git describe --tags --always)
go build -ldflags="-X main.Version=$VERSION -s -w" -o citary-backend .

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o citary-backend .
```

---

## 📦 Dependency Management

```bash
# Add new package
go get github.com/package/name

# Add specific version
go get github.com/package/name@v1.2.3

# Update all dependencies
go get -u ./...
go mod tidy

# View all dependencies
go list -m all

# Check for updates
go list -u -m all

# Why is package X needed?
go mod why github.com/package/name
```

---

## 🧹 Cleanup

```bash
# Remove binaries
rm -f citary-backend citary-backend.exe

# Clean Go cache
go clean -cache -testcache

# Full cleanup (careful!)
go clean -i -r -cache -testcache -modcache
```

---

## 🔥 Troubleshooting

### "Cannot find package"
```bash
go mod download
go mod tidy
```

### "Build fails"
```bash
go clean -cache
go mod verify
go build -v .
```

### "Tests fail"
```bash
go clean -testcache
go test -v ./...
```

### "Import cycle"
- Check dependency direction
- Domain should not import infrastructure
- See ARCHITECTURE.md for layer rules

---

## 📝 Environment Variables

```env
# Required
DATABASE_URL=postgres://user:pass@localhost:5432/citary?sslmode=disable

# Optional
PORT=3001  # Default: 3001
```

---

## 🔗 Useful Links

- [Full README](README.md) - Complete documentation
- [Architecture Guide](ARCHITECTURE.md) - Detailed architecture
- [Effective Go](https://golang.org/doc/effective_go) - Go best practices
- [Project Layout](https://github.com/golang-standards/project-layout) - Standard structure

---

**Need help?** Check the full [README.md](README.md) or [ARCHITECTURE.md](ARCHITECTURE.md)
