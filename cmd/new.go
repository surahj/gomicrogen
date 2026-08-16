package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Choplife-group/gomicrogen/internal/config"
	"github.com/Choplife-group/gomicrogen/internal/generator"
	"github.com/spf13/cobra"
)

var (
	moduleName          string
	description         string
	serviceType         string
	version             string
	databaseDriver      string
	port                string
	grpcPort            string
	databaseHost        string
	databasePort        string
	databasePassword    string
	redisHost           string
	redisPort           string
	redisDatabaseNumber string
	redisPassword       string
	environment         string
	outputDir           string
	initGit             bool
	runGoMod            bool
	forceOverwrite      bool
)

var newCmd = &cobra.Command{
	Use:   "new [service-name]",
	Short: "Create a new Go microservice project",
	Long: `Create a new Go microservice project with the specified name.
This will generate a complete project structure with all necessary files.

The generated project includes:
• Complete folder structure (app/, docs/, migrations/)
• Docker and Docker Compose configurations
• Migrations applied automatically at startup
• API documentation with Swagger/OpenAPI
• Hot reload development setup
• Tracing, JSON access logs and a Prometheus /metrics endpoint
• Real client-IP resolution and an IP-keyed rate limiter
• Git repository initialization
• Go module management

Examples:
  # Basic microservice
  gomicrogen new user-service --module github.com/choplife-group/user-service

  # Typed service
  gomicrogen new pawapay-service \
    --module github.com/choplife-group/pawapay-service \
    --type payment

  # With custom configuration
  gomicrogen new payment-service \
    --module github.com/choplife-group/payment-service \
    --type payment \
    --description "Payment processing microservice" \
    --version "2.1.0" \
    --port "3000" \
    --grpc-port "3001" \
    --db-driver "postgres" \
    --db-host "localhost" \
    --db-password "secret" \
    --redis-host "localhost" \
    --redis-port "6379" \
    --env "development"

  # In custom directory
  gomicrogen new auth-service \
    --module github.com/choplife-group/auth-service \
    --output-dir /path/to/projects

  # Force overwrite existing project
  gomicrogen new my-service --module github.com/choplife-group/my-service --force

  # Skip Git and Go module initialization
  gomicrogen new my-service --module github.com/choplife-group/my-service --git=false --go-mod=false`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		serviceName := args[0]

		// Determine target directory
		var targetDir string
		if outputDir != "" {
			// Use specified output directory
			targetDir = filepath.Join(outputDir, serviceName)
		} else {
			// Use current working directory
			cwd, err := os.Getwd()
			if err != nil {
				return fmt.Errorf("failed to get current directory: %w", err)
			}
			targetDir = filepath.Join(cwd, serviceName)
		}

		// Find templates directory
		templatesDir := findTemplatesDir()
		if templatesDir == "" {
			execPath, _ := os.Executable()
			execDir := filepath.Dir(execPath)

			return fmt.Errorf(`❌ Templates directory not found

🔍 Searched in the following locations:
   • %s/templates
   • ./templates
   • ../templates
   • ../../templates

💡 This usually happens when:
   • The binary was installed without templates
   • You're running from a different directory
   • Templates were not included in the release

🛠️ Solutions:
   • Reinstall gomicrogen: curl -fsSL https://raw.githubusercontent.com/surahj/gomicrogen/main/install-oneline.sh | bash
   • Download from GitHub releases: https://github.com/surahj/gomicrogen/releases
   • Clone the repository and run from project directory: git clone https://github.com/surahj/gomicrogen && cd gomicrogen && ./gomicrogen new ...`,
				execDir)
		}

		// Resolve the templates layout and the requested service type BEFORE
		// touching the target directory, so a mistyped --type can never trigger
		// the --force removal below
		layout := generator.ResolveLayout(templatesDir)

		canonicalType, overlayDir, err := layout.ResolveType(serviceType)
		if err != nil {
			return err
		}

		// Check if directory already exists
		if err := checkExistingService(serviceName, targetDir); err != nil {
			if !forceOverwrite {
				return err
			} else {
				fmt.Printf("⚠️  Service '%s' already exists. Overwriting due to --force flag...\n", serviceName)
				// Remove existing directory
				if err := os.RemoveAll(targetDir); err != nil {
					return fmt.Errorf("failed to remove existing directory %s: %w", targetDir, err)
				}
			}
		}

		// Create service configuration
		serviceConfig := config.NewServiceConfig(serviceName)
		serviceConfig.Type = canonicalType

		// Override defaults with provided flags
		if moduleName != "" {
			serviceConfig.ModuleName = moduleName
		}
		if description != "" {
			serviceConfig.Description = description
		}
		if version != "" {
			serviceConfig.Version = version
		}
		if port != "" {
			serviceConfig.Port = port
		}
		if grpcPort != "" {
			serviceConfig.GRPCPort = grpcPort
		}
		if databaseDriver != "" {

			if err := config.ValidateDriver(databaseDriver); err != nil {
				return err
			}

			serviceConfig.DatabaseDriver = databaseDriver

			// the conventional port follows the driver unless --db-port says otherwise
			serviceConfig.DatabasePort = config.DefaultDatabasePort(databaseDriver)
		}
		if databaseHost != "" {
			serviceConfig.DatabaseHost = databaseHost
		}
		if databasePort != "" {
			serviceConfig.DatabasePort = databasePort
		}
		if databasePassword != "" {
			serviceConfig.DatabasePassword = databasePassword
		}
		if redisHost != "" {
			serviceConfig.RedisHost = redisHost
		}
		if redisPort != "" {
			serviceConfig.RedisPort = redisPort
		}
		if redisDatabaseNumber != "" {
			serviceConfig.RedisDatabaseNumber = redisDatabaseNumber
		}
		if redisPassword != "" {
			serviceConfig.RedisPassword = redisPassword
		}
		if environment != "" {
			serviceConfig.Environment = environment
		}

		// Create template generator
		gen := generator.NewTemplateGenerator(layout, overlayDir, serviceConfig)

		// Generate the service
		fmt.Printf("Generating %s microservice...\n", serviceName)
		if err := gen.GenerateService(targetDir); err != nil {
			return fmt.Errorf("failed to generate service: %w", err)
		}

		// Initialize Go module if requested
		if runGoMod {
			if err := initializeGoModule(targetDir, serviceConfig.ModuleName); err != nil {
				return fmt.Errorf("failed to initialize Go module: %w", err)
			}
		}

		// Initialize Git repository if requested
		if initGit {
			if err := initializeGitRepo(targetDir); err != nil {
				return fmt.Errorf("failed to initialize Git repository: %w", err)
			}
		}

		fmt.Printf("\n✅ Successfully created %s microservice!\n", serviceName)
		fmt.Printf("📁 Project location: %s\n", targetDir)
		fmt.Printf("🚀 To get started:\n")
		fmt.Printf("   cd %s\n", targetDir)
		if !runGoMod {
			fmt.Printf("   go mod tidy\n")
		}
		fmt.Printf("   go run main.go\n")

		return nil
	},
}

