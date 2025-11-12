# Repository Integration Tests

## Setup

Las pruebas de integración de repositorios requieren una base de datos PostgreSQL de prueba.

### Opción 1: Usar Docker (Recomendado)

```bash
# Levantar PostgreSQL en Docker
docker run --name postgres-test \
  -e POSTGRES_USER=test_user \
  -e POSTGRES_PASSWORD=test_password \
  -e POSTGRES_DB=citary_test \
  -p 5433:5432 \
  -d postgres:15-alpine

# Ejecutar las pruebas
TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5433/citary_test?sslmode=disable" go test -v ./test/integration/repositories/...

# Detener y eliminar el contenedor
docker stop postgres-test
docker rm postgres-test
```

### Opción 2: Usar PostgreSQL local

1. Crear una base de datos de prueba:
```sql
CREATE DATABASE citary_test;
CREATE USER test_user WITH PASSWORD 'test_password';
GRANT ALL PRIVILEGES ON DATABASE citary_test TO test_user;
```

2. Ejecutar el script de migración en la base de datos de prueba

3. Ejecutar las pruebas:
```bash
TEST_DATABASE_URL="postgres://test_user:test_password@localhost:5432/citary_test?sslmode=disable" go test -v ./test/integration/repositories/...
```

## Notas

- Las pruebas de integración de repositorios están **deshabilitadas por defecto**
- Solo se ejecutan si la variable de entorno `TEST_DATABASE_URL` está configurada
- Se recomienda usar una base de datos separada para pruebas
- Las pruebas limpian los datos después de cada ejecución
