package config

import (
	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"
)

func Register(db *gorm.DB) *cli.Command {
	return &cli.Command{
		Name:     "config",
		Aliases:  []string{"c"},
		Usage:    "manage configures",
		Category: "resources",
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "list configures",
				Action:  list(db),
			},
			{
				Name:    "set",
				Aliases: []string{"s"},
				Usage:   "set configures",
				Action:  set(db),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "from-file",
						Aliases: []string{"f", "file"},
						Usage:   "add user from file (json, xml, yaml, xls, xlsx)",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "key",
						Aliases: []string{"k"},
						Usage:   "the config key",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "value",
						Aliases: []string{"v", "val"},
						Usage:   "the config value",
						Value:   "",
					},
				},
			},
			{
				Name:    "get",
				Aliases: []string{"g"},
				Usage:   "get configures",
				Action:  get(db),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "the config key",
						Required: true,
						Value:    "",
					},
				},
			},
		},
	}
}