// checkExistingService checks if a service with the given name already exists
func checkExistingService(serviceName, targetDir string) error {
	// Check if directory exists
	if _, err := os.Stat(targetDir); err == nil {
		// Check if it contains Go files (indicating it's a Go project)
		goFiles, err := filepath.Glob(filepath.Join(targetDir, "*.go"))
		if err == nil && len(goFiles) > 0 {
			return fmt.Errorf(`❌ Go service "%s" already exists at: %s

📁 This appears to be an existing Go project with files:
   %s

💡 To resolve this, you can:
   • Use a different service name
   • Remove the existing directory: rm -rf %s
   • Use --force flag to overwrite: gomicrogen new %s --force`,
				serviceName, targetDir, strings.Join(goFiles, "\n   "), serviceName, serviceName)
		}

		// Check if it contains a go.mod file
		if _, err := os.Stat(filepath.Join(targetDir, "go.mod")); err == nil {
			return fmt.Errorf(`❌ Go module "%s" already exists at: %s

📁 This appears to be an existing Go module with go.mod file.

💡 To resolve this, you can:
   • Use a different service name
   • Remove the existing directory: rm -rf %s
   • Use --force flag to overwrite: gomicrogen new %s --force`,
				serviceName, targetDir, serviceName, serviceName)
		}

		// Generic directory exists error
		return fmt.Errorf(`❌ Directory "%s" already exists at: %s

💡 To resolve this, you can:
   • Use a different service name
   • Remove the existing directory: rm -rf %s
   • Use --force flag to overwrite: gomicrogen new %s --force`,
			serviceName, targetDir, serviceName, serviceName)
	}
	return nil
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Required flags
	newCmd.Flags().StringVarP(&moduleName, "module", "m", "", "Go module name (e.g., github.com/choplife-group/service-name)")
	newCmd.MarkFlagRequired("module")

	// Service configuration flags
	newCmd.Flags().StringVarP(&serviceType, "type", "t", "general", typeFlagUsage())
	newCmd.Flags().StringVarP(&description, "description", "d", "", "Service description (e.g., 'User management microservice')")
	newCmd.Flags().StringVarP(&version, "version", "v", "1.0.0", "Service version (e.g., '2.1.0')")
	newCmd.Flags().StringVarP(&port, "port", "p", "8080", "HTTP port for the service")
	newCmd.Flags().StringVarP(&grpcPort, "grpc-port", "g", "8081", "gRPC port for the service")
	newCmd.Flags().StringVarP(&environment, "env", "e", "development", "Environment (development, staging, production)")

	// Database configuration flags
	newCmd.Flags().StringVarP(&databaseDriver, "db-driver", "", "mysql", "Database driver (mysql, postgres)")
	newCmd.Flags().StringVarP(&databaseHost, "db-host", "", "localhost", "Database host")
	newCmd.Flags().StringVarP(&databasePort, "db-port", "", "", "Database port (3306 for MySQL, 5432 for PostgreSQL)")
	newCmd.Flags().StringVarP(&databasePassword, "db-password", "", "", "Database password")

	// Redis configuration flags
	newCmd.Flags().StringVarP(&redisHost, "redis-host", "", "localhost", "Redis host")
	newCmd.Flags().StringVarP(&redisPort, "redis-port", "", "6379", "Redis port")
	newCmd.Flags().StringVarP(&redisDatabaseNumber, "redis-db-number", "", "0", "Redis database number (0-15)")
	newCmd.Flags().StringVarP(&redisPassword, "redis-password", "", "", "Redis password")

	// Output and behavior flags
	newCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "", "Output directory (default: current directory)")
	newCmd.Flags().BoolVarP(&initGit, "git", "", true, "Initialize Git repository with dev branch")
	newCmd.Flags().BoolVarP(&runGoMod, "go-mod", "", true, "Run go mod init and go mod tidy")
	newCmd.Flags().BoolVarP(&forceOverwrite, "force", "", false, "Force overwrite if service already exists")
}

