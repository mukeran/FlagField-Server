package resource

import (
	"reflect"

	"github.com/FlagField/FlagField-Server/internal/app/setup/handler"
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

func Register(h *handler.Handler) {
	hook := Hooks{}
	s := h.DefaultSection("resource", "set the path to store resource")
	{
		s.Item(KeyBaseDir, "The base directory of uploaded resource", reflect.String, hooks.Default(), hook.CheckBaseDir())
	}
}
