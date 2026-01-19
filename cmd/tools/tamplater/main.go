package tamplater

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/nbintang/goscaff/tools"
)

//
//
// go run . <template-folder> <old-module>
//
func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  go run . <template-folder> <old-module>")
		fmt.Println("Example:")
		fmt.Println("  go run . ./templates/full github.com/nbintang/myapp")
		os.Exit(1)
	}

	root := os.Args[1]
	oldModule := os.Args[2]

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// skip go.sum
		if strings.HasSuffix(path, "go.sum") {
			return nil
		}

		return tools.ReplaceModule(path, oldModule)
	})

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	fmt.Println("Done replacing module with {{.MODULE_PATH}}")
}