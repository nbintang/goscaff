package renames

import (
	"fmt"
	"os"

	"github.com/nbintang/goscaff/tools"
)

func main() {
	dir := "./internal/scaffold/templates"

	// 1) preview dulu
	// if err := RenameAllGoToTmpl(dir, true); err != nil {
	// 	fmt.Println("error:", err)
	// 	os.Exit(1)
	// }

	// 2) kalau sudah yakin, jalankan real rename
	if err := tools.RenameAllGoToTmpl(dir, false); err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}
}
