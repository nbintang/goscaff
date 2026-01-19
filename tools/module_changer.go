package tools

import ( 
	"bytes"
	"fmt" 
	"os" 
)



func ReplaceModule(path, old string) error {
	input, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if !bytes.Contains(input, []byte(old)) {
		return nil
	}

	output := bytes.ReplaceAll(input, []byte(old), []byte("{{.MODULE_PATH}}"))

	info := "updated: " + path
	fmt.Println(info)

	return os.WriteFile(path, output, 0o644)
}
