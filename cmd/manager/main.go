package main

import (
	"github.com/FlagField/FlagField-Server/internal/app/manager"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
)

func main() {
	m, err := manager.New(config.FromFile("./config.json"))
	if err != nil {
		panic(err)
	}
	m.Run()
}
