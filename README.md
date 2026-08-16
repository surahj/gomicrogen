# gomicrogen

A powerful CLI tool to scaffold Go microservice projects with a predefined folder structure and production-ready templates.

## 🚀 Features

- **Service Types**: `--type casino` or `--type payment` scaffolds the wiring that family
  needs; no flag gives a neutral service. Adding a type is a new directory, not a code change
- **Quick Project Scaffolding**: Generate complete Go microservice projects in seconds
- **Production-Ready Templates**: Docker, hot reload, migrations, Swagger and observability
- **MySQL or PostgreSQL**: `--db-driver` switches the driver, DSN, migration backend and dependency
- **Git Integration**: Automatic Git repository initialization with dev branch
- **Go Module Management**: Automatic `go mod init` and `go mod tidy` execution
- **Hot Reload Support**: Includes Air configuration for development
- **API Documentation**: Swagger/OpenAPI integration
- **Observability**: OpenTelemetry tracing, JSON access logs and a Prometheus `/metrics` endpoint
- **Hardened by default**: real client-IP resolution behind proxies, plus an IP-keyed rate limiter
- **Database Migrations**: Built-in migration support with golang-migrate

## 📋 Prerequisites

- Go 1.23.0 or higher
- Git (for repository initialization)
- Docker (optional, for containerization)

## 🛠️ Installation

### Quick Install (Recommended)

#### Linux/macOS

```bash
# One-liner installation (recommended)
curl -fsSL https://raw.githubusercontent.com/surahj/gomicrogen/main/install-oneline.sh | bash

# Or download and run the full installer
curl -fsSL https://raw.githubusercontent.com/surahj/gomicrogen/main/install.sh | bash
```

#### Windows (PowerShell)

```powershell
# Run the PowerShell installer
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/surahj/gomicrogen/main/install.ps1" -UseBasicParsing).Content
```

### Manual Installation

