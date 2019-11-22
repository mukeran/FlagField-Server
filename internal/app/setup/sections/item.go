package sections

import (
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
	"reflect"
)

type Item struct {
	Key, Description string
	Type             reflect.Kind
	Before, After    hooks.HookFunc
}
