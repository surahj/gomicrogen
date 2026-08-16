package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	appVersion = "dev"
	appCommit  = "dev"
	appDate    = ""
)

var rootCmd = &cobra.Command{
	Use:   "gomicrogen",
	Short: "A CLI tool to scaffold Go microservice projects",
	Long: `gomicrogen is a command line tool that helps you quickly scaffold 
new Go microservice projects with a predefined folder structure and files.

This tool generates complete microservice projects with:
• Complete project structure (app/, docs/, migrations/)
• Docker and Docker Compose configurations
• Migrations applied automatically at startup
• API documentation with Swagger/OpenAPI
• Hot reload development setup
• Observability integration (OpenTelemetry, Uptrace, Prometheus)
• Real client-IP resolution and an IP-keyed rate limiter
• Git repository initialization
• Go module management

Examples:
  # Create a basic microservice
  gomicrogen new user-service --module github.com/choplife-group/user-service

  # Create a typed service
  gomicrogen new pawapay-service \
    --module github.com/choplife-group/pawapay-service \
    --type payment

  # Create with custom configuration
  gomicrogen new reporting-service \
    --module github.com/choplife-group/reporting-service \
    --description "Reporting microservice" \
    --port 3000 \
    --db-driver postgres \
    --db-host localhost \
    --redis-host localhost

  # Create in specific directory
  gomicrogen new auth-service \
    --module github.com/choplife-group/auth-service \
    --output-dir /path/to/projects

List the available service types with: gomicrogen types
For detailed help on any command, use: gomicrogen [command] --help`,
	Version: fmt.Sprintf("%s-%s (%s)", appVersion, appCommit, appDate),
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Make built-in --version also available as -v
	if f := rootCmd.Flags().Lookup("version"); f != nil && f.Shorthand == "" {
		f.Shorthand = "v"
	}

	// Explicit version subcommand so users can run: gomicrogen version
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("gomicrogen %s-%s (%s)\n", appVersion, appCommit, appDate)
		},
	})
}
