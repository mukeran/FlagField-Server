package user

import (
	"github.com/jinzhu/gorm"
	"github.com/urfave/cli/v2"
)

func Register(db *gorm.DB) *cli.Command {
	return &cli.Command{
		Name:     "user",
		Aliases:  []string{"u"},
		Usage:    "manage users",
		Category: "resources",
		Subcommands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"l", "ls"},
				Usage:   "list users",
				Action:  list(db),
				Flags: []cli.Flag{
					&cli.UintFlag{
						Name:    "offset",
						Aliases: []string{"o"},
						Usage:   "set the offset of list",
						Value:   0,
					},
					&cli.UintFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "set the limit of list",
						Value:   10,
					},
					&cli.StringFlag{
						Name:    "query",
						Aliases: []string{"q"},
						Usage:   "query parameters",
						Value:   "",
					},
				},
			},
			{
				Name:    "add",
				Aliases: []string{"a", "new"},
				Usage:   "add users",
				Action:  add(db),
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "from-file",
						Aliases: []string{"f", "file"},
						Usage:   "add user from file (json, xml, yaml, xls, xlsx)",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "username",
						Aliases: []string{"u", "user"},
						Usage:   "input the username",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "password",
						Aliases: []string{"p", "pass"},
						Usage:   "input the password. (empty means to randomly generate one)",
						Value:   "",
					},
					&cli.StringFlag{
						Name:    "email",
						Aliases: []string{"e"},
						Usage:   "input the email",
						Value:   "",
					},
					&cli.BoolFlag{Name: "admin"},
				},
			},
			{
				Name:    "delete",
				Aliases: []string{"del", "rm", "remove"},
				Usage:   "delete users",
				Action:  del(db),
			},
		},
	}
}
