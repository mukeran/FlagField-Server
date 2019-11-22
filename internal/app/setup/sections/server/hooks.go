package server

import (
	"os"

	"github.com/fatih/color"

	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

type Hooks struct {
}

func (*Hooks) CheckPort() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		port := (*m)[KeyPort].(uint64)
		if !(port >= 1 && port <= 65535) {
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Invalid port value, please check and reinput.")
			return hooks.BeginItem
		}
		return hooks.Normal
	}
}
