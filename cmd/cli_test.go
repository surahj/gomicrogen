package cmd_test

// End-to-end tests for the CLI: they build the real binary and run it against
// the real templates, then assert on the generated tree.
//
// These are hermetic — no network, no docker — because every generation passes
// --go-mod=false and --git=false. Compiling and running the generated services
// lives in the separate test/e2e module.

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var (
	binary   string
	repoRoot string
)

func TestMain(m *testing.M) {

	root, err := filepath.Abs("..")
	if err != nil {
		panic(err)
	}
	repoRoot = root

	dir, err := os.MkdirTemp("", "gomicrogen-bin")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)

	binary = filepath.Join(dir, "gomicrogen")

	build := exec.Command("go", "build", "-o", binary, ".")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build gomicrogen: " + string(out))
	}

	os.Exit(m.Run())
}

// generate runs the CLI and returns the target directory plus combined output.
func generate(t *testing.T, name string, args ...string) (string, string, error) {
	t.Helper()

	out := t.TempDir()

	full := append([]string{
		"new", name,
		"--module", "github.com/test-org/" + name,
		"--output-dir", out,
		"--git=false",
		"--go-mod=false",
	}, args...)

	cmd := exec.Command(binary, full...)
	cmd.Dir = repoRoot // so findTemplatesDir picks up ./templates

	combined, err := cmd.CombinedOutput()

	return filepath.Join(out, name), string(combined), err
}

func mustGenerate(t *testing.T, name string, args ...string) string {
	t.Helper()

	dir, out, err := generate(t, name, args...)
	if err != nil {
		t.Fatalf("generation failed: %v\n%s", err, out)
	}

	return dir
}

func exists(t *testing.T, dir, rel string) bool {
	t.Helper()

	_, err := os.Stat(filepath.Join(dir, rel))

	return err == nil
}

// fileContains reads a generated file directly. Recursive greps are unreliable
// here because the generated .gitignore hides docker-compose-local.yml.
func fileContains(t *testing.T, dir, rel, want string) bool {
	t.Helper()

	body, err := os.ReadFile(filepath.Join(dir, rel))
	if err != nil {
		t.Fatalf("read %s: %v", rel, err)
	}

	return strings.Contains(string(body), want)
}

func treeContains(t *testing.T, dir, want string) []string {
	t.Helper()

	var hits []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}

		body, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		if strings.Contains(string(body), want) {
			rel, _ := filepath.Rel(dir, path)
			hits = append(hits, rel)
		}

		return nil
	})
	if err != nil {
		t.Fatalf("walk: %v", err)
	}

	return hits
}

// --- every service type ------------------------------------------------------

func TestEachTypeGeneratesItsOwnShape(t *testing.T) {

	cases := []struct {
		serviceType string
		wantFiles   []string
		wantAbsent  []string
	}{
		{
			serviceType: "general",
			wantFiles: []string{
				"main.go", "go.mod", "Dockerfile", "docker-compose-local.yml",
				"app/router/router.go", "app/router/middleware.go", "app/router/status.go",
				"app/controllers/controller.go", "app/database/database.go",
			},
			// general is HTTP only and has no upstream
			wantAbsent: []string{"app/grpc", "app/queue", "app/publisher", "app/rabbitmq"},
		},
		{
			serviceType: "casino",
			wantFiles: []string{
				"app/router/router.go",
				"app/grpc/casino/casino-service.proto",
			},
			wantAbsent: []string{
				"app/queue", "app/publisher", "app/rabbitmq",
				// generated protobuf is the developer's job, not the generator's
				"app/grpc/casino/casino-service.pb.go",
			},
		},
		{
			serviceType: "payment",
			wantFiles: []string{
				"app/router/router.go",
				"app/grpc/wallet/wallet-service.proto",
				"app/publisher/publisher.go", "app/queue/queue.go",
				"app/queue/consumer.go", "app/rabbitmq/rabbitmq.go",
			},
			wantAbsent: []string{
				"app/grpc/wallet/wallet-service.pb.go",
				"app/crontask",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.serviceType, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", tc.serviceType)

			for _, rel := range tc.wantFiles {
				if !exists(t, dir, rel) {
					t.Errorf("missing %s", rel)
				}
			}
			for _, rel := range tc.wantAbsent {
				if exists(t, dir, rel) {
					t.Errorf("%s must not be generated for type %s", rel, tc.serviceType)
				}
			}
		})
	}
}

