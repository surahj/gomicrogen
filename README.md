# gomicrogen

A powerful CLI tool to scaffold Go microservice projects with a predefined folder structure and production-ready templates.

## 🚀 Features

- **Quick Project Scaffolding**: Generate complete Go microservice projects in seconds
- **Production-Ready Templates**: Includes Docker, Kubernetes, database migrations, and more
- **Customizable Configuration**: Flexible flags for database, Redis, ports, and environment settings
- **Git Integration**: Automatic Git repository initialization with dev branch
- **Go Module Management**: Automatic `go mod init` and `go mod tidy` execution
- **Comprehensive Structure**: Well-organized folder structure following Go best practices
- **Hot Reload Support**: Includes Air configuration for development
- **API Documentation**: Swagger/OpenAPI integration
- **Observability**: OpenTelemetry and Uptrace integration
- **Database Migrations**: Built-in migration support with golang-migrate

## 📋 Prerequisites

- Go 1.23.0 or higher
- Git (for repository initialization)
- Docker (optional, for containerization)

## ��️ Installation

### Quick Install (Recommended)

#### Linux/macOS

```bash
# One-liner installation
curl -fsSL https://raw.githubusercontent.com/Choplife-group/gomicrogen/main/install-oneline.sh | bash

# Or download and run the full installer
curl -fsSL https://raw.githubusercontent.com/Choplife-group/gomicrogen/main/install.sh | bash
```

#### Windows (PowerShell)

```powershell
# Run the PowerShell installer
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/Choplife-group/gomicrogen/main/install.ps1" -UseBasicParsing).Content
```

### Manual Installation

