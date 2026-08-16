package e2e

// Shared harness: builds the CLI once, generates services, compiles them, and
// starts the backing services every generated project expects.

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	binary   string
	repoRoot string

	infra *backing
)

// backing holds the dependencies a generated service connects to.
type backing struct {
	mysqlHost, mysqlPort       string
	postgresHost, postgresPort string
	redisHost, redisPort       string
	rabbitHost, rabbitPort     string

	terminate []func()
}

func TestMain(m *testing.M) {

	if os.Getenv("GOMICROGEN_E2E") == "" {
		fmt.Println("skipping e2e: set GOMICROGEN_E2E=1 (requires docker)")
		os.Exit(0)
	}

	root, err := filepath.Abs("../..")
	if err != nil {
		panic(err)
	}
	repoRoot = root

	dir, err := os.MkdirTemp("", "gomicrogen-e2e")
	if err != nil {
		panic(err)
	}

	binary = filepath.Join(dir, "gomicrogen")

	build := exec.Command("go", "build", "-o", binary, ".")
	build.Dir = repoRoot
	if out, err := build.CombinedOutput(); err != nil {
		panic("failed to build gomicrogen: " + string(out))
	}

	ctx := context.Background()

	infra, err = startBacking(ctx)
	if err != nil {
		fmt.Printf("failed to start backing services: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()

	for _, stop := range infra.terminate {
		stop()
	}
	os.RemoveAll(dir)

	os.Exit(code)
}

func startBacking(ctx context.Context) (*backing, error) {

	b := &backing{}

	start := func(req testcontainers.ContainerRequest, port string) (string, string, error) {

		c, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: req,
			Started:          true,
		})
		if err != nil {
			return "", "", err
		}

		b.terminate = append(b.terminate, func() { _ = c.Terminate(context.Background()) })

		host, err := c.Host(ctx)
		if err != nil {
			return "", "", err
		}

		mapped, err := c.MappedPort(ctx, port+"/tcp")
		if err != nil {
			return "", "", err
		}

		return host, mapped.Port(), nil
	}

	var err error

	b.mysqlHost, b.mysqlPort, err = start(testcontainers.ContainerRequest{
		Image:        "mysql:8.0",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "secret",
			"MYSQL_DATABASE":      "servicedb",
		},
		WaitingFor: wait.ForLog("port: 3306  MySQL Community Server").WithStartupTimeout(3 * time.Minute),
	}, "3306")
	if err != nil {
		return nil, fmt.Errorf("mysql: %w", err)
	}

	b.postgresHost, b.postgresPort, err = start(testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_PASSWORD": "secret",
			"POSTGRES_DB":       "servicedb",
			"POSTGRES_USER":     "postgres",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp").WithStartupTimeout(2 * time.Minute),
	}, "5432")
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}

	b.redisHost, b.redisPort, err = start(testcontainers.ContainerRequest{
		Image:        "redis:7-alpine",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForListeningPort("6379/tcp").WithStartupTimeout(2 * time.Minute),
	}, "6379")
	if err != nil {
		return nil, fmt.Errorf("redis: %w", err)
	}

	b.rabbitHost, b.rabbitPort, err = start(testcontainers.ContainerRequest{
		Image:        "rabbitmq:3-management-alpine",
		ExposedPorts: []string{"5672/tcp"},
		WaitingFor:   wait.ForListeningPort("5672/tcp").WithStartupTimeout(3 * time.Minute),
	}, "5672")
	if err != nil {
		return nil, fmt.Errorf("rabbitmq: %w", err)
	}

	return b, nil
}

// generate runs the CLI and returns the generated project directory.
func generate(t *testing.T, name string, args ...string) string {
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
	cmd.Dir = repoRoot

	if combined, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("generate %s: %v\n%s", name, err, combined)
	}

	return filepath.Join(out, name)
}

