package main

import (
	"fmt"
	"os"

	"github.com/nbintang/goscaff/tools"
)

func main() {
	dir := "./internal/templates"
	dryRun := false

	if len(os.Args) > 1 && os.Args[1] != "" {
		dir = os.Args[1]
	}

	if len(os.Args) > 2 && (os.Args[2] == "--dry-run" || os.Args[2] == "-n") {
		dryRun = true
	}

	if err := tools.RenameAllGoToTmpl(dir, dryRun); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
