package main

import (
	"github.com/FlagField/FlagField-Server/internal/app/server"
	"os"

	"github.com/FlagField/FlagField-Server/internal/pkg/config"
)

func main() {
	instance := server.New()
	workDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	err = os.Setenv("FLAGFIELD_HOME", workDir)
	if err != nil {
		panic(err)
	}
	conf := config.FromFile(workDir + "/config.json")
	instance.SetConfig(*conf)
	instance.Run()
}
