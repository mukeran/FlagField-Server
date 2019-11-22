package manager

import (
	configCommand "github.com/FlagField/FlagField-Server/internal/app/manager/commands/config"
	"github.com/FlagField/FlagField-Server/internal/app/manager/commands/user"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/urfave/cli/v2"
	"os"
)

type Manger struct {
	_gorm *gorm.DB
	app   *cli.App
}

func (m *Manger) Run() {
	m.app = &cli.App{
		Name:     "FlagField-Server Manager",
		HelpName: "manager",
		Usage:    "manage server resources",
		Version:  "0.0.1",
		Commands: []*cli.Command{
			configCommand.Register(m._gorm),
			user.Register(m._gorm),
		},
		EnableBashCompletion: true,
	}
	err := m.app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}

func New(conf *config.Config) (*Manger, error) {
	db, err := gorm.Open(conf.Database.Type, conf.Database.Parameter)
	if err != nil {
		return nil, err
	}
	db.SingularTable(true)
	if v := db.Where("name = ?", cpt.DefaultContest.Name).First(&cpt.DefaultContest); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Where("username = ?", cpt.DefaultUser.Username).First(&cpt.DefaultUser); v.Error != nil {
		panic(v.Error)
	}
	return &Manger{_gorm: db}, nil
}
