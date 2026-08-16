package config

import "fmt"

// Supported database drivers. go-utils speaks both dialects: goutils.Db has a
// Dialect field that switches placeholders, RETURNING and MySQL-only SQL modes.
const (
	DriverMySQL    = "mysql"
	DriverPostgres = "postgres"
)

// SupportedDrivers lists the drivers a generated service can be built for.
var SupportedDrivers = []string{DriverMySQL, DriverPostgres}

// ValidateDriver reports whether the driver is one a generated service supports.
func ValidateDriver(driver string) error {

	for _, supported := range SupportedDrivers {
		if driver == supported {
			return nil
		}
	}

	return fmt.Errorf(`❌ Unsupported database driver %q

📦 Supported drivers: %v`, driver, SupportedDrivers)
}

// DefaultDatabasePort is the conventional port for a driver.
func DefaultDatabasePort(driver string) string {

	if driver == DriverPostgres {
		return "5432"
	}

	return "3306"
}

// IsPostgres lets templates branch on the driver.
func (c *ServiceConfig) IsPostgres() bool { return c.DatabaseDriver == DriverPostgres }

// ServiceConfig holds the configuration for generating a new microservice
type ServiceConfig struct {
	ServiceName         string
	ModuleName          string
	Description         string
	Type                string
	Version             string
	Port                string
	GRPCPort            string
	DatabaseDriver      string
	DatabaseHost        string
	DatabasePort        string
	DatabasePassword    string
	RedisHost           string
	RedisPort           string
	RedisDatabaseNumber string
	RedisPassword       string
	Environment         string
}

// NewServiceConfig creates a new ServiceConfig with default values
func NewServiceConfig(serviceName string) *ServiceConfig {
	return &ServiceConfig{
		ServiceName:         serviceName,
		ModuleName:          "github.com/Choplife-group/" + serviceName,
		Description:         serviceName + " microservice",
		Type:                "general",
		Version:             "1.0.0",
		Port:                "8080",
		GRPCPort:            "8081",
		DatabaseDriver:      DriverMySQL,
		DatabaseHost:        "localhost",
		DatabasePort:        DefaultDatabasePort(DriverMySQL),
		DatabasePassword:    "mysql",
		RedisHost:           "localhost",
		RedisPort:           "6379",
		RedisDatabaseNumber: "0",
		RedisPassword:       "",
		Environment:         "development",
	}
}
