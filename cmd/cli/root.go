package cli

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goscaff",
	Short: "A Go CLI tool to scaffold backend projects",
	Long: `
                                            
 ██████╗  ██████╗ ███████╗ ██████╗ █████╗ ███████╗███████╗
██╔════╝ ██╔═══██╗██╔════╝██╔════╝██╔══██╗██╔════╝██╔════╝
██║  ███╗██║   ██║███████╗██║     ███████║█████╗  █████╗  
██║   ██║██║   ██║╚════██║██║     ██╔══██║██╔══╝  ██╔══╝  
╚██████╔╝╚██████╔╝███████║╚██████╗██║  ██║██║     ██║     
 ╚═════╝  ╚═════╝ ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝     ╚═╝     

                                            
Goscaff is an instant Go scaffolding CLI.
Use it to generate a production-ready Go backend boilerplate.

Wizard:
  - framework     : Gin or Fiber
  - database      : PostgreSQL or MySQL
  - architecture  : Modular, Layered, or Full Setup
  - dependency DI : None, Uber Dig, or Uber Fx

Quick start:
  goscaff new myapp
  goscaff new myapp --module github.com/you/myapp
  goscaff new myapp --framework gin --db postgres --architecture modular --di uber-fx

Tips:
  - Run "goscaff new --help" to see all flags and examples.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	enableColoredHelp(rootCmd)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	defer func() {
		if r := recover(); r != nil {
			fmt.Fprintln(os.Stderr, "Error: interactive prompt crashed. Try using flags --framework/--db/--architecture/--di/--module.")
			os.Exit(1)
		}
	}()
	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {

		if errors.Is(err, context.Canceled) || errors.Is(err, promptui.ErrInterrupt) {
			fmt.Fprintln(os.Stderr, "✋ Cancelled by user. Nothing was generated.")
			return
		}

		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
