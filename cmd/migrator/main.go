package main

import (
	"flag"
	"github.com/FlagField/FlagField-Server/internal/app/migrator"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
)

func main() {
	template := flag.String("template", "initial", "The template that migration applies")
	flag.Parse()
	m, err := migrator.New(config.FromFile("./config.json"), *template)
	if err != nil {
		panic(err)
	}
	m.Execute()
}