// setRouters belongs in router.go, matching every service in the fleet.
func TestSetRoutersLivesInRouterGo(t *testing.T) {

	for _, serviceType := range []string{"general", "casino", "payment"} {
		t.Run(serviceType, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", serviceType)

			if !fileContains(t, dir, "app/router/router.go", "func (a *App) setRouters()") {
				t.Error("setRouters must be defined in app/router/router.go")
			}
			if exists(t, dir, "app/router/routes.go") || exists(t, dir, "app/router/wiring.go") {
				t.Error("routes.go/wiring.go must not exist; wiring is inline in router.go")
			}
		})
	}
}

// Only casino and payment serve gRPC.
func TestGRPCServerOnlyForTypedServices(t *testing.T) {

	cases := map[string]bool{"general": false, "casino": true, "payment": true}

	for serviceType, wantGRPC := range cases {
		t.Run(serviceType, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", serviceType)

			hasGRPCRun := fileContains(t, dir, "app/router/router.go", "func (a *App) GRPCRun()")
			if hasGRPCRun != wantGRPC {
				t.Errorf("GRPCRun present = %v, want %v", hasGRPCRun, wantGRPC)
			}
		})
	}
}

// The old template leaked casino structs into every service.
func TestNoCasinoLeakageIntoOtherTypes(t *testing.T) {

	for _, serviceType := range []string{"general", "payment"} {
		t.Run(serviceType, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", serviceType)

			for _, symbol := range []string{"CasinoGame", "LaunchURL", "TRANSACTION_TYPE_BET"} {
				if hits := treeContains(t, dir, symbol); len(hits) > 0 {
					t.Errorf("%s leaked into a %s service: %v", symbol, serviceType, hits)
				}
			}
		})
	}
}

func TestScaffoldingNeverLeaksIntoService(t *testing.T) {

	for _, serviceType := range []string{"general", "casino", "payment"} {
		t.Run(serviceType, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", serviceType)

			for _, leaked := range []string{"base", "types", "type.json", ".DS_Store", "go.sum"} {
				if exists(t, dir, leaked) {
					t.Errorf("%s must not appear in a generated service", leaked)
				}
			}
		})
	}
}

// --- type resolution ---------------------------------------------------------

func TestTypeAliasesProduceACompleteService(t *testing.T) {

	// A base-only generation would omit router.go and never compile.
	for _, alias := range []string{"general", "none", "base", "GENERAL"} {
		t.Run("alias="+alias, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", alias)

			if !exists(t, dir, "app/router/router.go") {
				t.Error("alias must resolve to the general overlay, which supplies router.go")
			}
		})
	}
}

func TestNoTypeFlagDefaultsToGeneral(t *testing.T) {

	dir := mustGenerate(t, "svc")

	if !exists(t, dir, "app/router/router.go") {
		t.Fatal("default generation is incomplete")
	}
	if exists(t, dir, "app/grpc") {
		t.Error("default type must be general, which has no gRPC")
	}
}

func TestUnknownTypeIsRejected(t *testing.T) {

	_, out, err := generate(t, "svc", "--type", "crypto")
	if err == nil {
		t.Fatal("an unknown type must fail")
	}

	for _, want := range []string{"crypto", "casino", "payment"} {
		if !strings.Contains(out, want) {
			t.Errorf("error output must mention %q, got:\n%s", want, out)
		}
	}
}

