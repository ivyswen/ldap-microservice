# LDAP Microservice

A lightweight Go microservice for LDAP authentication and user management. This service provides REST APIs for authenticating users against an LDAP directory server.

## Features

- **LDAP Authentication**: Authenticate users against LDAP/Active Directory servers
- **Flexible Configuration**: Support for LDAP, LDAPS, and StartTLS connections
- **Health Checks**: Built-in health and readiness probes for Kubernetes
- **Structured Logging**: Comprehensive logging using zerolog
- **Error Handling**: Detailed error types and messages for debugging
- **Connection Timeout**: Configurable connection and request timeouts
- **Service Account Binding**: Optional service account for user searches

## Requirements

- Go 1.25 or later
- LDAP/Active Directory server

## Building

### Local Build

```bash
go build -o ldap-svc
```

### Docker Build

```bash
docker build -t ldap-svc:latest .
```

## Running

### Local Execution

```bash
./ldap-svc
```

### Docker Execution

```bash
docker run -p 8080:8080 \
  -e LDAP_URL=ldaps://ldap.example.com:636 \
  -e LDAP_BIND_DN=cn=admin,dc=example,dc=com \
  -e LDAP_BIND_PASSWORD=password \
  ldap-svc:latest
```

## Configuration

Configuration is managed through environment variables. The service supports loading configuration from a `.env` file for easier development and testing.

### Using .env File (Recommended for Development)

The service automatically loads environment variables from a `.env` file if it exists in the working directory.

#### Quick Start with .env

1. **Copy the example file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your LDAP settings:**
   ```env
   LDAP_URL=ldap://ldap.example.com:389
   LDAP_BIND_DN=cn=admin,dc=example,dc=com
   LDAP_BIND_PASSWORD=password123
   LDAP_USER_BASE=dc=example,dc=com
   LDAP_USER_FILTER=(uid=%s)
   SERVICE_PORT=8080
   ```

3. **Run the service:**
   ```bash
   go run .
   # or
   ./ldap-microservice
   ```

The service will automatically load the configuration from `.env` file.

### Environment Variables

### LDAP Connection Settings

- `LDAP_URL` (default: `ldap://ldap.example.com:389`): LDAP server URL
- `LDAP_CONN_TIMEOUT` (default: `5s`): Connection timeout duration
- `LDAP_USE_LDAPS` (default: `0`): Use LDAPS (1=true, 0=false)
- `LDAP_USE_STARTTLS` (default: `0`): Use StartTLS (1=true, 0=false)
- `LDAP_INSECURE_SKIP_VERIFY` (default: `0`): Skip TLS verification (1=true, 0=false)

### LDAP Credentials

- `LDAP_BIND_DN`: Service account DN for searches (optional)
- `LDAP_BIND_PASSWORD`: Service account password (optional)

### User Search Configuration

- `LDAP_USER_BASE` (default: `dc=example,dc=com`): Base DN for user searches
- `LDAP_USER_FILTER` (default: `(uid=%s)`): LDAP filter for user searches
- `LDAP_USER_DN_ATTR`: Attribute name for user DN (optional, uses entry DN if not set)
- `LDAP_RETURN_ATTRIBUTES` (default: `uid,mail,cn`): Comma-separated list of attributes to return

### Service Configuration

- `SERVICE_PORT` (default: `8080`): HTTP server port
- `BASE_PATH` (default: empty): URL path prefix for all API endpoints
  - Examples:
    - `BASE_PATH=/api` → endpoints become `/api/v1/auth`, `/api/v1/healthz`, etc.
    - `BASE_PATH=/ldap` → endpoints become `/ldap/v1/auth`, `/ldap/v1/healthz`, etc.
    - `BASE_PATH=` (empty) → endpoints remain `/v1/auth`, `/v1/healthz`, etc.
  - Note: The prefix is automatically normalized (trailing `/` removed, leading `/` ensured)
- `LDAP_REQUEST_TIMEOUT` (default: `10s`): Request timeout duration
- `LOG_LEVEL` (default: `info`): Log level (debug, info, warn, error)

## API Endpoints

> **Note**: All endpoints below are shown without the `BASE_PATH` prefix. If you configure `BASE_PATH=/api`, prepend `/api` to all paths (e.g., `/api/v1/auth`).

### POST /v1/auth

Authenticate a user against LDAP.

