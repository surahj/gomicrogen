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
• Complete project structure (app/, cmd/, docs/, k8s/, etc.)
• Docker and Docker Compose configurations
• Kubernetes deployment manifests
• Database migrations and models
• API documentation with Swagger/OpenAPI
• Hot reload development setup
• Observability integration (OpenTelemetry, Uptrace)
• Redis caching and session management
• Git repository initialization
• Go module management

Examples:
  # Create a basic microservice
  gomicrogen new user-service --module github.com/myorg/user-service

  # Create with custom configuration
  gomicrogen new payment-service \
    --module github.com/myorg/payment-service \
    --description "Payment processing microservice" \
    --port 3000 \
    --db-driver mysql \
    --db-host localhost \
    --redis-host localhost

  # Create in specific directory
  gomicrogen new auth-service \
    --module github.com/myorg/auth-service \
    --output-dir /path/to/projects

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

	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.gomicrogen.yaml)")

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