Download the appropriate binary for your platform from the [latest release](https://github.com/surahj/gomicrogen/releases/latest):

#### Linux

```bash
# Download and extract the package
wget https://github.com/surahj/gomicrogen/releases/latest/download/gomicrogen-linux-amd64.tar.gz
tar -xzf gomicrogen-linux-amd64.tar.gz

# Install binary and templates
sudo cp gomicrogen-linux-amd64-package/gomicrogen-linux-amd64 /usr/local/bin/gomicrogen
sudo chmod +x /usr/local/bin/gomicrogen
sudo mkdir -p /usr/local/bin/templates
sudo cp -r gomicrogen-linux-amd64-package/templates/* /usr/local/bin/templates/
```

#### macOS

```bash
# Download and extract the package
curl -L https://github.com/surahj/gomicrogen/releases/latest/download/gomicrogen-darwin-amd64.tar.gz | tar -xz

# Install binary and templates
sudo cp gomicrogen-darwin-amd64-package/gomicrogen-darwin-amd64 /usr/local/bin/gomicrogen
sudo chmod +x /usr/local/bin/gomicrogen
sudo mkdir -p /usr/local/bin/templates
sudo cp -r gomicrogen-darwin-amd64-package/templates/* /usr/local/bin/templates/
```

#### Windows

```powershell
# Download and extract from the latest release
# Then add the extracted directory to your PATH
```

### From Source

```bash
# Clone the repository
git clone https://github.com/surahj/gomicrogen.git
cd gomicrogen

# Build the binary
go build -o gomicrogen

# Make it executable
chmod +x gomicrogen

# Move to a directory in your PATH (optional)
sudo mv gomicrogen /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/surahj/gomicrogen@latest
```

> `go install` places only the binary on your PATH. gomicrogen reads its templates from
> disk at generation time, so you also need a `templates/` directory next to the binary or
> in the working directory — use one of the installers above, or run from a clone.

### Using Docker

```bash
# Pull the latest image
docker pull ghcr.io/surahj/gomicrogen:latest

# Run gomicrogen
docker run --rm -it -v $(pwd):/workspace ghcr.io/surahj/gomicrogen:latest new my-service --module github.com/choplife-group/my-service
```

## 🎯 Quick Start

Create a new microservice with minimal configuration:

```bash
gomicrogen new my-service --module github.com/choplife-group/my-service
```

This will create a new microservice with:

- Complete project structure
- Docker configuration
- Database migrations
- API documentation
- Hot reload setup

Or scaffold for a specific domain:

```bash
gomicrogen new pawapay-service --module github.com/choplife-group/pawapay-service --type payment
```

## 📖 Usage

### Basic Command

```bash
gomicrogen new [service-name] [flags]
```

### Service Types

`--type` selects an overlay applied on top of the shared base. Run `gomicrogen types` to
list what is installed.

| type | adds | use for |
|---|---|---|
| `general` (default) | nothing beyond the base | anything that isn't one of the below |
| `casino` | gRPC server + `casino-service.proto` | game provider integrations, which reach wallet/identity/bonus through casino-service |
| `payment` | gRPC server + `wallet-service.proto`, RabbitMQ publisher, queue consumers | PSP integrations |

Only `casino` and `payment` open a gRPC port; a `general` service is HTTP only.

The `.proto` is shipped, not the generated `.pb.go` — run `protoc` yourself, then add the
client to the `Controller` and to `Initialize` in `app/router/router.go`:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       app/grpc/wallet/wallet-service.proto
```

### Required Flags

- `--module, -m`: Go module name (e.g., `github.com/choplife-group/service-name`)

### Optional Flags

Every flag below is covered by the test suite, which asserts the value reaches the file that
consumes it.

#### Service Configuration

| flag | default | lands in |
|---|---|---|
| `--type, -t` | `general` | which overlay is applied |
| `--description, -d` | `<name> microservice` | `docs.SwaggerInfo.Description` in `main.go` |
| `--version, -v` | `1.0.0` | `docs.SwaggerInfo.Version` in `main.go` |
| `--port, -p` | `8080` | `docker-compose-local.yml` |
| `--grpc-port, -g` | `8081` | `docker-compose-local.yml` and `GRPCRun` |
| `--env, -e` | `development` | uptrace deployment environment in `main.go` |

#### Database Configuration

| flag | default | lands in |
|---|---|---|
| `--db-driver` | `mysql` | driver import, DSN format, migration backend, `go.mod` |
| `--db-host` | `localhost` | `docker-compose-local.yml` |
| `--db-port` | follows the driver | `docker-compose-local.yml` |
| `--db-password` | | `docker-compose-local.yml` |

`--db-driver` accepts `mysql` or `postgres`; anything else is rejected. It sets the DSN
format, `otelsql` attributes, the golang-migrate backend, the `go.mod` dependency
(`go-sql-driver/mysql` or `lib/pq`), and exports `database.Driver` so handlers can pass it
to `goutils.Db{Dialect: database.Driver}` and get the right placeholder style.

`--db-port` defaults to `3306` for MySQL and `5432` for PostgreSQL; passing it explicitly
overrides that.

#### Redis Configuration

| flag | default | lands in |
|---|---|---|
| `--redis-host` | `localhost` | `REDIS_HOST` and `GLOBAL_REDIS_HOST` |
| `--redis-port` | `6379` | `REDIS_PORT` and `GLOBAL_REDIS_PORT` |
| `--redis-db-number` | `0` | `REDIS_DATABASE_NUMBER` and `GLOBAL_REDIS_DATABASE_NUMBER` |
| `--redis-password` | | `REDIS_PASSWORD` and `GLOBAL_REDIS_PASSWORD` |

Two Redis instances are configured: `REDIS_*` is the service-local cache, `GLOBAL_REDIS_*`
is the shared platform instance that `auth.Authenticate` reads tokens from. Both are seeded
from these flags; point them at different hosts in your deployment.

#### Output Options

- `--output-dir, -o`: Output directory (default: current directory)
- `--force`: Force overwrite if service already exists
- `--git`: Initialize Git repository with dev branch (default: true)
- `--go-mod`: Run go mod init and go mod tidy (default: true)

### Examples

#### Minimal Service

```bash
gomicrogen new user-service --module github.com/choplife-group/user-service
```

#### Service with Custom Configuration

```bash
gomicrogen new payment-service \
  --module github.com/choplife-group/payment-service \
  --type payment \
  --description "Payment processing microservice" \
  --version "2.1.0" \
  --port "3000" \
  --grpc-port "3001" \
  --db-driver "mysql" \
  --db-host "localhost" \
  --db-password "secret" \
  --redis-host "localhost" \
  --redis-port "6379" \
  --env "development"
```

#### PostgreSQL Service

```bash
gomicrogen new reporting-service \
  --module github.com/choplife-group/reporting-service \
  --db-driver postgres
```

#### Service in Custom Directory

```bash
gomicrogen new auth-service \
  --module github.com/choplife-group/auth-service \
  --output-dir /path/to/projects
```

#### Force Overwrite Existing Project

```bash
gomicrogen new my-service --module github.com/choplife-group/my-service --force
```

#### Skip Git and Go Module Initialization

```bash
gomicrogen new my-service --module github.com/choplife-group/my-service --git=false --go-mod=false
```

### Getting Help

```bash
# General help
gomicrogen --help

# Command-specific help
gomicrogen new --help

# List the installed service types
gomicrogen types

# Show all available commands
gomicrogen help
```

## 📁 Generated Project Structure

Every service gets this:

```
service-name/
├── app/
│   ├── auth/           # token validation middleware
│   ├── constants/      # application constants
│   ├── controllers/    # HTTP handlers and the Controller struct
│   ├── database/       # MySQL/Postgres, Redis and RabbitMQ connections
│   ├── library/        # shared helpers
│   ├── models/         # request/response structs
│   ├── router/         # router.go (App, Initialize, setRouters, Run)
│   │                   # middleware.go (client IP, rate limiters), status.go
│   └── utils/          # utility functions
├── docs/               # swagger, regenerated by `swag init`
├── migrations/         # golang-migrate .up.sql, applied at startup
├── test/
├── .gitignore
├── air.toml            # hot reload configuration
├── docker-compose-local.yml
├── Dockerfile          # production image
├── Dockerfile.dev      # development image
├── go.mod
├── main.go             # OTel setup → migrations → Initialize → Run
└── Makefile
```

`--type casino` additionally brings:

```
├── app/grpc/casino/casino-service.proto
└── app/router/router.go   # + GRPCRun and getGrpcConn
```

`--type payment` additionally brings:

```
├── app/grpc/wallet/wallet-service.proto
├── app/publisher/         # RabbitMQ publisher
├── app/queue/             # consumers, driven by the QUEUES env var
├── app/rabbitmq/          # connection handling
└── app/router/router.go   # + GRPCRun, getGrpcConn, publisher and queue wiring
```

Routes go in `setRouters` in `app/router/router.go`, and collaborators are wired inline in
`Initialize` — the same shape every service in the fleet uses.

## 🔧 Configuration

### Environment Variables

The generated service supports the following environment variables:

```bash
# Database Configuration
DATABASE_USERNAME=your_username
DATABASE_HOST=localhost
DATABASE_PORT=3306
DATABASE_PASSWORD=your_password
DATABASE_NAME=your_database
DATABASE_MAX_CONNECTION=150
DATABASE_IDLE_CONNECTION=100
DATABASE_CONNECTION_LIFETIME=60

# Redis Configuration — service-local cache
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE_NUMBER=0
REDIS_PASSWORD=your_redis_password

# Redis Configuration — shared platform instance that auth reads tokens from
GLOBAL_REDIS_HOST=localhost
GLOBAL_REDIS_PORT=6379
GLOBAL_REDIS_DATABASE_NUMBER=1
GLOBAL_REDIS_PASSWORD=your_redis_password

# Service Configuration
SYSTEM_HOST=0.0.0.0
SYSTEM_PORT=8080
SYSTEM_GRPC_PORT=8081        # casino and payment types only
ENV=development
SESSION_SECRET=change_me

# Rate limiting — all optional, defaults shown
RATE_LIMIT=20                        # requests per second, keyed on the resolved client IP
RATE_LIMIT_BURST=5
RATE_LIMIT_EXPIRES_IN_SECONDS=5
STRICT_RATE_LIMIT=5                  # for CustomStrictRateLimiterConfig, applied per route
STRICT_RATE_LIMIT_BURST=3
STRICT_RATE_LIMIT_EXPIRES_IN_SECONDS=60

# RabbitMQ — payment type
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASS=guest
RABBITMQ_VHOST=
QUEUES=                              # comma-separated queues to consume
QUEUE_PREFIX=

# Observability
UPTRACE_DSN=your_uptrace_dsn
BASE_URL=https://your-api-domain.com
```

`DATABASE_SSL_MODE` is also read when the service was generated with
`--db-driver postgres` (defaults to `disable`).

The `/metrics` path is exempt from the rate limiter, so a tight `RATE_LIMIT` can never turn
a Prometheus scrape into a 429.

### Docker Compose

For local development, the generated service includes a `docker-compose-local.yml` file:

```bash
# Start the service locally
docker-compose -f docker-compose-local.yml up

# Build and start
docker-compose -f docker-compose-local.yml up --build
```

## 🚀 Getting Started with Generated Service

1. **Navigate to the project directory**:

   ```bash
   cd your-service-name
   ```

2. **Install dependencies**:

   ```bash
   go mod tidy
   ```

3. **Set up environment variables**:

   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

4. **Start the service**:

   ```bash
   # Development with hot reload
   make run-air

   # Or directly
   go run .
   ```

   Migrations run automatically at startup — `main.go` applies every `.up.sql` in
   `migrations/` before the server starts. There is no separate migrate step.

## 🛠️ Development Commands

The generated service ships this Makefile:

```bash
make install      # go mod vendor && go mod download
make build        # swag init, then build with version ldflags
make run          # run the built binary
make run-air      # hot reload via air
make swagger      # regenerate docs/ from the @ annotations
make docker-up    # docker compose -f docker-compose-local.yml up -d
make docker-down
```

Re-run `make swagger` (or `swag init`) after changing the Swagger annotations in
`app/router/`; `make build` and the Dockerfile both do it for you.

## 📚 API Documentation

The generated service includes Swagger/OpenAPI documentation:

1. **Start the service**
2. **Visit**: `http://localhost:8080/docs/index.html`

## 🔍 Observability

Every generated service exposes:

| endpoint | purpose |
|---|---|
| `GET \| POST /` | status — pings the database and Redis |
| `GET /docs/*` | Swagger UI |
| `GET /metrics` | Prometheus, with the fleet-wide `echo_requests_total` and `echo_request_duration_seconds` |

Wired in via `go-utils/observability`: JSON access logs to stdout (one line per request with
method, status, latency, resolved client IP and request id), OpenTelemetry tracing through
Uptrace, and gzip that skips `/metrics` so scrapes aren't double-compressed.

The client IP is resolved from the proxy headers before anything else runs, preferring the
first *public* address in `X-Original-Client-Ip`, `X-Client-Ip`, `Cf-Connecting-Ip`,
`True-Client-Ip`, then the leftmost public entry of `X-Forwarded-For`, then the L3 peer.
Private and loopback addresses are ignored, so a spoofed header cannot mask the caller and
the rate limiter keys on the real client rather than the proxy.

## 🐳 Docker Support

### Development

```bash
# Build development image
docker build -f Dockerfile.dev -t my-service:dev .

# Run with docker-compose
docker-compose -f docker-compose-local.yml up
```

### Production

```bash
# Build production image
docker build -t my-service:latest .

# Run production container
docker run -p 8080:80 -p 8081:81 my-service:latest
```

The image listens on port 80 for HTTP and 81 for gRPC; `docker-compose-local.yml` maps your
`--port` and `--grpc-port` onto those.

## 🧪 Tests

```bash
make test        # fast suite — no docker, no network
make test-e2e    # end-to-end — requires docker
make test-all
```

The fast suite covers the generator and every CLI flag by running the real binary against
the real templates. The end-to-end suite spins up MySQL, Postgres, Redis and RabbitMQ with
testcontainers, then generates, compiles, runs and drives each service type over HTTP.

See [test/README.md](test/README.md) for what each tier asserts.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Troubleshooting

### Common Issues

1. **Installation fails**

   - Ensure you have write permissions to `/usr/local/bin/` or use the one-line installer
   - Check that curl/wget is available on your system
   - Verify your internet connection can access GitHub

2. **Go module initialization fails**

   - Check that Go is properly installed and in your PATH
   - Ensure you have write permissions in the target directory
   - Verify the module name follows Go module naming conventions

3. **Database connection issues**

   - Verify database credentials and connection settings
   - Ensure the database server is running and accessible
   - `--db-driver` accepts `mysql` and `postgres` only

4. **`Unknown service type` on generation**

   - Run `gomicrogen types` to see what is installed
   - If it reports legacy templates, reinstall so the binary and templates match

5. **Port already in use**

   - Use different ports with `--port` and `--grpc-port` flags
   - Check for existing services using the same ports
   - Use `netstat -tulpn | grep :8080` to check port usage

6. **Git initialization fails**
   - Ensure Git is installed and configured
   - Check that you have write permissions in the target directory
   - Use `--git=false` to skip Git initialization if needed

### Getting Help

- Create an issue on [GitHub](https://github.com/surahj/gomicrogen/issues)
- Check the existing issues for similar problems
- Review the generated code and configuration files
- Use `gomicrogen --help` for command-line help

## 🔗 Related Projects

- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [Air](https://github.com/cosmtrek/air) - Live reload for Go apps
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [Uptrace](https://uptrace.dev/) - Observability platform