// A mistyped --type must never reach the --force removal.
func TestBadTypeWithForceDoesNotDeleteExistingService(t *testing.T) {

	out := t.TempDir()
	target := filepath.Join(out, "svc")

	if err := os.MkdirAll(target, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}

	canary := filepath.Join(target, "main.go")
	if err := os.WriteFile(canary, []byte("package main // precious"), 0o644); err != nil {
		t.Fatalf("write canary: %v", err)
	}

	cmd := exec.Command(binary, "new", "svc",
		"--module", "github.com/test-org/svc",
		"--output-dir", out,
		"--type", "casion", // typo
		"--force",
		"--git=false", "--go-mod=false")
	cmd.Dir = repoRoot

	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("a bad --type must fail, got success:\n%s", out)
	}

	body, err := os.ReadFile(canary)
	if err != nil {
		t.Fatal("the existing service was deleted before the type was validated")
	}
	if !strings.Contains(string(body), "precious") {
		t.Error("the existing service was overwritten")
	}
}

func TestForceOverwriteSwapsTypeCleanly(t *testing.T) {

	out := t.TempDir()

	run := func(serviceType string, extra ...string) {
		args := append([]string{"new", "svc",
			"--module", "github.com/test-org/svc",
			"--output-dir", out,
			"--type", serviceType,
			"--git=false", "--go-mod=false"}, extra...)

		cmd := exec.Command(binary, args...)
		cmd.Dir = repoRoot
		if o, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("generate %s: %v\n%s", serviceType, err, o)
		}
	}

	run("casino")
	run("payment", "--force")

	dir := filepath.Join(out, "svc")

	if exists(t, dir, "app/grpc/casino") {
		t.Error("--force must remove the previous type's files")
	}
	if !exists(t, dir, "app/grpc/wallet") {
		t.Error("--force must apply the new type")
	}
}

// --- flags -------------------------------------------------------------------

func TestFlagsReachGeneratedOutput(t *testing.T) {

	dir := mustGenerate(t, "flagsvc",
		"--type", "payment",
		"--description", "SENTINEL_DESCRIPTION",
		"--version", "9.9.9",
		"--port", "7777",
		"--grpc-port", "7778",
		"--env", "SENTINEL_ENV",
		"--db-host", "SENTINEL_DBHOST",
		"--db-port", "7306",
		"--db-password", "SENTINEL_DBPASS",
		"--redis-host", "SENTINEL_REDISHOST",
		"--redis-port", "7379",
		"--redis-db-number", "7",
		"--redis-password", "SENTINEL_REDISPASS",
	)

	cases := []struct {
		flag  string
		file  string
		value string
	}{
		{"--description", "main.go", "SENTINEL_DESCRIPTION"},
		{"--version", "main.go", "9.9.9"},
		{"--env", "main.go", "SENTINEL_ENV"},
		{"--port", "docker-compose-local.yml", "7777"},
		{"--grpc-port", "docker-compose-local.yml", "7778"},
		{"--grpc-port", "app/router/router.go", "7778"},
		{"--db-host", "docker-compose-local.yml", "SENTINEL_DBHOST"},
		{"--db-port", "docker-compose-local.yml", "7306"},
		{"--db-password", "docker-compose-local.yml", "SENTINEL_DBPASS"},
		{"--redis-host", "docker-compose-local.yml", "SENTINEL_REDISHOST"},
		{"--redis-port", "docker-compose-local.yml", "7379"},
		{"--redis-db-number", "docker-compose-local.yml", "7"},
		{"--redis-password", "docker-compose-local.yml", "SENTINEL_REDISPASS"},
	}

	for _, tc := range cases {
		t.Run(tc.flag+" -> "+tc.file, func(t *testing.T) {
			if !fileContains(t, dir, tc.file, tc.value) {
				t.Errorf("%s value %q never reached %s", tc.flag, tc.value, tc.file)
			}
		})
	}
}

