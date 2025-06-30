# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added

- Initial release of gomicrogen CLI tool
- Cross-platform binary builds (Linux, macOS, Windows)
- Automated installation scripts
- GitHub Actions CI/CD workflows
- Docker support
- Comprehensive project scaffolding
- Database migration support
- Kubernetes manifests generation
- API documentation with Swagger
- Hot reload configuration with Air
- Observability integration (OpenTelemetry, Uptrace)

### Changed

### Deprecated

### Removed

### Fixed

### Security

## [1.0.0] - 2024-01-XX

### Added

- Initial release of gomicrogen CLI tool
- Cross-platform binary builds (Linux, macOS, Windows)
- Automated installation scripts for Linux/macOS and Windows
- GitHub Actions workflows for CI/CD and releases
- Docker containerization support
- Comprehensive Go microservice project scaffolding
- Database migration support with golang-migrate
- Kubernetes deployment manifests
- API documentation with Swagger/OpenAPI
- Hot reload development setup with Air
- Observability integration (OpenTelemetry, Uptrace)
- Redis configuration support
- Environment-specific configurations
- Git repository initialization
- Go module management
- Production-ready Docker configurations
- Local development with Docker Compose
- Comprehensive Makefile with build targets
- Cross-platform release artifacts (tar.gz, zip)

### Features

- **Project Generation**: Create complete microservice projects with `gomicrogen new`
- **Database Support**: MySQL, PostgreSQL with migration support
- **Redis Integration**: Caching and session management
- **API Documentation**: Automatic Swagger/OpenAPI generation
- **Containerization**: Docker and Docker Compose support
- **Kubernetes**: Ready-to-deploy manifests
- **Development Tools**: Hot reload, linting, testing setup
- **Observability**: Distributed tracing and metrics
- **Multi-platform**: Support for Linux, macOS, and Windows
- **Easy Installation**: One-liner installation scripts

### Installation

#### Linux/macOS

```bash
curl -fsSL https://raw.githubusercontent.com/surahj/gomicrogen/main/install-oneline.sh | bash
```

#### Windows (PowerShell)

```powershell
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/surahj/gomicrogen/main/install.ps1" -UseBasicParsing).Content
```

### Quick Start

```bash
gomicrogen new my-service --module github.com/your-org/my-service
```

---

## Version History

- **1.0.0**: Initial release with full microservice scaffolding capabilities

## Release Notes

### Version 1.0.0

- First stable release of gomicrogen
- Complete microservice project generation
- Cross-platform support
- Automated installation and deployment
- Production-ready templates and configurations

## Contributing

When contributing to this project, please update this changelog with:

1. A new entry under `[Unreleased]` for your changes
2. Move `[Unreleased]` to a new version section when releasing
3. Follow the [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format

## Links

- [GitHub Repository](https://github.com/surahj/gomicrogen)
- [Documentation](https://github.com/surahj/gomicrogen#readme)
- [Issues](https://github.com/surahj/gomicrogen/issues)
- [Releases](https://github.com/surahj/gomicrogen/releases)
