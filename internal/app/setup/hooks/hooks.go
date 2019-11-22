package hooks

const (
	Normal = iota + 1
	Begin
	End
	BeginSection
	EndSection
	NextSection
	BeginItem
	EndItem
	NextItem
)

type Map map[string]interface{}

type HookFunc func(mapping *Map) uint

func Default() HookFunc {
	return func(mapping *Map) uint {
		return Normal
	}
}
