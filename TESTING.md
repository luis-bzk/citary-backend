# 🧪 Guía Rápida de Testing - Citary Backend

## ✅ Resumen de Implementación

Se implementaron **66 pruebas** distribuidas en 5 categorías siguiendo **TDD (Test-Driven Development)** y principios de **Clean Architecture**.

### 📊 Cobertura de Pruebas

| Componente | Tipo | Pruebas | Archivo | Prioridad |
|-----------|------|---------|---------|-----------|
| **DTOs** | Unitaria | 27 | [test/unit/dtos/auth/signup_request_test.go](test/unit/dtos/auth/signup_request_test.go) | ⚠️ CRÍTICA |
| **Mappers** | Unitaria | 7 | [test/unit/mappers/user_mapper_test.go](test/unit/mappers/user_mapper_test.go) | 🟡 ALTA |
| **Use Cases** | Unitaria + Mocks | 12 | [test/unit/usecases/auth/signup_user_test.go](test/unit/usecases/auth/signup_user_test.go) | ⚠️ CRÍTICA |
| **Handlers** | Integración + httptest | 11 | [test/integration/handlers/auth/auth_handler_test.go](test/integration/handlers/auth/auth_handler_test.go) | 🟡 ALTA |
| **Repositories** | Integración + DB | 9 | [test/integration/repositories/user_repository_test.go](test/integration/repositories/user_repository_test.go) | ⚠️ CRÍTICA |

**Total:** 66 pruebas (57 activas + 9 condicionales)

---

## 🚀 Comandos Rápidos

### Usando Makefile (Recomendado)

```bash
# Ver todos los comandos disponibles
make help

# Ejecutar todas las pruebas
make test

# Solo pruebas unitarias (rápidas)
make test-unit

# Solo pruebas de integración (sin DB)
make test-integration

# Generar reporte de cobertura HTML
make test-coverage

# Ejecutar con race detector
make test-race
```

### Usando go test directamente

```bash
# Todas las pruebas
go test -v ./test/...

# Solo unitarias
go test -v ./test/unit/...

# Solo integración de handlers
go test -v ./test/integration/handlers/...

# Con cobertura
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

---

## 📁 Estructura de Carpetas de Testing

```
test/
├── unit/                           # Pruebas sin dependencias externas
│   ├── dtos/
│   │   └── auth/
│   │       └── signup_request_test.go        (27 pruebas)
│   ├── mappers/
│   │   └── user_mapper_test.go               (7 pruebas)
│   └── usecases/
│       └── auth/
│           └── signup_user_test.go           (12 pruebas)
├── integration/                    # Pruebas de múltiples componentes
│   ├── handlers/
│   │   └── auth/
│   │       └── auth_handler_test.go          (11 pruebas)
│   └── repositories/
│       ├── README.md                         (Instrucciones de DB)
│       └── user_repository_test.go           (9 pruebas)
├── mocks/
│   └── repositories/
│       └── user_repository_mock.go           (Mock para tests)
└── README.md                                 (Documentación completa)
```

---

## 🎯 Enfoque TDD Implementado

Todas las pruebas siguen el patrón **Red → Green → Refactor**:

### 1. **Red Phase (Casos de error primero)**
```go
✗ Email vacío
✗ Email con formato inválido
✗ Password sin mayúsculas
✗ Usuario ya existe
✗ JSON malformado
```

### 2. **Green Phase (Casos exitosos)**
```go
✓ Validación correcta
✓ Usuario creado exitosamente
✓ Password hasheado correctamente
✓ Respuesta HTTP 201 con formato correcto
```

### 3. **Refactor**
- Código DRY (Don't Repeat Yourself)
- Nombres descriptivos
- Setup/Cleanup claros

---

## 📦 Dependencias de Testing

```go
require (
    github.com/stretchr/testify v1.11.1  // Assertions + Mocks
    github.com/lib/pq v1.10.9            // PostgreSQL driver
)
```

**Total de dependencias:** Solo 5 (4 indirectas)

---

## 🔍 Ejemplos de Pruebas

### Prueba Unitaria de DTO (TDD - Error primero)
```go
func TestSignupRequest_Validate_EmailEmpty(t *testing.T) {
    // Arrange - Caso de error
    request := auth.SignupRequest{
        Email:    "",
        Password: "ValidPass123!",
    }

    // Act
    err := request.Validate()

    // Assert
    assert.Error(t, err)
    assert.Equal(t, auth.ErrEmailEmpty, err)
}
```

### Prueba de Use Case con Mock
```go
func TestSignupUserUseCase_Execute_Success(t *testing.T) {
    // Arrange
    mockRepository := new(mockRepo.MockUserRepository)
    useCase := authUseCase.NewSignupUserUseCase(mockRepository)

    mockRepository.On("FindByEmail", ctx, "test@example.com").
        Return(nil, errors.ErrNotFound("User not found"))

    mockRepository.On("Create", ctx, mock.AnythingOfType("*entities.User")).
        Return(nil)

    // Act
    user, err := useCase.Execute(ctx, validRequest)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.NotEmpty(t, user.PasswordHash)
    mockRepository.AssertExpectations(t)
}
```

### Prueba de Integración de Handler
```go
func TestAuthHandler_SignupUser_Success(t *testing.T) {
    // Arrange
    handler := setupHandlerWithMocks()

    payload := map[string]string{
        "email":    "test@example.com",
        "password": "ValidPass123!",
    }
    body, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/auth/signup", bytes.NewBuffer(body))
    w := httptest.NewRecorder()

    // Act
    handler.SignupUser(w, req)

    // Assert
    assert.Equal(t, http.StatusCreated, w.Code)
    assert.Contains(t, w.Body.String(), "test@example.com")
}
```

---

## 🗄️ Pruebas con Base de Datos

Las pruebas de repositorios se **saltan automáticamente** si no hay base de datos configurada.

### Setup rápido con Docker

```bash
# 1. Levantar PostgreSQL
docker run --name postgres-test \
  -e POSTGRES_USER=test_user \
  -e POSTGRES_PASSWORD=test_password \
  -e POSTGRES_DB=citary_test \
  -p 5433:5432 \
  -d postgres:15-alpine

