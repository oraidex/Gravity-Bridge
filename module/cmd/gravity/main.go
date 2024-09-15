package main

import (
	"os"

	"github.com/Gravity-Bridge/Gravity-Bridge/module/cmd/gravity/cmd"
	_ "github.com/Gravity-Bridge/Gravity-Bridge/module/config"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := cmd.Execute(rootCmd); err != nil {
		os.Exit(1)
	}
}
