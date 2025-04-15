# Flow Wallet API - Suite de Pruebas HTTP

Esta suite de pruebas está diseñada para validar la funcionalidad de la API de Flow Wallet a través de solicitudes HTTP.

## Requisitos

- Docker y docker-compose instalados
- Cliente HTTP compatible con archivos .http (como REST Client en VSCode, Insomnia, o Postman)
- Variables de entorno configuradas (especialmente `FLOW_WALLET_ADMIN_ADDRESS`)

## Cómo usar REST Client en VSCode

1. **Instala la extensión REST Client**:
   - Busca "REST Client" en la pestaña de extensiones de VSCode e instálala
   - Reinicia VSCode si es necesario

2. **Abre el archivo de pruebas**:
   - Abre el archivo `test-suite/run-tests.http` en VSCode

3. **Ejecuta peticiones**:
   - Verás el texto "Send Request" encima de cada petición HTTP
   - Haz clic en "Send Request" para ejecutar esa petición específica
   - También puedes usar el atajo `Ctrl+Alt+R` (`Cmd+Alt+R` en Mac)

4. **Visualiza respuestas**:
   - Las respuestas aparecerán en una nueva pestaña a la derecha
   - Las variables como `{{response.body.jobId}}` se actualizarán automáticamente

5. **Ejecuta las pruebas en orden**:
   - Es importante seguir el orden de las pruebas ya que algunas dependen de los resultados de otras
   - Para pruebas más específicas, puedes usar los archivos individuales en cada subcarpeta

## Entorno de Pruebas Aislado

Para ejecutar estas pruebas sin afectar los datos de producción, se ha creado un entorno Docker aislado:

```bash
# Iniciar el entorno de pruebas
make -f test.mk test

# Ver logs del entorno de pruebas
make -f test.mk test-logs

# Limpiar completamente el entorno después de las pruebas
make -f test.mk test-clean
```

El entorno de pruebas incluye:
- API en puerto `3001`
- Base de datos PostgreSQL en puerto `5433`
- Emulador de Flow en puerto `3570`
- Redis en puerto `6380`

Todos los servicios usan volúmenes y nombres distintos para evitar conflictos con el entorno de producción.

## Estructura de la Suite

```
test-suite/
├── README.md
├── run-tests.http           # Archivo principal con pruebas básicas
├── variables.http           # Variables compartidas (ahora integradas en run-tests.http)
├── accounts/                # Pruebas detalladas de gestión de cuentas
├── health/                  # Pruebas detalladas de estado del servicio
├── scripts/                 # Pruebas detalladas de ejecución de scripts
├── system/                  # Pruebas detalladas de configuración del sistema
├── tokens/                  # Pruebas detalladas de tokens y transferencias
└── transactions/            # Pruebas detalladas de transacciones
```

## Solución de problemas

- **Variables no definidas**: Si ves errores sobre variables no definidas, asegúrate de que las peticiones se ejecutan en el orden correcto.
- **Error "Header name must be a valid HTTP token"**: Este error puede ocurrir si hay una directiva `@import` dentro de una petición HTTP. Ahora hemos integrado todas las variables directamente en `run-tests.http`.
- **Servidor no disponible**: Verifica que el entorno de pruebas está corriendo con `docker ps` o `make -f test.mk test-logs`.

## Flujo de Trabajo Recomendado

1. Inicia el entorno de pruebas con `make -f test.mk test`
2. Abre `test-suite/run-tests.http` en VSCode con la extensión REST Client
3. Ejecuta las pruebas en orden, una por una
4. Limpia todo con `make -f test.mk test-clean` cuando termines 