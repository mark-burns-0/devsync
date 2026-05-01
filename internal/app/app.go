package app

import (
	"github.com/k0kubun/pp/v3"
	"github.com/mark-burns-0/devsync/internal/config"
	"github.com/mark-burns-0/devsync/internal/scanner"
)

func Run(cfg *config.SyncConfig) {
	scanner := scanner.New(cfg)
	dirs, err := scanner.ScanDirs()

	if err != nil {
		panic(err)
	}
	pp.Print(len(dirs))

	// cmd := exec.Command("ls", "-la")
	// cmd.Dir = "."
	// cmd.Stdout = os.Stdout

	// err = cmd.Run()
	// if err != nil {
	// 	fmt.Println(err)
	// }
}
