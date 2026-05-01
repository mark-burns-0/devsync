package main

import (
	"fmt"

	"github.com/mark-burns-0/devsync/internal/app"
	"github.com/mark-burns-0/devsync/internal/config"
)

func main() {
	cfg, err := config.New(".env")
	if err != nil {
		panic(err)
	}
	fmt.Println("Configuration loaded successfully")

	app.Run(cfg)
}
