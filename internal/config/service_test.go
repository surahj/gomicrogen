package config

import "testing"

func TestDefaultsAreUsable(t *testing.T) {

	c := NewServiceConfig("my-service")

	// These defaults land in docker-compose-local.yml, so a typo here ships a
	// service that cannot reach its dependencies.
	cases := map[string]struct{ got, want string }{
		"ServiceName":         {c.ServiceName, "my-service"},
		"Type":                {c.Type, "general"},
		"Version":             {c.Version, "1.0.0"},
		"Port":                {c.Port, "8080"},
		"GRPCPort":            {c.GRPCPort, "8081"},
		"DatabaseDriver":      {c.DatabaseDriver, DriverMySQL},
		"DatabaseHost":        {c.DatabaseHost, "localhost"},
		"DatabasePort":        {c.DatabasePort, "3306"},
		"RedisHost":           {c.RedisHost, "localhost"},
		"RedisPort":           {c.RedisPort, "6379"},
		"RedisDatabaseNumber": {c.RedisDatabaseNumber, "0"},
		"Environment":         {c.Environment, "development"},
	}

	for field, tc := range cases {
		if tc.got != tc.want {
			t.Errorf("%s = %q, want %q", field, tc.got, tc.want)
		}
	}
}

func TestValidateDriver(t *testing.T) {

	for _, ok := range SupportedDrivers {
		if err := ValidateDriver(ok); err != nil {
			t.Errorf("ValidateDriver(%q) errored: %v", ok, err)
		}
	}

	for _, bad := range []string{"sqlite", "mssql", "", "MySQL"} {
		if err := ValidateDriver(bad); err == nil {
			t.Errorf("ValidateDriver(%q) must reject an unsupported driver", bad)
		}
	}
}

func TestDefaultDatabasePortFollowsDriver(t *testing.T) {

	if got := DefaultDatabasePort(DriverPostgres); got != "5432" {
		t.Errorf("postgres default port = %q, want 5432", got)
	}
	if got := DefaultDatabasePort(DriverMySQL); got != "3306" {
		t.Errorf("mysql default port = %q, want 3306", got)
	}
}

func TestIsPostgres(t *testing.T) {

	c := NewServiceConfig("svc")

	if c.IsPostgres() {
		t.Error("default config must not be postgres")
	}

	c.DatabaseDriver = DriverPostgres
	if !c.IsPostgres() {
		t.Error("IsPostgres must report true for the postgres driver")
	}
}
