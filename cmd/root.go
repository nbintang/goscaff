package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goscaff",
	Short: "A Go CLI tool to scaffold backend projects",
	Long: `
                                            
 ▄████   ▄▄▄   ▄▄▄▄  ▄▄▄▄  ▄▄▄  ▄▄▄▄▄ ▄▄▄▄▄ 
██  ▄▄▄ ██▀██ ███▄▄ ██▀▀▀ ██▀██ ██▄▄  ██▄▄  
 ▀███▀  ▀███▀ ▄▄██▀ ▀████ ██▀██ ██    ██    
                                            
Goscaff is an instant Go scaffolding CLI.
Use it to generate a production-ready Go backend boilerplate.

Presets:
  - base : minimal starter (core structure only)
  - full : production-ready starter (default)

Quick start:
  goscaff new myapp --module github.com/you/myapp
  goscaff new myapp --preset base
  goscaff new myapp --preset full --db mysql --module github.com/you/myapp

Tips:
  - Run "goscaff new --help" to see all flags and examples.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	enableColoredHelp(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
