package goscaff

 

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nbintang/goscaff-cli/internal/scaffold"

	"github.com/spf13/cobra"
)

var (
	flagModule string
	flagDB     string
)

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new project from the boilerplate",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		outDir := filepath.Join(".", projectName)

		if _, err := os.Stat(outDir); err == nil {
			return fmt.Errorf("directory %s already exists", outDir)
		}

		modulePath := flagModule
		if modulePath == "" {
			modulePath = projectName
		}

		opts := scaffold.Options{
			ProjectName: projectName,
			ModulePath:  modulePath,
			DB:          flagDB,
		}

		return scaffold.Generate(outDir, opts)
},

}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().StringVar(&flagModule, "module", "", "Go module path (required)")
	newCmd.Flags().StringVar(&flagDB, "db", "postgres", "Database: postgres|mysql")
}
