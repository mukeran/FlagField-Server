package migrator

import (
	"fmt"
	"github.com/FlagField/FlagField-Server/internal/app/migrator/templates"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

const (
	VersionInitial = "initial"
)

type Migrator struct {
	_version string
	_gorm    *gorm.DB
}

func (m *Migrator) Execute() {
	var err error
	switch m._version {
	case VersionInitial:
		err = (&templates.Initial{}).Execute(m._gorm)
	}
	buildOutput(err)
}

func New(conf *config.Config, version string) (*Migrator, error) {
	db, err := gorm.Open(conf.Database.Type, conf.Database.Parameter)
	if err != nil {
		return nil, err
	}
	return &Migrator{
		_version: version,
		_gorm:    db,
	}, nil
}

func buildOutput(err error) {
	if err != nil {
		fmt.Println("Migrate failed!")
		fmt.Printf("error: %v\n", err)
	} else {
		fmt.Println("Migrate successful!")
	}
}

type Template interface {
	Execute(db *gorm.DB) error
}
