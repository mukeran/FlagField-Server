package sections

import (
	"os"

	"github.com/fatih/color"

	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

func AfterAll() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		color.Yellow("Creating install-lock file...\n")
		f, err := os.Create("setup-lock")
		if err != nil {
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Cannot create setup-lock file! Please create a file named \"setup-lock\" in the directory to prevent re-install.\n")
		} else {
			_ = f.Close()
		}
		color.Green("Successfully created install-lock file...\n")
		color.Green("\nSuccessfully configured FlagField-Server!\nRun server through \"dist/server\".\nEnjoy it!\n")
		return hooks.Normal
	}
}
