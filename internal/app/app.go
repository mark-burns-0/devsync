package app

import (
	"fmt"
	"os/exec"
)

func Run() {
	cmd := exec.Command("ls")
	cmd.Dir = "."
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
	}
}
