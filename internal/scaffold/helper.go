package scaffold

import (
	"fmt"
	"path/filepath"
)

func printNextSteps(outDir string, opts Options) {
	fmt.Println()

	// Status line
	cOk.Printf("✓ ")
	cOk.Println("Project generated successfully")

	cTip.Print("⚑ ")
	cTip.Println("Tip: review your .env values before running")

	fmt.Println()
	cHeader.Println("Next steps")
	cBullet.Println("────────────────────────────────────────")

	// Normalize outDir biar ga tampil "./myapp"
	projectDir := filepath.Base(filepath.Clean(outDir))

	// Step 1
	printStep("Go to project directory",
		fmt.Sprintf("cd %s", projectDir),
	)

	// Step 2
	printStep("Setup environment",
		"cp .env.example .env",
	)

	// Full preset extras
	if opts.Preset == "full" {
		printStep("Start dependencies (optional)",
			"docker compose up -d",
		)
	}

	// Run
	printStep("Run the app",
		"go run ./cmd/api",
	)

	fmt.Println()
	cNote.Println("  • Server: http://localhost:8080")
	cNote.Println("  • If you changed ports, check .env")
	fmt.Println()
}

func printStep(title string, commands ...string) {
	fmt.Println()
	cStepTitle.Printf("  %s\n", title)
	for _, c := range commands {
		cBullet.Print("    $ ")
		cCmd.Println(c)
	}
}
