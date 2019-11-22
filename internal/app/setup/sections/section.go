package sections

import (
	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
	"reflect"
)

type Section struct {
	Name, Description string
	Items             []*Item
	Before, After     hooks.HookFunc
}

func (s *Section) Item(_key, _description string, _type reflect.Kind, _before, _after hooks.HookFunc) *Item {
	item := &Item{
		Key:         _key,
		Description: _description,
		Type:        _type,
		Before:      _before,
		After:       _after,
	}
	s.Items = append(s.Items, item)
	return item
}
