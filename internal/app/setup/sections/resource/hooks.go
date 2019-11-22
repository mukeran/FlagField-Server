package resource

import (
	"os"

	"github.com/fatih/color"

	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

type Hooks struct {
}

func (*Hooks) CheckBaseDir() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		baseDir := (*m)[KeyBaseDir].(string)
		_, err := os.Stat(baseDir)
		if err != nil {
			if os.IsNotExist(err) {
				err = os.Mkdir(baseDir, os.ModePerm)
				if err != nil {
					red := color.Set(color.FgHiRed)
					_, _ = red.Fprintf(os.Stderr, "Error when creating directory, please create this directory by your self or specify another directory.\n")
					return hooks.BeginItem
				}
				color.Green("Successfully created upload directory.\n")
				return hooks.Normal
			}
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Error when checking directory (%v), please check or specify another directory.\n", err)
			return hooks.BeginItem
		}
		color.Green("Using existed directory.\n")
		return hooks.Normal
	}
}