// compile resolves dependencies and builds the generated service.
func compile(t *testing.T, dir string) string {
	t.Helper()

	for _, step := range [][]string{
		{"go", "mod", "tidy"},
		{"go", "vet", "./..."},
	} {
		cmd := exec.Command(step[0], step[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), "GOFLAGS=-mod=mod")

		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("%s failed:\n%s", strings.Join(step, " "), out)
		}
	}

	bin := filepath.Join(t.TempDir(), "service")

	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = dir
	build.Env = append(os.Environ(), "GOFLAGS=-mod=mod")

	if out, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build failed:\n%s", out)
	}

	return bin
}

// serviceEnv is the environment a generated service reads on startup.
func serviceEnv(driver, dbName, httpPort, grpcPort string) []string {

	dbHost, dbPort := infra.mysqlHost, infra.mysqlPort
	dbUser := "root"

	if driver == "postgres" {
		dbHost, dbPort = infra.postgresHost, infra.postgresPort
		dbUser = "postgres"
	}

	return append(os.Environ(),
		"DATABASE_HOST="+dbHost,
		"DATABASE_PORT="+dbPort,
		"DATABASE_USERNAME="+dbUser,
		"DATABASE_PASSWORD=secret",
		"DATABASE_NAME="+dbName,
		"DATABASE_SSL_MODE=disable",
		"REDIS_HOST="+infra.redisHost,
		"REDIS_PORT="+infra.redisPort,
		"REDIS_DATABASE_NUMBER=0",
		"GLOBAL_REDIS_HOST="+infra.redisHost,
		"GLOBAL_REDIS_PORT="+infra.redisPort,
		"GLOBAL_REDIS_DATABASE_NUMBER=1",
		"RABBITMQ_HOST="+infra.rabbitHost,
		"RABBITMQ_PORT="+infra.rabbitPort,
		"RABBITMQ_USER=guest",
		"RABBITMQ_PASS=guest",
		"RABBITMQ_VHOST=",
		"SYSTEM_HOST=0.0.0.0",
		"SYSTEM_PORT="+httpPort,
		"SYSTEM_GRPC_PORT="+grpcPort,
		"SESSION_SECRET=test-secret",
		"QUEUES=",
		"ENV=test",
	)
}

// runService starts the binary and waits for it to answer, returning its log.
func runService(t *testing.T, bin, dir string, env []string, httpPort string) string {
	t.Helper()

	logPath := filepath.Join(t.TempDir(), "service.log")

	logFile, err := os.Create(logPath)
	if err != nil {
		t.Fatalf("create log: %v", err)
	}

	cmd := exec.Command(bin)
	cmd.Dir = dir // so GetRootPath finds migrations/
	cmd.Env = env
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	if err := cmd.Start(); err != nil {
		t.Fatalf("start service: %v", err)
	}

	t.Cleanup(func() {
		_ = cmd.Process.Kill()
		_, _ = cmd.Process.Wait()
		logFile.Close()
	})

	if !waitForHTTP("http://127.0.0.1:"+httpPort+"/", 60*time.Second) {
		body, _ := os.ReadFile(logPath)
		t.Fatalf("service never became healthy on port %s\n--- log ---\n%s", httpPort, body)
	}

	return logPath
}

func waitForHTTP(url string, timeout time.Duration) bool {

	deadline := time.Now().Add(timeout)

	for time.Now().Before(deadline) {

		resp, err := http.Get(url)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return true
			}
		}

		time.Sleep(500 * time.Millisecond)
	}

	return false
}

func get(t *testing.T, url string) (int, string) {
	t.Helper()

	resp, err := http.Get(url)
	if err != nil {
		t.Fatalf("GET %s: %v", url, err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	return resp.StatusCode, string(body)
}

func portIsOpen(port string) bool {

	conn, err := net.DialTimeout("tcp", "127.0.0.1:"+port, 2*time.Second)
	if err != nil {
		return false
	}

	conn.Close()

	return true
}
