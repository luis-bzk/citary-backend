# Testing Guide - Citary Backend

Este proyecto sigue una estrategia de testing completa basada en **TDD (Test-Driven Development)** con pruebas organizadas por tipo y responsabilidad.

## 📁 Estructura de Carpetas

```
test/
├── unit/                          # Pruebas unitarias (sin dependencias externas)
│   ├── dtos/
│   │   └── auth/                  # Pruebas de validación de DTOs
│   ├── mappers/                   # Pruebas de conversión Domain ↔ DB
│   └── usecases/
│       └── auth/                  # Pruebas de lógica de negocio con mocks
├── integration/                   # Pruebas de integración (múltiples componentes)
│   ├── handlers/
│   │   └── auth/                  # Pruebas HTTP con httptest
│   └── repositories/              # Pruebas con base de datos real
└── mocks/
    └── repositories/              # Mocks generados con testify/mock
```

## 🧪 Tipos de Pruebas

### 1. **Pruebas Unitarias de DTOs** ⚠️ CRÍTICA
- **Ubicación:** `test/unit/dtos/auth/`
- **Qué prueban:** Validaciones de entrada (email, password, etc.)
- **Dependencias:** Ninguna
- **Cobertura:** 27 pruebas

```bash
go test -v ./test/unit/dtos/auth/...
```

### 2. **Pruebas Unitarias de Mappers** 🟡 ALTA
- **Ubicación:** `test/unit/mappers/`
- **Qué prueban:** Conversión correcta entre entidades de dominio y base de datos
- **Dependencias:** Ninguna
- **Cobertura:** 7 pruebas

```bash
go test -v ./test/unit/mappers/...
```

### 3. **Pruebas Unitarias de Use Cases** ⚠️ CRÍTICA
- **Ubicación:** `test/unit/usecases/auth/`
- **Qué prueban:** Lógica de negocio (signup, validaciones, hash de password)
- **Dependencias:** Mocks del repositorio
- **Cobertura:** 12 pruebas

```bash
go test -v ./test/unit/usecases/auth/...
```

### 4. **Pruebas de Integración de Handlers** 🟡 ALTA
- **Ubicación:** `test/integration/handlers/auth/`
- **Qué prueban:** Capa HTTP completa (parsing JSON, status codes, respuestas)
- **Dependencias:** httptest + mocks del repositorio
- **Cobertura:** 11 pruebas

```bash
go test -v ./test/integration/handlers/auth/...
```

### 5. **Pruebas de Integración de Repositories** ⚠️ CRÍTICA
- **Ubicación:** `test/integration/repositories/`
- **Qué prueban:** Queries SQL reales contra PostgreSQL
- **Dependencias:** Base de datos PostgreSQL de prueba
- **Cobertura:** 9 pruebas (se saltan si no hay DB configurada)

```bash
# Sin base de datos configurada (se saltan las pruebas)
go test -v ./test/integration/repositories/...

# Con base de datos configurada
TEST_DATABASE_URL="postgres://user:pass@localhost:5433/citary_test?sslmode=disable" go test -v ./test/integration/repositories/...
```

Ver [test/integration/repositories/README.md](integration/repositories/README.md) para configurar la base de datos de prueba.

## 🚀 Comandos de Ejecución

### Ejecutar TODAS las pruebas

```bash
# Todas las pruebas (incluye solo las que no requieren DB)
go test -v ./test/...

# Con cobertura
go test -v -cover ./test/...

# Con reporte de cobertura detallado
go test -v -coverprofile=coverage.out ./test/...
go tool cover -html=coverage.out
```

### Ejecutar pruebas por tipo

```bash
# Solo pruebas unitarias
go test -v ./test/unit/...

# Solo pruebas de integración (sin DB)
go test -v ./test/integration/handlers/...

# Pruebas de repositorios (con DB)
TEST_DATABASE_URL="postgres://..." go test -v ./test/integration/repositories/...
```

### Ejecutar pruebas específicas

```bash
# Una prueba específica
go test -v -run TestSignupRequest_Validate_Success ./test/unit/dtos/auth/

# Todas las pruebas que contengan "Email" en el nombre
go test -v -run Email ./test/...

# Pruebas de un paquete específico
go test -v ./test/unit/usecases/auth/
```

### Modo watch (ejecutar al guardar cambios)

```bash
# Instalar nodemon o similar
npm install -g nodemon

# Ejecutar en modo watch
nodemon --exec "go test -v ./test/..." --ext go
```

## 📊 Resumen de Cobertura

| Componente | Tipo | Pruebas | Prioridad |
|-----------|------|---------|-----------|
| DTOs (Validate) | Unitaria | 27 | ⚠️ CRÍTICA |
| Mappers | Unitaria | 7 | 🟡 ALTA |
| Use Cases | Unitaria con mocks | 12 | ⚠️ CRÍTICA |
| Handlers | Integración con httptest | 11 | 🟡 ALTA |
| Repositories | Integración con DB | 9 | ⚠️ CRÍTICA |
| **TOTAL** | | **66** | |

