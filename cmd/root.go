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
Use it to generate a production-ready Go backend boilerplate.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
