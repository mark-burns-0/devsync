package main

import (
	"github.com/mark-burns-0/devsync/internal/app"
	"github.com/mark-burns-0/devsync/internal/config"
)

func main() {
	cfg, err := config.New(".env")
	if err != nil {
		panic(err)
	}

	app.Run(cfg)
}