// typeFlagUsage builds the --type help text. Cobra assembles usage strings at
// init(), before the templates directory is known, so this is best-effort: when
// discovery fails it points at the types subcommand instead.
func typeFlagUsage() string {

	const fallback = "Service type; run 'gomicrogen types' to list available types"

	templatesDir := findTemplatesDir()
	if templatesDir == "" {
		return fallback
	}

	names := []string{generator.GeneralType}

	for _, t := range generator.ResolveLayout(templatesDir).Types() {

		if t.Name == generator.GeneralType {
			continue
		}

		names = append(names, t.Name)
	}

	if len(names) == 1 {
		return fallback
	}

	return fmt.Sprintf("Service type: %s", strings.Join(names, ", "))
}

// findTemplatesDir finds the templates directory relative to the executable
func findTemplatesDir() string {
	// Get the executable path
	execPath, err := os.Executable()
	if err != nil {
		// Fallback to current directory
		execPath = "."
	}

	// Get the directory of the executable
	execDir := filepath.Dir(execPath)

	// Try to find templates in the executable directory
	templatesPath := filepath.Join(execDir, "templates")
	if _, err := os.Stat(templatesPath); err == nil {
		return templatesPath
	}

	// Try to find templates in the current directory
	if _, err := os.Stat("templates"); err == nil {
		return "templates"
	}

	// Try to find templates in the parent directory (for development)
	if _, err := os.Stat("../templates"); err == nil {
		return "../templates"
	}

	// Try to find templates in the go-template directory (for this project)
	if _, err := os.Stat("../go-template"); err == nil {
		return "../go-template"
	}

	// If we're in development mode, try to find templates relative to the project root
	// Look for templates in common development locations
	devPaths := []string{
		"templates",
		"../templates",
		"../../templates",
		"./templates",
	}

	for _, path := range devPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	return ""
}

