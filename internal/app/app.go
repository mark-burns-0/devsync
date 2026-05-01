package app

import (
	"log"

	"github.com/mark-burns-0/devsync/internal/config"
	"github.com/mark-burns-0/devsync/internal/scanner"
	"github.com/mark-burns-0/devsync/internal/sync"
)

func Run(cfg *config.SyncConfig) {
	scanner := scanner.New(cfg)
	dirs, err := scanner.ScanDirs()
	if err != nil {
		panic(err)
	}

	syncer := sync.New(dirs)
	if err := syncer.Sync(); err != nil {
		log.Fatal(err)
	}
}
