package database

import (
	"os"

	"github.com/fatih/color"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	// _ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

type Hooks struct {
}

func (*Hooks) CheckDialect() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		if (*m)[KeyType] != "mysql" && (*m)[KeyType] != "sqlite3" && (*m)[KeyType] != "postgres" && (*m)[KeyType] != "mssql" {
			red := color.New(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Wrong type! Please check and reinput.\n")
			return hooks.BeginItem
		}
		return hooks.Normal
	}
}

func (*Hooks) CheckParameter() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		return hooks.Normal
	}
}

func (*Hooks) CheckConnection() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		color.Yellow("Checking your database config...\n")
		dialect := (*m)[KeyType].(string)
		parameter := (*m)[KeyParameter].(string)
		db, err := gorm.Open(dialect, parameter)
		if err != nil {
			red := color.New(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Cannot connect to the database, please check your configure.\n")
			return hooks.BeginSection
		}
		color.Green("Successfully connected to database\n")
		defer db.Close()
		return hooks.Normal
	}
}
