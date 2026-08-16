# Tests

Two tiers, split by what they cost to run.

## Fast suite — `make test`

Runs in seconds. No docker, no network, no module downloads: every generation
passes `--go-mod=false --git=false`.

| package | covers |
|---|---|
| `internal/config` | defaults, driver validation, driver-specific ports |
| `internal/generator` | layout resolution, type resolution and aliases, file filters, overlay replacement, template substitution |
| `cmd` | the real binary against the real templates — every flag, every type, every error path |

`cmd/cli_test.go` builds the CLI once in `TestMain` and executes it, so it exercises
cobra wiring, flag parsing and the generator together rather than calling internals.

What the fast suite locks down, in the order things have actually broken:

- a mistyped `--type` must not reach the `--force` removal and delete a service
- `--type none` / `base` / `""` must resolve to the general overlay, not a base-only
  tree that omits `router.go` and cannot compile
- `.pb.go` must be emitted (overlays ship gRPC clients) but never templated
- `.DS_Store`, `go.sum`, `base/`, `types/` and `type.json` must never reach a service
- a `tmp` ancestor in the templates path must not silently skip every file
- casino symbols must not leak into general or payment services
- each flag's value must reach the file that consumes it

## End-to-end suite — `make test-e2e`

Requires docker. Spins up MySQL, Postgres, Redis and RabbitMQ with testcontainers,
generates each service type, compiles it, runs it, and drives it over HTTP.

```bash
make test-e2e
# or
cd test/e2e && GOMICROGEN_E2E=1 go test -v -timeout 45m ./...
```

Without `GOMICROGEN_E2E=1` the suite prints a skip notice and exits 0, so it stays
out of the way in environments with no docker.

It asserts that a generated service:

- serves `/` with a real database ping and Redis `PONG`
- serves `/docs/index.html` and `/metrics` with the fleet-wide `echo_requests_total`
- applies its migrations (`schema_migrations` present and not dirty)
- opens a gRPC port for casino and payment, and **not** for general
- rate limits ordinary routes while never limiting `/metrics`
- resolves the client IP from `X-Forwarded-For`, preferring a public entry and
  ignoring a private one
- works on both `--db-driver mysql` and `--db-driver postgres`
- still finds `migrations/` when the binary is moved away from where it was built

### Why a separate module

`test/e2e` has its own `go.mod`. testcontainers pulls a large dependency tree, and
the CLI itself depends only on cobra — anyone running `go install` on gomicrogen
should not inherit those. The nested module also keeps `go test ./...` in the root
fast and docker-free.

`templates/` likewise carries a `go.mod` so the toolchain stops descending into the
template tree, which is not compilable Go on its own.
