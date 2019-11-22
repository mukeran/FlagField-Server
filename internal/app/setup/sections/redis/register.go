package redis

import (
	"reflect"

	"github.com/FlagField/FlagField-Server/internal/app/setup/handler"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

func Register(h *handler.Handler) {
	hook := Hooks{}
	s := h.Section("redis", "configure and check redis connection settings", hooks.Default(), hook.CheckRedisConnection())
	{
		s.Item(KeyURI, "URI of redis (host:port)", reflect.String, hooks.Default(), hooks.Default())
		s.Item(KeyPassword, "Password of redis (normally empty)", reflect.String, hooks.Default(), hooks.Default())
		s.Item(KeyDB, "DB of redis (normally 0)", reflect.Uint, hooks.Default(), hooks.Default())
	}
}