## 🎯 Filosofía TDD Aplicada

Las pruebas están organizadas siguiendo TDD (Test-Driven Development):

1. **Red Phase (Casos de error primero):**
   - Email vacío
   - Email con formato inválido
   - Password sin mayúsculas
   - Usuario ya existe
   - Errores de base de datos

2. **Green Phase (Casos exitosos):**
   - Validación correcta
   - Usuario creado exitosamente
   - Diferentes formatos válidos de email
   - Password hasheado correctamente

3. **Refactor:**
   - Tests organizados por funcionalidad
   - Nombres descriptivos
   - Setup y Cleanup claros

## 🔧 Herramientas de Testing

### Instaladas en el proyecto

```go
require (
    github.com/stretchr/testify v1.11.1  // Assertions + Mocks
)
```

### Testify/Assert - Assertions
```go
assert.NoError(t, err)
assert.Equal(t, expected, actual)
assert.Contains(t, haystack, needle)
```

### Testify/Mock - Mocking
```go
mockRepo := new(MockUserRepository)
mockRepo.On("FindByEmail", ctx, "test@example.com").Return(nil, nil)
mockRepo.AssertExpectations(t)
```

### Httptest - HTTP Testing
```go
req := httptest.NewRequest(http.MethodPost, "/auth/signup", body)
w := httptest.NewRecorder()
handler.SignupUser(w, req)
assert.Equal(t, http.StatusCreated, w.Code)
```

## 📝 Convenciones de Nombres

### Formato de nombres de pruebas
```
Test<ComponentName>_<MethodName>_<Scenario>
```

Ejemplos:
- `TestSignupRequest_Validate_EmailEmpty`
- `TestUserMapper_ToDomainEntity_WithNullFields`
- `TestSignupUserUseCase_Execute_UserAlreadyExists`
- `TestAuthHandler_SignupUser_InvalidJSON`

### Estructura AAA (Arrange-Act-Assert)
```go
func TestExample(t *testing.T) {
    // Arrange - Preparar datos y mocks
    mockRepo := new(MockUserRepository)
    useCase := NewSignupUserUseCase(mockRepo)

    // Act - Ejecutar la función bajo prueba
    result, err := useCase.Execute(ctx, request)

    // Assert - Verificar resultados
    assert.NoError(t, err)
    assert.NotNil(t, result)
}
```

## 🐛 Debugging de Pruebas

### Ver output detallado
```bash
go test -v ./test/...
```

### Ejecutar una sola prueba
```bash
go test -v -run TestNombreEspecifico ./test/unit/dtos/auth/
```

### Ver logs de la aplicación en pruebas
```go
t.Logf("Debug info: %v", someVariable)
```

### Imprimir cobertura por función
```bash
go test -coverprofile=coverage.out ./test/...
go tool cover -func=coverage.out
```

## 🔄 CI/CD Integration

### GitHub Actions (ejemplo)
```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:15-alpine
        env:
          POSTGRES_USER: test_user
          POSTGRES_PASSWORD: test_password
          POSTGRES_DB: citary_test
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.24'
      - run: go test -v -cover ./test/...
      - run: TEST_DATABASE_URL="postgres://..." go test -v ./test/integration/repositories/...
```

## 📚 Referencias

- [Testing en Go - Documentación oficial](https://golang.org/pkg/testing/)
- [Testify - GitHub](https://github.com/stretchr/testify)
- [Clean Architecture Testing](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [TDD Best Practices](https://martinfowler.com/bliki/TestDrivenDevelopment.html)

## ❓ Preguntas Frecuentes

### ¿Por qué las pruebas de repositorios se saltan por defecto?

Las pruebas de integración de repositorios requieren una base de datos PostgreSQL real. Para evitar fallos en desarrollo local, se saltan automáticamente si no está configurada la variable `TEST_DATABASE_URL`.

### ¿Cómo ejecuto solo las pruebas rápidas?

```bash
go test -short ./test/...
```

Las pruebas lentas (como las de integración con DB) deben marcar `testing.Short()` para saltarse.

### ¿Necesito mockear todo?

No. Solo mockea dependencias externas (base de datos, APIs). Las dependencias internas del dominio (entities, DTOs) deben usarse reales.

### ¿Cuándo debo agregar nuevas pruebas?

- Antes de agregar nueva funcionalidad (TDD)
- Cuando encuentres un bug (escribe primero el test que expone el bug)
- Cuando refactorices código crítico

## 🎉 ¡Listo para probar!

Ejecuta todas las pruebas con:
```bash
go test -v ./test/...
```

**Resultado esperado:** ✅ 57 pruebas pasando (66 total, 9 saltadas sin DB)
