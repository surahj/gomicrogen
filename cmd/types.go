package cmd

import (
	"fmt"

	"github.com/Choplife-group/gomicrogen/internal/generator"
	"github.com/spf13/cobra"
)

var typesCmd = &cobra.Command{
	Use:   "types",
	Short: "List the service types available to --type",
	Long: `List the service types that can be passed to 'gomicrogen new --type'.

Types are discovered from the installed templates directory, so adding a new one
is a matter of creating templates/types/<name>/ — no rebuild is required.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		templatesDir := findTemplatesDir()
		if templatesDir == "" {
			return fmt.Errorf("❌ Templates directory not found — reinstall gomicrogen or run from the project directory")
		}

		layout := generator.ResolveLayout(templatesDir)

		if layout.Legacy {

			cmd.Println("📦 This installation ships legacy templates with no type overlays.")
			cmd.Println("   Only the default general service can be generated.")
			cmd.Println("\n💡 Reinstall gomicrogen to pick up the current templates.")

			return nil
		}

		cmd.Println("📦 Available service types:")
		cmd.Printf("   • %-10s base microservice (default)\n", generator.GeneralType)

		for _, t := range layout.Types() {

			if t.Name == generator.GeneralType {
				continue
			}

			cmd.Printf("   • %-10s %s\n", t.Name, t.Description)
		}

		cmd.Println("\n💡 Add your own: create templates/types/<name>/ with files mirroring the")
		cmd.Println("   base layout; a file at the same path replaces the base version.")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(typesCmd)
}
