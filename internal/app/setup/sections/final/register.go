package final

import (
	"reflect"

	"github.com/FlagField/FlagField-Server/internal/app/setup/handler"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

func Register(h *handler.Handler) {
	hook := Hooks{}
	s := h.Section("final", "Final settings", hook.Before(), hook.After())
	{
		s.Item(KeyUsername, "System admin username", reflect.String, hooks.Default(), hooks.Default())
		s.Item(KeyPassword, "System admin password", reflect.String, hooks.Default(), hooks.Default())
		s.Item(KeyEmail, "System admin email", reflect.String, hooks.Default(), hooks.Default())
	}
}
