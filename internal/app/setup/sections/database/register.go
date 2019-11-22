package database

import (
	"reflect"

	"github.com/FlagField/FlagField-Server/internal/app/setup/handler"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

func Register(h *handler.Handler) {
	hook := Hooks{}
	s := h.Section("database", "configure and check database connection settings", hooks.Default(), hook.CheckConnection())
	{
		s.Item(KeyType, "Your database's dialect [mysql/sqlite3/postgres/mssql]", reflect.String, hooks.Default(), hook.CheckDialect())
		s.Item(KeyParameter, "Database's connection parameter", reflect.String, hooks.Default(), hook.CheckParameter())
	}
}
