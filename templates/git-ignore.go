package templates

import (
	"fmt"
	"os"
)

var gitIgnore = `
.env
`

func CreateGitIgnore() {
	file, err := os.Create(".gitignore")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	_, err = file.WriteString(gitIgnore)
	if err != nil {
		fmt.Println(err)
		return
	}
}