# 2. Ejecutar migraciones (crear tabla data.data_user)

# 3. Ejecutar pruebas
TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5433/citary_test?sslmode=disable" \
  go test -v ./test/integration/repositories/...

# 4. Cleanup
docker stop postgres-test && docker rm postgres-test
```

Ver [test/integration/repositories/README.md](test/integration/repositories/README.md) para más opciones.

---

## 📈 Ejecutar con Cobertura

```bash
# Generar reporte de cobertura
make test-coverage

# O manualmente
go test -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out -o coverage.html

# Abrir coverage.html en el navegador
```

---

## ✅ Verificación Rápida

```bash
# Ejecutar todas las pruebas y verificar que pasen
go test ./test/...

# Resultado esperado:
# ✅ 57 pruebas pasando
# ⏭️ 9 pruebas saltadas (repositorios sin DB)
# ❌ 0 pruebas fallidas
```

---

## 🎓 Aprende Más

- [Test README completo](test/README.md) - Documentación detallada
- [Repository Tests README](test/integration/repositories/README.md) - Setup de DB
- [Testing en Go](https://golang.org/pkg/testing/) - Docs oficiales
- [Testify](https://github.com/stretchr/testify) - Librería de testing

---

## 🐛 Troubleshooting

### Las pruebas no encuentran los paquetes
```bash
go mod tidy
go mod download
```

### Error: "package not found"
```bash
# Asegúrate de estar en la raíz del proyecto
cd d:\Documents\citary\citary-backend
go test ./test/...
```

### Las pruebas de DB se ejecutan pero fallan
```bash
# Verifica que la tabla existe
psql $TEST_DATABASE_URL -c "\dt data.*"

# Si no existe, ejecuta las migraciones
```

---

## 📝 Convenciones

### Nombres de pruebas
```
Test<ComponentName>_<MethodName>_<Scenario>
```

Ejemplos:
- `TestSignupRequest_Validate_EmailEmpty`
- `TestUserMapper_ToDomainEntity_WithNullFields`
- `TestSignupUserUseCase_Execute_Success`

### Estructura AAA
- **Arrange:** Preparar datos y mocks
- **Act:** Ejecutar función bajo prueba
- **Assert:** Verificar resultados

---

**¿Dudas?** Revisa la documentación completa en [test/README.md](test/README.md)
