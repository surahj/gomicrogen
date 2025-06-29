# Deployment Guide

This guide explains how to deploy and release the gomicrogen CLI tool.

## 🚀 Release Process

### 1. Prepare for Release

Before creating a release, ensure:

- All tests pass: `make test`
- Code is linted: `make lint`
- Code is formatted: `make fmt`
- Version is updated in code (if needed)

### 2. Create a Release

#### Automated Release (Recommended)

1. **Create and push a tag:**

   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically:**
   - Build binaries for all platforms
   - Create release archives
   - Publish to GitHub Releases
   - Update installation scripts

#### Manual Release

1. **Build all platforms:**

   ```bash
   make release
   ```

2. **Create GitHub Release:**
   - Go to GitHub repository
   - Click "Releases" → "Create a new release"
   - Tag: `v1.0.0`
   - Title: `Release v1.0.0`
   - Upload files from `release/` directory

### 3. Update Installation Scripts

The installation scripts automatically fetch the latest version from GitHub Releases. No manual updates needed.

## 📦 Distribution Methods

### 1. GitHub Releases

- **Automatic**: GitHub Actions creates releases with all platform binaries
- **Manual**: Upload built artifacts to GitHub Releases

### 2. Package Managers

#### Homebrew (macOS)

```bash
# Create formula
brew create --go https://github.com/Choplife-group/gomicrogen

# Install
brew install gomicrogen
```

#### Snap (Linux)

```bash
# Create snap package
snapcraft init

# Build and publish
snapcraft build
snapcraft upload gomicrogen_*.snap
```

#### Chocolatey (Windows)

```powershell
# Create package
choco new gomicrogen

# Install
choco install gomicrogen
```

### 3. Docker Hub

```bash
# Build and push Docker image
docker build -t gomicrogen:latest .
docker tag gomicrogen:latest ghcr.io/choplife-group/gomicrogen:latest
docker push ghcr.io/choplife-group/gomicrogen:latest
```

## 🔧 CI/CD Setup

### GitHub Actions

The repository includes:

- **CI Workflow** (`.github/workflows/ci.yml`):

  - Runs on PRs and pushes to main
  - Tests, lints, and builds
  - Uploads build artifacts

- **Release Workflow** (`.github/workflows/release.yml`):
  - Triggers on tag push
  - Builds cross-platform binaries
  - Creates GitHub release
  - Uploads release artifacts

### Local Development

```bash
# Development
make dev

# Testing
make test
make test-coverage

# Building
make build
make build-all

# Release preparation
make release
```

## 📋 Release Checklist

- [ ] All tests pass
- [ ] Code is linted and formatted
- [ ] Documentation is updated
- [ ] Version is updated (if needed)
- [ ] Changelog is updated
- [ ] Tag is created and pushed
- [ ] GitHub Actions complete successfully
- [ ] Release notes are reviewed
- [ ] Installation scripts work correctly
- [ ] Cross-platform binaries are tested

## 🐛 Troubleshooting

### Build Issues

```bash
# Clean and rebuild
make clean
make build

# Check Go version
go version

# Update dependencies
go mod tidy
```

### Release Issues

```bash
# Check GitHub Actions logs
# Verify tag format (v*)
# Ensure repository permissions
```

### Installation Issues

```bash
# Test installation script locally
./install-oneline.sh

# Check binary permissions
chmod +x gomicrogen

# Verify PATH
echo $PATH
which gomicrogen
```

## 📊 Monitoring

### Release Metrics

- Download counts per platform
- Installation success rate
- User feedback and issues
- Version adoption rate

### Health Checks

```bash
# Test binary functionality
gomicrogen --help
gomicrogen version

# Test project generation
gomicrogen new test-service --module github.com/test/service
```

## 🔄 Update Process

### For Users

1. **Automatic updates**: Installation scripts fetch latest version
2. **Manual updates**: Download new binary from releases
3. **Package managers**: Use standard update commands

### For Maintainers

1. **Patch releases**: Bug fixes and minor improvements
2. **Minor releases**: New features, backward compatible
3. **Major releases**: Breaking changes

## 📞 Support

- **Issues**: GitHub Issues
- **Discussions**: GitHub Discussions
- **Documentation**: README.md and this guide
- **Examples**: Generated project templates
