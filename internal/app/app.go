package app

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/k0kubun/pp/v3"
)

var root = ""
var gitDir = ".git"

func Run() {
	dirs := []string{}
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && d.Name() == ".git" {
			path := strings.TrimSuffix(strings.TrimPrefix(path, root), gitDir)
			cleanPath := strings.Trim(path, string(filepath.Separator))
			dirs = append(dirs, cleanPath)
			return fs.SkipDir
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Ошибка при поиске: %v\n", err)
	}

	pp.Print(dirs)

	// cmd := exec.Command("ls", "-la")
	// cmd.Dir = "."
	// cmd.Stdout = os.Stdout

	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