// Redis must be configured for both the local cache and the shared platform
// instance that auth reads from.
func TestRedisFlagsPopulateBothRedisBlocks(t *testing.T) {

	dir := mustGenerate(t, "svc", "--redis-host", "redis.internal", "--redis-port", "6380")

	for _, key := range []string{
		"REDIS_HOST: redis.internal",
		"GLOBAL_REDIS_HOST: redis.internal",
		"REDIS_PORT: 6380",
		"GLOBAL_REDIS_PORT: 6380",
	} {
		if !fileContains(t, dir, "docker-compose-local.yml", key) {
			t.Errorf("compose is missing %q", key)
		}
	}
}

func TestModuleFlagRewritesImports(t *testing.T) {

	dir := mustGenerate(t, "svc", "--type", "payment")

	if !fileContains(t, dir, "go.mod", "module github.com/test-org/svc") {
		t.Error("go.mod does not carry the requested module path")
	}
	if !fileContains(t, dir, "app/router/router.go", "github.com/test-org/svc/app/controllers") {
		t.Error("internal imports were not rewritten to the module path")
	}
}

func TestOutputDirIsHonoured(t *testing.T) {

	// relative output dirs once broke go mod init and .gitignore creation
	rel := filepath.Join("files", "testtmp")

	full := filepath.Join(repoRoot, rel)
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer os.RemoveAll(full)

	cmd := exec.Command(binary, "new", "relsvc",
		"--module", "github.com/test-org/relsvc",
		"--output-dir", rel,
		"--git=false", "--go-mod=false")
	cmd.Dir = repoRoot

	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("relative output-dir failed: %v\n%s", err, out)
	}

	if _, err := os.Stat(filepath.Join(full, "relsvc", "main.go")); err != nil {
		t.Error("service was not generated into the relative output dir")
	}
}

// --git writes .gitignore into the target. A relative --output-dir once made
// that path resolve twice and generation failed after creating the service.
func TestGitInitWithRelativeOutputDir(t *testing.T) {

	rel := filepath.Join("files", "testtmp-git")

	full := filepath.Join(repoRoot, rel)
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer os.RemoveAll(full)

	cmd := exec.Command(binary, "new", "gitsvc",
		"--module", "github.com/test-org/gitsvc",
		"--output-dir", rel,
		"--go-mod=false")
	cmd.Dir = repoRoot

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git init with a relative output dir failed: %v\n%s", err, out)
	}

	svc := filepath.Join(full, "gitsvc")

	if _, err := os.Stat(filepath.Join(svc, ".gitignore")); err != nil {
		t.Error(".gitignore was not created")
	}

	branch := exec.Command("git", "branch", "--show-current")
	branch.Dir = svc

	name, err := branch.Output()
	if err != nil {
		t.Fatalf("git branch: %v", err)
	}
	if strings.TrimSpace(string(name)) != "dev" {
		t.Errorf("branch = %q, want dev", strings.TrimSpace(string(name)))
	}
}

// --- database driver ---------------------------------------------------------

