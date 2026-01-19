package pkg

import "github.com/fatih/color"

var (
	// Header
	ColorHeader = color.New(color.FgCyan, color.Bold)

	// Status
	ColorOk     = color.New(color.FgGreen, color.Bold)
	ColorTip    = color.New(color.FgYellow)
	ColorInfo   = color.New(color.FgHiBlack)
	ColorAction = color.New(color.FgWhite, color.Bold)

	// Sections
	ColorStepTitle = color.New(color.FgWhite, color.Bold)
	ColorBullet    = color.New(color.FgHiBlack)

	// Commands & notes
	ColorCmd  = color.New(color.FgHiBlue)
	ColorNote = color.New(color.FgHiBlack)

	ColorSep = color.New(color.FgHiBlack)
)

func Header(title string) {
	ColorHeader.Println(title)
	ColorSep.Println("────────────────────────────────────────")
}

func Action(label string) {
	ColorInfo.Print("  ")
	ColorAction.Println(label)
}

func Success(msg string) {
	ColorInfo.Print("  ")
	ColorOk.Print("✓ ")
	ColorOk.Println(msg)
}

func Info(msg string) {
	ColorInfo.Print("  ")
	ColorInfo.Println(msg)
}
