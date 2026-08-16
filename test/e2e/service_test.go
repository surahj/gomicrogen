package e2e

// Runtime tests: generate a service, compile it, run it against real MySQL /
// Postgres / Redis / RabbitMQ, and assert it actually serves traffic.

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	// registered so the assertions below can open both databases directly
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

// Every service type must compile and serve its endpoints.
func TestServiceTypesRunAgainstRealDependencies(t *testing.T) {

	cases := []struct {
		serviceType string
		httpPort    string
		grpcPort    string
		wantGRPC    bool
	}{
		{"general", "18080", "18081", false},
		{"casino", "18090", "18091", true},
		{"payment", "18100", "18101", true},
	}

	for _, tc := range cases {
		t.Run(tc.serviceType, func(t *testing.T) {

			dbName := "svc_" + tc.serviceType
			mustCreateMySQLDatabase(t, dbName)

			dir := generate(t, "svc", "--type", tc.serviceType)
			bin := compile(t, dir)

			runService(t, bin, dir, serviceEnv("mysql", dbName, tc.httpPort, tc.grpcPort), tc.httpPort)

			base := "http://127.0.0.1:" + tc.httpPort

			t.Run("status endpoint reports healthy dependencies", func(t *testing.T) {

				code, body := get(t, base+"/")

				if code != http.StatusOK {
					t.Fatalf("GET / = %d, want 200", code)
				}
				for _, want := range []string{"successful ping", "PONG", "svc"} {
					if !strings.Contains(body, want) {
						t.Errorf("status body missing %q: %s", want, body)
					}
				}
			})

			t.Run("swagger docs are served", func(t *testing.T) {
				if code, _ := get(t, base+"/docs/index.html"); code != http.StatusOK {
					t.Errorf("GET /docs/index.html = %d, want 200", code)
				}
			})

			t.Run("prometheus metrics are exposed", func(t *testing.T) {

				code, body := get(t, base+"/metrics")

				if code != http.StatusOK {
					t.Fatalf("GET /metrics = %d, want 200", code)
				}
				if !strings.Contains(body, "echo_requests_total") {
					t.Error("metrics missing the fleet-wide echo_requests_total series")
				}
			})

			t.Run("migrations applied", func(t *testing.T) {
				assertMigrated(t, "mysql", mysqlDSN(dbName))
			})

			t.Run("grpc server only where the type calls for it", func(t *testing.T) {

				if got := portIsOpen(tc.grpcPort); got != tc.wantGRPC {
					t.Errorf("gRPC port open = %v, want %v", got, tc.wantGRPC)
				}
			})
		})
	}
}

// The rate limiter must protect ordinary routes and never the scrape path.
func TestRateLimiterProtectsRoutesButNotMetrics(t *testing.T) {

	dbName := "svc_ratelimit"
	mustCreateMySQLDatabase(t, dbName)

	dir := generate(t, "svc")
	bin := compile(t, dir)

	env := append(serviceEnv("mysql", dbName, "18200", "18201"),
		"RATE_LIMIT=5", "RATE_LIMIT_BURST=2", "RATE_LIMIT_EXPIRES_IN_SECONDS=60")

	runService(t, bin, dir, env, "18200")

	base := "http://127.0.0.1:18200"

	t.Run("ordinary route is limited", func(t *testing.T) {

		limited := 0
		for i := 0; i < 40; i++ {
			if code, _ := get(t, base+"/"); code == http.StatusTooManyRequests {
				limited++
			}
		}

		if limited == 0 {
			t.Error("no request was rate limited; the limiter is not active")
		}
	})

	t.Run("metrics path is exempt", func(t *testing.T) {

		for i := 0; i < 40; i++ {
			if code, _ := get(t, base+"/metrics"); code != http.StatusOK {
				t.Fatalf("/metrics returned %d; scrapes must never be rate limited", code)
			}
		}
	})
}

// The IP extractor must prefer a public forwarded address and ignore private
// ones, otherwise the limiter keys on the proxy instead of the caller.
func TestClientIPExtraction(t *testing.T) {

	dbName := "svc_clientip"
	mustCreateMySQLDatabase(t, dbName)

	dir := generate(t, "svc")
	bin := compile(t, dir)

	logPath := runService(t, bin, dir, serviceEnv("mysql", dbName, "18300", "18301"), "18300")

	cases := []struct {
		name       string
		header     string
		wantLogged string
	}{
		{"public forwarded address is used", "203.0.113.9, 10.0.0.1", "203.0.113.9"},
		{"private forwarded address is ignored", "10.1.2.3", "127.0.0.1"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {

			req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:18300/"+tc.name, nil)
			if err != nil {
				t.Fatalf("request: %v", err)
			}
			req.Header.Set("X-Forwarded-For", tc.header)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("do: %v", err)
			}
			resp.Body.Close()

			time.Sleep(500 * time.Millisecond)

			body, err := os.ReadFile(logPath)
			if err != nil {
				t.Fatalf("read log: %v", err)
			}

			// the access log is JSON and carries remote_ip per request
			want := `"remote_ip":"` + tc.wantLogged + `"`
			if !strings.Contains(string(body), want) {
				t.Errorf("access log has no entry with %s", want)
			}
		})
	}
}