func TestDatabaseDriverRendersTheRightStack(t *testing.T) {

	cases := []struct {
		driver     string
		wantInDB   []string
		wantInMain string
		wantInMod  string
		wantPort   string
	}{
		{
			driver:     "mysql",
			wantInDB:   []string{`_ "github.com/go-sql-driver/mysql"`, `const Driver = "mysql"`, "DBSystemMySQL", "@tcp("},
			wantInMain: "migrate/v4/database/mysql",
			wantInMod:  "github.com/go-sql-driver/mysql",
			wantPort:   "DATABASE_PORT: 3306",
		},
		{
			driver:     "postgres",
			wantInDB:   []string{`_ "github.com/lib/pq"`, `const Driver = "postgres"`, "DBSystemPostgreSQL", "sslmode"},
			wantInMain: "migrate/v4/database/postgres",
			wantInMod:  "github.com/lib/pq",
			wantPort:   "DATABASE_PORT: 5432",
		},
	}

	for _, tc := range cases {
		t.Run(tc.driver, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--db-driver", tc.driver)

			for _, want := range tc.wantInDB {
				if !fileContains(t, dir, "app/database/database.go", want) {
					t.Errorf("database.go missing %q", want)
				}
			}
			if !fileContains(t, dir, "main.go", tc.wantInMain) {
				t.Errorf("main.go missing migrate driver %q", tc.wantInMain)
			}
			if !fileContains(t, dir, "go.mod", tc.wantInMod) {
				t.Errorf("go.mod missing driver dependency %q", tc.wantInMod)
			}
			if !fileContains(t, dir, "docker-compose-local.yml", tc.wantPort) {
				t.Errorf("compose missing %q: the port default must follow the driver", tc.wantPort)
			}
		})
	}
}

func TestExplicitDBPortOverridesDriverDefault(t *testing.T) {

	dir := mustGenerate(t, "svc", "--db-driver", "postgres", "--db-port", "6543")

	if !fileContains(t, dir, "docker-compose-local.yml", "DATABASE_PORT: 6543") {
		t.Error("--db-port must win over the driver default")
	}
}

func TestUnsupportedDriverIsRejected(t *testing.T) {

	_, out, err := generate(t, "svc", "--db-driver", "sqlite")
	if err == nil {
		t.Fatal("an unsupported driver must fail")
	}
	if !strings.Contains(out, "sqlite") {
		t.Errorf("error must name the rejected driver, got:\n%s", out)
	}
}

func TestDriverAppliesToEveryType(t *testing.T) {

	for _, serviceType := range []string{"general", "casino", "payment"} {
		t.Run(serviceType, func(t *testing.T) {

			dir := mustGenerate(t, "svc", "--type", serviceType, "--db-driver", "postgres")

			if !fileContains(t, dir, "app/database/database.go", `const Driver = "postgres"`) {
				t.Error("the driver choice must apply regardless of service type")
			}
		})
	}
}

// --- subcommands -------------------------------------------------------------

func TestTypesSubcommandListsEveryType(t *testing.T) {

	cmd := exec.Command(binary, "types")
	cmd.Dir = repoRoot

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("types failed: %v\n%s", err, out)
	}

	for _, want := range []string{"general", "casino", "payment"} {
		if !strings.Contains(string(out), want) {
			t.Errorf("types output missing %q:\n%s", want, out)
		}
	}
}

func TestVersionSubcommand(t *testing.T) {

	cmd := exec.Command(binary, "version")
	cmd.Dir = repoRoot

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("version failed: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "gomicrogen") {
		t.Errorf("unexpected version output: %s", out)
	}
}

func TestModuleFlagIsRequired(t *testing.T) {

	cmd := exec.Command(binary, "new", "svc", "--output-dir", t.TempDir(), "--git=false", "--go-mod=false")
	cmd.Dir = repoRoot

	if out, err := cmd.CombinedOutput(); err == nil {
		t.Fatalf("--module is required, expected failure:\n%s", out)
	}
}

func TestExistingServiceIsNotOverwrittenWithoutForce(t *testing.T) {

	out := t.TempDir()

	args := []string{"new", "svc", "--module", "github.com/test-org/svc",
		"--output-dir", out, "--git=false", "--go-mod=false"}

	first := exec.Command(binary, args...)
	first.Dir = repoRoot
	if o, err := first.CombinedOutput(); err != nil {
		t.Fatalf("first generation failed: %v\n%s", err, o)
	}

	second := exec.Command(binary, args...)
	second.Dir = repoRoot
	if o, err := second.CombinedOutput(); err == nil {
		t.Fatalf("regenerating without --force must fail:\n%s", o)
	}
}
