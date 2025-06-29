package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gomicrogen",
	Short: "A CLI tool to scaffold Go microservice projects",
	Long: `gomicrogen is a command line tool that helps you quickly scaffold 
new Go microservice projects with a predefined folder structure and files.`,
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
}
