package server

import (
	"reflect"

	"github.com/FlagField/FlagField-Server/internal/app/setup/handler"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

func Register(h *handler.Handler) {
	hook := Hooks{}
	s := h.DefaultSection("server", "basic configure of server")
	{
		s.Item(KeyPort, "the listening port", reflect.Uint, hooks.Default(), hook.CheckPort())
	}
}