// initializeGoModule initializes the Go module in the target directory
func initializeGoModule(targetDir, moduleName string) error {
	fmt.Println("Initializing Go module...")

	// Change to target directory
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(targetDir); err != nil {
		return err
	}

	// Check if go.mod already exists. This is relative to the working directory,
	// which is now targetDir: joining targetDir again would break for a relative
	// --output-dir and send us down the go mod init path on a module that exists.
	if _, err := os.Stat("go.mod"); err == nil {
		// go.mod already exists, just run go mod tidy
		fmt.Println("📁 go.mod already exists, running go mod tidy...")
		cmd := exec.Command("go", "mod", "tidy")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run go mod tidy: %w", err)
		}
		fmt.Println("✅ go mod tidy completed successfully")
	} else {
		// go.mod doesn't exist, run go mod init first
		fmt.Println("📁 Creating new go.mod file...")
		cmd := exec.Command("go", "mod", "init", moduleName)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run go mod init: %w", err)
		}

		// Then run go mod tidy
		fmt.Println("📁 Running go mod tidy...")
		cmd = exec.Command("go", "mod", "tidy")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run go mod tidy: %w", err)
		}
		fmt.Println("✅ Go module initialized successfully")
	}

	return nil
}

// initializeGitRepo initializes a Git repository with dev branch
func initializeGitRepo(targetDir string) error {
	fmt.Println("Initializing Git repository...")

	// Change to target directory
	originalDir, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(originalDir)

	if err := os.Chdir(targetDir); err != nil {
		return err
	}

	// Initialize Git repository
	cmd := exec.Command("git", "init")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create .gitignore file
	gitignoreContent := `# Binaries for programs and plugins
*.exe
*.exe~
*.dll
*.so
*.dylib

# Test binary, built with 'go test -c'
*.test

# Output of the go coverage tool, specifically when used with LiteIDE
*.out

# Dependency directories (remove the comment below to include it)
# vendor/

# Go workspace file
go.work

# Environment variables
.env
.env.local
.env.*.local

# Docker compose local file
docker-compose-local.yml

# IDE files
.vscode/
.idea/
*.swp
*.swo

# OS generated files
.DS_Store
.DS_Store?
._*
.Spotlight-V100
.Trashes
ehthumbs.db
Thumbs.db

# Logs
*.log

# Air live reload
tmp/

# Docker
.dockerignore

# Database
*.db
*.sqlite

# Build artifacts
build/
dist/

docker-compose-local.yml
`

	// relative to the working directory, which is now targetDir
	if err := os.WriteFile(".gitignore", []byte(gitignoreContent), 0644); err != nil {
		return fmt.Errorf("failed to create .gitignore: %w", err)
	}
	fmt.Println("📁 Created .gitignore")

	// Add all files except .env and docker-compose-local.yml
	cmd = exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	// Create and switch to dev branch (skip master)
	cmd = exec.Command("git", "checkout", "-b", "dev")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	fmt.Println("✅ Git repository initialized with dev branch")
	fmt.Println("📁 .gitignore created")
	fmt.Println("📁 .env and docker-compose-local.yml excluded from tracking")
	return nil
}
