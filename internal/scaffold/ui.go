package scaffold

import "github.com/fatih/color"

var (
	// Header
	cHeader = color.New(color.FgCyan, color.Bold)

	// Status
	cOk     = color.New(color.FgGreen, color.Bold)
	cTip    = color.New(color.FgYellow)
	cInfo   = color.New(color.FgHiBlack)  
	cAction = color.New(color.FgWhite, color.Bold)

	// Sections
	cStepTitle = color.New(color.FgWhite, color.Bold)
	cBullet    = color.New(color.FgHiBlack)  

	// Commands & notes
	cCmd  = color.New(color.FgHiBlue)   
	cNote = color.New(color.FgHiBlack) 

	cSep = color.New(color.FgHiBlack)
)

func header(title string) {
	cHeader.Println(title)
	cSep.Println("────────────────────────────────────────")
}

func action(label string) {
	// contoh: "Render preset (base)"
	cInfo.Print("  ")
	cAction.Println(label)
}

func ok(msg string) {
	cInfo.Print("  ")
	cOk.Print("✓ ")
	cOk.Println(msg)
}

func infoLine(msg string) {
	cInfo.Print("  ")
	cInfo.Println(msg)
}
