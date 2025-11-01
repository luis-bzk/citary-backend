# Architecture Documentation

## Clean Architecture Overview

This project follows **Clean Architecture** (also known as **Hexagonal Architecture** or **Ports and Adapters**) principles.

---

## Layer Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                         HTTP Layer                              │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │  Middleware  │→ │   Handlers   │→ │   Response   │         │
│  │  - CORS      │  │   - Auth     │  │   Helpers    │         │
│  │  - Logging   │  │              │  │              │         │
│  │  - Recovery  │  │              │  │              │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                      Use Case Layer                             │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              Business Logic (Use Cases)                   │  │
│  │  - SignupUserUseCase                                      │  │
│  │  - Validation                                             │  │
│  │  - Business Rules                                         │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
                            ↓
┌─────────────────────────────────────────────────────────────────┐
│                    Domain Layer (Core)                          │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐         │
│  │   Entities   │  │ Repositories │  │     DTOs     │         │
│  │   - User     │  │  (Interfaces)│  │  - SignupReq │         │
│  │              │  │              │  │              │         │
│  └──────────────┘  └──────────────┘  └──────────────┘         │
└─────────────────────────────────────────────────────────────────┘
                            ↑
┌─────────────────────────────────────────────────────────────────┐
│                 Infrastructure Layer                            │
│  ┌──────────────────────────────────────────────────────────┐  │
│  │              Repository Implementations                   │  │
│  │  - UserRepositoryImpl (PostgreSQL)                        │  │
│  │  - Mappers (Domain ↔ DB)                                  │  │
│  │  - Database Connection                                    │  │
│  └──────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

---

## Dependency Flow

```
main.go
  │
  ├─→ DI Container (Initializes everything)
  │     │
  │     ├─→ Config Loader
  │     ├─→ Database Connection
  │     ├─→ Repositories
  │     ├─→ Use Cases
  │     ├─→ Handlers
  │     ├─→ Router
  │     └─→ HTTP Server
  │
  └─→ Graceful Shutdown Handler
```

---

## Request Flow (Signup Example)

```
1. HTTP Request arrives
   POST /auth/signup
   Body: {"email": "user@example.com", "password": "Pass123!"}

   ↓

2. Middleware Chain
   Recovery → CORS → Logging → Route Handler

   ↓

3. AuthHandler.SignupUser()
   - Validates HTTP method
   - Decodes JSON into SignupRequest DTO
   - Calls use case with context

   ↓

4. SignupUserUseCase.Execute(ctx, dto)
   - Validates input data
   - Checks if user exists (via repository)
   - Hashes password
   - Creates User entity
   - Persists user (via repository)
   - Returns User entity

   ↓

5. UserRepositoryImpl.Create(ctx, user)
   - Maps domain entity → database entity
   - Executes SQL INSERT
   - Returns generated ID

   ↓

6. Handler maps entity → SignupResponse DTO

   ↓

7. Response Helper formats JSON
   {
     "success": true,
     "message": "User created successfully",
     "data": {
       "id": 1,
       "email": "user@example.com",
       "emailVerified": false,
       "createdDate": "2024-10-31T..."
     }
   }

   ↓

8. HTTP Response sent (201 Created)
```

---

## Package Dependencies

### Domain Layer (No External Dependencies)
```
internal/domain/
├── entities/        (depends on: pkg/constants)
├── repositories/    (depends on: entities, context)
├── usecases/        (depends on: entities, repositories, dtos, errors)
├── dtos/            (depends on: nothing - pure validation)
└── errors/          (depends on: pkg/constants)
```

### Infrastructure Layer (Implements Domain Interfaces)
```
internal/infrastructure/
├── persistence/
│   └── postgres/
│       ├── entities/       (depends on: database/sql)
│       ├── mappers/        (depends on: domain/entities, postgres/entities)
│       └── repositories/   (depends on: domain/repositories, mappers)
│
├── http/
│   ├── handlers/           (depends on: usecases, http/dto, response)
│   ├── middleware/         (depends on: net/http)
│   ├── dto/                (depends on: nothing)
│   ├── response/           (depends on: domain/errors, http/dto)
│   └── router/             (depends on: handlers, middleware)
│
├── config/                 (depends on: godotenv)
└── di/                     (depends on: ALL layers - wires everything)
```