Download the appropriate binary for your platform from the [latest release](https://github.com/Choplife-group/gomicrogen/releases/latest):

#### Linux

```bash
# AMD64
wget https://github.com/Choplife-group/gomicrogen/releases/latest/download/gomicrogen-linux-amd64.tar.gz
tar -xzf gomicrogen-linux-amd64.tar.gz
sudo mv gomicrogen-linux-amd64 /usr/local/bin/gomicrogen
sudo chmod +x /usr/local/bin/gomicrogen

# ARM64
wget https://github.com/Choplife-group/gomicrogen/releases/latest/download/gomicrogen-linux-arm64.tar.gz
tar -xzf gomicrogen-linux-arm64.tar.gz
sudo mv gomicrogen-linux-arm64 /usr/local/bin/gomicrogen
sudo chmod +x /usr/local/bin/gomicrogen
```

#### macOS

```bash
# AMD64
curl -L https://github.com/Choplife-group/gomicrogen/releases/latest/download/gomicrogen-darwin-amd64.tar.gz | tar -xz
sudo mv gomicrogen-darwin-amd64 /usr/local/bin/gomicrogen
sudo chmod +x /usr/local/bin/gomicrogen

# ARM64 (Apple Silicon)
curl -L https://github.com/Choplife-group/gomicrogen/releases/latest/download/gomicrogen-darwin-arm64.tar.gz | tar -xz
sudo mv gomicrogen-darwin-arm64 /usr/local/bin/gomicrogen
sudo chmod +x /usr/local/bin/gomicrogen
```

#### Windows

```powershell
# Download and extract from the latest release
# Then add the extracted directory to your PATH
```

### From Source

```bash
# Clone the repository
git clone https://github.com/Choplife-group/gomicrogen.git
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
go install github.com/Choplife-group/gomicrogen@latest
```

### Using Docker

```bash
# Pull the latest image
docker pull ghcr.io/choplife-group/gomicrogen:latest

# Run gomicrogen
docker run --rm -it -v $(pwd):/workspace ghcr.io/choplife-group/gomicrogen:latest new my-service --module github.com/your-org/my-service
```

## 🎯 Quick Start

Create a new microservice with minimal configuration:

```bash
gomicrogen new my-service --module github.com/your-org/my-service
```

This will create a new microservice with:

- Complete project structure
- Docker configuration
- Kubernetes manifests
- Database migrations
- API documentation
- Hot reload setup

## 📖 Usage

### Basic Command

```bash
gomicrogen new [service-name] [flags]
```

### Required Flags

- `--module, -m`: Go module name (e.g., `github.com/your-org/service-name`)

### Optional Flags

#### Service Configuration

- `--description, -d`: Service description
- `--version, -v`: Service version (default: "1.0.0")
- `--author, -a`: Author name
- `--port, -p`: HTTP port (default: "8080")
- `--grpc-port, -g`: gRPC port (default: "8081")
- `--env, -e`: Environment (development, staging, production)

#### Database Configuration

- `--db-driver`: Database driver (mysql, postgres, etc.)
- `--db-url`: Database connection URL
- `--db-host`: Database host
- `--db-port`: Database port
- `--db-password`: Database password

#### Redis Configuration

- `--redis-url`: Redis connection URL
- `--redis-host`: Redis host
- `--redis-port`: Redis port
- `--redis-db-number`: Redis database number
- `--redis-password`: Redis password

#### Output Options

- `--output-dir, -o`: Output directory (default: current directory)
- `--force`: Force overwrite if service already exists
- `--git`: Initialize Git repository with dev branch (default: true)
- `--go-mod`: Run go mod init and go mod tidy (default: true)

### Examples

#### Minimal Service

```bash
gomicrogen new user-service --module github.com/mycompany/user-service
```

#### Service with Custom Configuration

```bash
gomicrogen new payment-service \
  --module github.com/mycompany/payment-service \
  --description "Payment processing microservice" \
  --version "2.1.0" \
  --author "John Doe" \
  --port "3000" \
  --grpc-port "3001" \
  --db-driver "mysql" \
  --db-host "localhost" \
  --db-port "3306" \
  --db-password "secret" \
  --redis-host "localhost" \
  --redis-port "6379" \
  --env "development"
```

#### Service in Custom Directory

```bash
gomicrogen new auth-service \
  --module github.com/mycompany/auth-service \
  --output-dir /path/to/projects
```

## 📁 Generated Project Structure

```
service-name/
├── app/
│   ├── auth/           # Authentication middleware
│   ├── constants/      # Application constants
│   ├── controllers/    # HTTP controllers
│   ├── database/       # Database connection and models
│   ├── grpc/          # gRPC server and services
│   ├── library/       # Shared libraries and utilities
│   ├── models/        # Data models
│   ├── router/        # HTTP router setup
│   └── utils/         # Utility functions
├── cmd/               # Command-line tools
├── docs/              # API documentation
├── k8s/               # Kubernetes manifests
├── migrations/        # Database migrations
├── test/              # Test files
├── .gitignore         # Git ignore file
├── air.toml           # Hot reload configuration
├── docker-compose-local.yml  # Local development setup
├── Dockerfile         # Production Docker image
├── Dockerfile.dev     # Development Docker image
├── go.mod             # Go module file
├── go.sum             # Go dependencies checksum
├── main.go            # Application entry point
└── Makefile           # Build and deployment commands
```

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

# Redis Configuration
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_DATABASE_NUMBER=0
REDIS_PASSWORD=your_redis_password

# Service Configuration
SYSTEM_HOST=0.0.0.0
SYSTEM_PORT=8080
SYSTEM_GRPC_PORT=8081
ENV=development

# Observability
UPTRACE_DSN=your_uptrace_dsn
BASE_URL=https://your-api-domain.com
```

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

4. **Run database migrations**:

   ```bash
   make migrate-up
   ```

5. **Start the service**:

   ```bash
   # Development with hot reload
   make dev

   # Production
   make run

   # Or directly
   go run main.go
   ```

## 🛠️ Development Commands

The generated service includes a comprehensive Makefile with useful commands:

```bash
# Development
make dev          # Run with hot reload using Air
make run          # Run the service
make build        # Build the binary
make test         # Run tests
make test-coverage # Run tests with coverage

# Database
make migrate-up   # Run database migrations
make migrate-down # Rollback database migrations
make migrate-create # Create new migration

# Docker
make docker-build # Build Docker image
make docker-run   # Run Docker container
make docker-stop  # Stop Docker container

# Kubernetes
make k8s-deploy   # Deploy to Kubernetes
make k8s-delete   # Delete from Kubernetes

# Utilities
make clean        # Clean build artifacts
make lint         # Run linter
make fmt          # Format code
```

## 📚 API Documentation

The generated service includes Swagger/OpenAPI documentation:

1. **Start the service**
2. **Visit**: `http://localhost:8080/swagger/index.html`

## 🔍 Observability

The generated service includes:

- **OpenTelemetry**: Distributed tracing
- **Uptrace**: Metrics and monitoring
- **Structured Logging**: Using logrus

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
docker run -p 8080:8080 -p 8081:8081 my-service:latest
```

## ☸️ Kubernetes Deployment

The generated service includes Kubernetes manifests in the `k8s/` directory:

```bash
# Deploy to Kubernetes
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -l app=my-service
```

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

1. **Template directory not found**

   - Ensure the `templates/` directory is in the same directory as the binary
   - Or use the `--output-dir` flag to specify a different location

2. **Go module initialization fails**

   - Check that Go is properly installed and in your PATH
   - Ensure you have write permissions in the target directory

3. **Database connection issues**

   - Verify database credentials and connection settings
   - Ensure the database server is running and accessible

4. **Port already in use**
   - Use different ports with `--port` and `--grpc-port` flags
   - Check for existing services using the same ports

### Getting Help

- Create an issue on GitHub
- Check the existing issues for similar problems
- Review the generated code and configuration files

## 🔗 Related Projects

- [Cobra](https://github.com/spf13/cobra) - CLI framework for Go
- [Air](https://github.com/cosmtrek/air) - Live reload for Go apps
- [golang-migrate](https://github.com/golang-migrate/migrate) - Database migrations
- [Uptrace](https://uptrace.dev/) - Observability platform
