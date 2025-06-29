package config

// ServiceConfig holds the configuration for generating a new microservice
type ServiceConfig struct {
	ServiceName         string
	ModuleName          string
	Description         string
	Version             string
	Author              string
	Port                string
	GRPCPort            string
	DatabaseDriver      string
	DatabaseURL         string
	DatabaseHost        string
	DatabasePort        string
	DatabasePassword    string
	RedisURL            string
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
		Version:             "1.0.0",
		Author:              "Choplife Group",
		Port:                "8080",
		GRPCPort:            "8081",
		DatabaseDriver:      "mysql",
		DatabaseURL:         "root:password@tcp(localhost:3306)/" + serviceName + "?parseTime=true",
		DatabaseHost:        "localhost",
		DatabasePort:        "5432",
		DatabasePassword:    "mysql",
		RedisURL:            "localhost:6379",
		RedisHost:           "loc7alhost",
		RedisPort:           "639",
		RedisDatabaseNumber: "0",
		RedisPassword:       "",
		Environment:         "development",
	}
}
