# gomicrogen Build and Release Setup

Maintainer-facing summary of how gomicrogen is built, tested, packaged and released.
For usage, see [README.md](README.md).

## The one invariant

**gomicrogen reads its templates from disk at generation time.** The binary on its own
cannot generate anything. Everything below exists to keep the binary and its `templates/`
directory together:

- Release archives contain a `gomicrogen-<os>-<arch>-package/` directory holding both.
- The installers look for exactly that directory and copy `templates/` next to the binary.
- CI and the release workflow both verify the layout before anything is published.

An archive without it produces `❌ Failed to extract package from archive` at install time.

## Build system

`Makefile` is the single packaging path — CI calls it rather than reimplementing packaging.

```bash
make build         # current platform
make build-all     # all six platform/arch combinations into dist/
make release       # dist/ -> release/*.tar.gz and release/*.zip, with templates
make test          # fast suite: generator + every CLI flag, no docker
make test-e2e      # end-to-end: real MySQL/Postgres/Redis/RabbitMQ via testcontainers
make test-all      # both
make clean
```

Version information is stamped in through ldflags:

```
-X github.com/Choplife-group/gomicrogen/cmd.appVersion=${VERSION}
-X github.com/Choplife-group/gomicrogen/cmd.appCommit=${COMMIT}
-X github.com/Choplife-group/gomicrogen/cmd.appDate=${BUILD_TIME}
```

`VERSION` comes from `git describe --tags`, and the release workflow overrides it with the
tag being built (`make release VERSION=v1.2.3`). The variables live in the `cmd` package —
stamping `main.version` silently does nothing.

macOS note: the tar steps set `COPYFILE_DISABLE=1`. Without it, macOS stores extended
attributes that GNU tar materialises as `._*` files on extraction, which then get copied
into every generated service.

## Installation scripts

| script | platform |
|---|---|
| `install-oneline.sh` | Linux/macOS, the documented one-liner |
| `install.sh` | Linux/macOS, verbose installer with flags |
| `install.ps1` | Windows PowerShell |

All install the binary **and** `templates/`, replacing any existing templates directory
rather than merging into it, so an upgrade cannot leave stale files from an older layout.

`GOMICROGEN_BASE_URL` overrides where archives are fetched from — used to test an installer
against a local `make release` build, or to serve from an internal mirror:

```bash
GOMICROGEN_BASE_URL="file://$PWD/release" bash install-oneline.sh
```

Verified on Ubuntu 22.04, Debian 12 and Alpine 3.19. Alpine has no `bash`, so the one-liner
must be piped to `sh` there.

## CI/CD

### `.github/workflows/ci.yml` — on PRs and pushes to main

- **test** — `go test -race ./...`
- **package** — runs `make release`, verifies every archive has its `-package/` directory,
  `templates/base`, `templates/types` and no AppleDouble files, then installs from the built
  archive and generates one service of each type

The package job exists because packaging bugs used to surface only after a release was
published. It now fails the PR instead.

### `.github/workflows/release.yml` — on tag push, and callable

Runs the tests, builds every platform with `make release`, verifies the archive layout, and
publishes the GitHub release. It is also a reusable workflow (`workflow_call`) so the manual
tag workflow shares the same packaging steps.

### `.github/workflows/tag.yml` — manual

Tagging is deliberate: merging to main never publishes.

```bash
gh workflow run tag.yml -f version=v1.1.0                          # tag and release
gh workflow run tag.yml -f version=v1.1.0 -f create_release=false  # tag only
gh workflow run tag.yml -f version=v1.1.0 -f create_tag=false      # release an existing tag
gh workflow run tag.yml -f bump=patch                              # v1.0.2 -> v1.0.3
```

It calls `release.yml` directly rather than relying on the tag push to trigger it — a tag
pushed with `GITHUB_TOKEN` does not fire other workflows.

Tag examples: `v1.0.3` patch (fixes), `v1.1.0` minor (a new flag or service type),
`v2.0.0` major (a removed flag or a template layout change).

## Release artifacts

```
gomicrogen-linux-amd64.tar.gz     gomicrogen-darwin-amd64.tar.gz
gomicrogen-linux-arm64.tar.gz     gomicrogen-darwin-arm64.tar.gz
gomicrogen-windows-amd64.zip      gomicrogen-windows-arm64.zip
```

Each expands to:

```
gomicrogen-<os>-<arch>-package/
├── gomicrogen-<os>-<arch>        # the binary
└── templates/
    ├── base/                     # shared by every service type
    └── types/{general,casino,payment}/
```

## Releasing by hand

If the workflows are unavailable:

```bash
make test
make release            # archives land in release/
# upload release/* to a GitHub release
```

Before uploading, confirm an archive extracts to a `-package/` directory containing
`templates/base` and `templates/types`.

## Not covered here

- Package managers (Homebrew, Snap, Chocolatey)
- Release signing and checksums
- Docker Hub publishing
- Any auto-update mechanism
