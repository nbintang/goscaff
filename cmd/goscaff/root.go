package goscaff

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goscaff",
	Short: "Generate a Go Fiber Boilerplate",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}

}