**Request:**
```json
{
  "username": "john.doe",
  "password": "password123"
}
```

**Success Response (200):**
```json
{
  "ok": true,
  "user": {
    "uid": "john.doe",
    "mail": "john@example.com",
    "cn": "John Doe"
  }
}
```

**Error Response (401/500):**
```json
{
  "ok": false,
  "error": "invalid_credentials",
  "detail": "error details"
}
```

### GET /v1/healthz

Health check endpoint.

**Response (200):**
```json
{
  "status": "ok"
}
```

### GET /v1/readyz

Readiness check endpoint.

**Response (200):**
```json
{
  "ready": "true"
}
```

## Deployment

### Kubernetes

See `deploy/example.yaml` for a sample Kubernetes deployment manifest.

```bash
kubectl apply -f deploy/example.yaml
```

## Testing

### Unit Tests

Run unit tests:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

### API Testing

> **Note on BASE_PATH**: If you configured `BASE_PATH=/api`, replace `/v1/` with `/api/v1/` in all examples below.

#### Windows Environment

For comprehensive Windows testing guide, see [WINDOWS_TESTING_GUIDE.md](WINDOWS_TESTING_GUIDE.md).

Quick test with PowerShell:

```powershell
# Health check (without BASE_PATH)
Invoke-RestMethod -Uri "http://localhost:8080/v1/healthz" -Method Get

# Health check (with BASE_PATH=/api)
Invoke-RestMethod -Uri "http://localhost:8080/api/v1/healthz" -Method Get

# Authentication test
$body = @{
    username = "john.doe"
    password = "password123"
} | ConvertTo-Json

Invoke-RestMethod `
    -Uri "http://localhost:8080/v1/auth" `
    -Method Post `
    -ContentType "application/json" `
    -Body $body
```

#### Using Test Script

Run the provided PowerShell test script:

```powershell
# Run all tests
.\test-api.ps1

# Run tests against custom URL
.\test-api.ps1 -BaseUrl "http://localhost:9090"
```

#### Using curl

```bash
# Health check (without BASE_PATH)
curl http://localhost:8080/v1/healthz

# Health check (with BASE_PATH=/api)
curl http://localhost:8080/api/v1/healthz

# Authentication (without BASE_PATH)
curl -X POST http://localhost:8080/v1/auth \
  -H "Content-Type: application/json" \
  -d '{"username":"john.doe","password":"password123"}'

# Authentication (with BASE_PATH=/api)
curl -X POST http://localhost:8080/api/v1/auth \
  -H "Content-Type: application/json" \
  -d '{"username":"john.doe","password":"password123"}'
```

## Development

### Project Structure

- `main.go`: Application entry point
- `config.go`: Configuration management
- `handlers.go`: HTTP request handlers
- `ldapclient.go`: LDAP client implementation
- `errors.go`: Custom error types
- `config_test.go`: Config tests
- `ldapclient_test.go`: LDAP client tests
- `.env.example`: Example environment configuration file
- `WINDOWS_TESTING_GUIDE.md`: Windows testing guide
- `test-api.ps1`: PowerShell API test script
- `run.ps1`: PowerShell startup script

### Configuration Files

- `.env.example`: Template for environment variables (commit to repository)
- `.env`: Local environment configuration (ignored by git, create from .env.example)
- `.env.local`: Local overrides (optional, ignored by git)

### Quick Start for Development

1. **Copy environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` with your settings:**
   ```bash
   # Edit with your favorite editor
   nano .env
   # or
   code .env
   ```

3. **Run the service:**
   ```bash
   go run .
   ```

4. **Test the API:**
   ```powershell
   # Windows PowerShell
   .\test-api.ps1

   # Or manually
   Invoke-RestMethod -Uri "http://localhost:8080/v1/healthz" -Method Get
   ```

### Code Style

This project follows Go best practices:
- Uses `any` instead of `interface{}` (Go 1.18+)
- Implements proper error handling with custom error types
- Uses structured logging with zerolog
- Includes comprehensive unit tests
- Supports `.env` file for configuration (using godotenv)

### Windows Development

For detailed Windows development and testing instructions, see [WINDOWS_TESTING_GUIDE.md](WINDOWS_TESTING_GUIDE.md).

Key features:
- `.env` file support for easy configuration
- PowerShell test scripts for API testing
- PowerShell startup script with build and test options
- Comprehensive Windows testing guide

## License

MIT
