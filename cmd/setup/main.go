package main

import "github.com/FlagField/FlagField-Server/internal/app/setup"

func main() {
	s := setup.New()
	s.Run()
}
