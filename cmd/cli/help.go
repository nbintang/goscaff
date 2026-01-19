package cli

import "github.com/spf13/cobra"

const (
	clrReset = "\x1b[0m"
	clrBold  = "\x1b[1m"
	clrDim   = "\x1b[2m"

	clrCyan   = "\x1b[36m"
	clrGreen  = "\x1b[32m"
	clrYellow = "\x1b[33m"
	clrBlue   = "\x1b[34m"
	clrGray   = "\x1b[90m"
)

func enableColoredHelp(root *cobra.Command) {
	// Help template (saat `goscaff help` / `goscaff`)
	root.SetHelpTemplate(clrBold + clrCyan + `{{with (or .Long .Short)}}{{.}}{{end}}` + clrReset + "\n\n" +
		clrBold + `Usage:` + clrReset + "\n" +
		"  " + clrBlue + `{{.UseLine}}` + clrReset + "\n\n" +
		clrBold + `Commands:` + clrReset + "\n" +
		`{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}` +
		"  " + clrGreen + `{{rpad .Name .NamePadding}}` + clrReset +
		" " + clrGray + `{{.Short}}` + clrReset + "\n" +
		`{{end}}{{end}}` + "\n" +
		clrBold + `Flags:` + clrReset + "\n" +
		`{{if .HasAvailableLocalFlags}}{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}` + "\n" +
		`{{if .HasAvailableInheritedFlags}}` + clrBold + `Global Flags:` + clrReset + "\n" +
		`{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}` + "\n")

	// Usage template (saat error / invalid args)
	root.SetUsageTemplate(clrBold + `Usage:` + clrReset + "\n" +
		"  " + clrBlue + `{{.UseLine}}` + clrReset + "\n\n" +
		`{{if .HasExample}}` + clrBold + `Examples:` + clrReset + "\n" +
		clrGray + `{{.Example}}` + clrReset + "\n\n" + `{{end}}` +
		`{{if .HasAvailableLocalFlags}}` + clrBold + `Flags:` + clrReset + "\n" +
		`{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}` + "\n\n" + `{{end}}`)
}