---

## Key Principles Applied

### 1. **Dependency Inversion**
- High-level modules (use cases) don't depend on low-level modules (repositories)
- Both depend on abstractions (interfaces)

### 2. **Single Responsibility**
- Each package has one reason to change
- Handlers: HTTP concerns
- Use cases: Business logic
- Repositories: Data access

### 3. **Interface Segregation**
- Small, focused interfaces (UserRepository)
- Clients depend only on methods they use

### 4. **Open/Closed Principle**
- Easy to add new features without modifying existing code
- Example: Add MongoDB repository by implementing UserRepository interface

---

## Testing Strategy (Recommended)

### Unit Tests
```go
// Use case tests with mock repository
func TestSignupUserUseCase_Execute(t *testing.T) {
    mockRepo := &MockUserRepository{}
    useCase := auth.NewSignupUserUseCase(mockRepo)

    mockRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, nil)
    mockRepo.On("Create", ctx, mock.Anything).Return(nil)

    user, err := useCase.Execute(ctx, validSignupRequest)

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Integration Tests
```go
// Handler tests with real use case and mock repository
func TestAuthHandler_SignupUser(t *testing.T) {
    req := httptest.NewRequest("POST", "/auth/signup", body)
    w := httptest.NewRecorder()

    handler.SignupUser(w, req)

    assert.Equal(t, http.StatusCreated, w.Code)
}
```

### E2E Tests
```go
// Full stack tests with test database
func TestE2E_UserSignup(t *testing.T) {
    testDB := setupTestDatabase()
    defer testDB.Cleanup()

    resp := callAPI("POST", "/auth/signup", signupPayload)

    assert.Equal(t, 201, resp.StatusCode)
    assertUserExistsInDB(testDB, "test@example.com")
}
```

---

## Adding New Features

### Example: Adding Login Endpoint

1. **Define DTO** (`internal/domain/dtos/auth/login_request.go`)
   ```go
   type LoginRequest struct {
       Email    string
       Password string
   }
   ```

2. **Create Use Case** (`internal/domain/usecases/auth/login_user.go`)
   ```go
   type LoginUserUseCase struct {
       userRepository repositories.UserRepository
   }

   func (uc *LoginUserUseCase) Execute(ctx context.Context, dto LoginRequest) (*LoginResponse, error) {
       // Business logic here
   }
   ```

3. **Create Handler** (`internal/infrastructure/http/handlers/auth/auth_handler.go`)
   ```go
   func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
       // HTTP handling here
   }
   ```

4. **Add Route** (`internal/infrastructure/http/router/router.go`)
   ```go
   mux.HandleFunc("/auth/login", rt.authHandler.LoginUser)
   ```

5. **Wire in DI** (`internal/infrastructure/di/container.go`)
   ```go
   loginUserUseCase := auth.NewLoginUserUseCase(userRepository)
   authHandlerInstance := authHandler.NewAuthHandler(signupUserUseCase, loginUserUseCase)
   ```

---

## Configuration

### Environment Variables
```env
DATABASE_URL=postgres://user:pass@localhost:5432/dbname?sslmode=disable
PORT=3005
```

### Database Connection Pool
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

---

## Security Considerations

1. **Password Hashing:** bcrypt with default cost (10)
2. **SQL Injection Protection:** Parameterized queries
3. **CORS:** Configurable via middleware
4. **Panic Recovery:** Prevents server crashes
5. **Context Timeouts:** Prevents long-running operations

---

## Performance Optimizations

1. **Connection Pooling:** Database connections reused
2. **Graceful Shutdown:** In-flight requests complete before shutdown
3. **Structured Logging:** Efficient log formatting
4. **Minimal Dependencies:** Fast build and startup times

---

## Monitoring & Observability

### Current Logging
- Request logging (method, path, IP, status, duration)
- Database connection status
- Error logging
- Startup/shutdown events

### Recommended Additions
- Structured logging (JSON format)
- Request ID tracking
- Metrics (Prometheus)
- Distributed tracing (OpenTelemetry)
- Health checks (liveness/readiness probes)

---

**Last Updated:** October 31, 2024