// The postgres driver must produce a service that migrates and serves.
func TestPostgresDriverEndToEnd(t *testing.T) {

	dir := generate(t, "svc", "--type", "payment", "--db-driver", "postgres")
	bin := compile(t, dir)

	runService(t, bin, dir, serviceEnv("postgres", "servicedb", "18400", "18401"), "18400")

	code, body := get(t, "http://127.0.0.1:18400/")
	if code != http.StatusOK {
		t.Fatalf("GET / = %d, want 200", code)
	}
	if !strings.Contains(body, "successful ping") {
		t.Errorf("postgres service is not talking to its database: %s", body)
	}

	assertMigrated(t, "postgres", postgresDSN("servicedb"))
}

// A generated service must find migrations/ when run from its own directory,
// regardless of where it was compiled.
func TestServiceRunsFromRelocatedBinary(t *testing.T) {

	dbName := "svc_relocated"
	mustCreateMySQLDatabase(t, dbName)

	dir := generate(t, "svc")
	bin := compile(t, dir)

	relocated := t.TempDir()

	if err := copyFile(bin, relocated+"/service"); err != nil {
		t.Fatalf("copy binary: %v", err)
	}
	if err := copyDir(dir+"/migrations", relocated+"/migrations"); err != nil {
		t.Fatalf("copy migrations: %v", err)
	}
	if err := os.Chmod(relocated+"/service", 0o755); err != nil {
		t.Fatalf("chmod: %v", err)
	}

	runService(t, relocated+"/service", relocated, serviceEnv("mysql", dbName, "18500", "18501"), "18500")

	if code, _ := get(t, "http://127.0.0.1:18500/"); code != http.StatusOK {
		t.Error("a relocated binary must still find its migrations")
	}
}

// --go-mod runs `go mod init`/`go mod tidy` inside the target. A relative
// --output-dir once made the go.mod existence check resolve twice, so the CLI
// ran `go mod init` on a module that already existed and aborted.
func TestGoModInitWithRelativeOutputDir(t *testing.T) {

	rel := filepath.Join("files", "testtmp-gomod")

	full := filepath.Join(repoRoot, rel)
	if err := os.MkdirAll(full, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	defer os.RemoveAll(full)

	cmd := exec.Command(binary, "new", "modsvc",
		"--module", "github.com/test-org/modsvc",
		"--output-dir", rel,
		"--type", "payment",
		"--git=false")
	cmd.Dir = repoRoot

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go mod init with a relative output dir failed: %v\n%s", err, out)
	}

	if !strings.Contains(string(out), "go mod tidy completed successfully") {
		t.Errorf("expected tidy to run, got:\n%s", out)
	}

	// tidy must have produced a lock file for the resolved dependencies
	if _, err := os.Stat(filepath.Join(full, "modsvc", "go.sum")); err != nil {
		t.Error("go.sum was not produced by go mod tidy")
	}
}

// --- helpers -----------------------------------------------------------------

func mysqlDSN(dbName string) string {

	return fmt.Sprintf("root:secret@tcp(%s:%s)/%s?parseTime=true&multiStatements=true",
		infra.mysqlHost, infra.mysqlPort, dbName)
}

func postgresDSN(dbName string) string {

	return fmt.Sprintf("host=%s port=%s user=postgres password=secret dbname=%s sslmode=disable",
		infra.postgresHost, infra.postgresPort, dbName)
}

func mustCreateMySQLDatabase(t *testing.T, name string) {
	t.Helper()

	db, err := sql.Open("mysql", mysqlDSN(""))
	if err != nil {
		t.Fatalf("open mysql: %v", err)
	}
	defer db.Close()

	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + name); err != nil {
		t.Fatalf("create database %s: %v", name, err)
	}
}

// assertMigrated checks golang-migrate recorded a clean run.
func assertMigrated(t *testing.T, driver, dsn string) {
	t.Helper()

	db, err := sql.Open(driver, dsn)
	if err != nil {
		t.Fatalf("open %s: %v", driver, err)
	}
	defer db.Close()

	var version int
	var dirty bool

	if err := db.QueryRow("SELECT version, dirty FROM schema_migrations").Scan(&version, &dirty); err != nil {
		t.Fatalf("migrations were not applied: %v", err)
	}
	if dirty {
		t.Errorf("schema_migrations is dirty at version %d", version)
	}
}

func copyFile(src, dst string) error {

	body, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, body, 0o755)
}

func copyDir(src, dst string) error {

	if err := os.MkdirAll(dst, 0o755); err != nil {
		return err
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, e := range entries {

		if e.IsDir() {
			continue
		}
		if err := copyFile(src+"/"+e.Name(), dst+"/"+e.Name()); err != nil {
			return err
		}
	}

	return nil
}
