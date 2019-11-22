package redis

import (
	"os"

	"github.com/fatih/color"
	"github.com/gomodule/redigo/redis"

	"github.com/FlagField/FlagField-Server/internal/app/setup/hooks"
)

type Hooks struct {
}

func (*Hooks) CheckRedisConnection() hooks.HookFunc {
	return func(m *hooks.Map) uint {
		color.Yellow("Checking your redis config...\n")
		uri := (*m)[KeyURI].(string)
		password := (*m)[KeyPassword].(string)
		db := int((*m)[KeyDB].(uint64))
		r, err := redis.Dial("tcp", uri,
			redis.DialPassword(password),
			redis.DialDatabase(db),
		)
		if err != nil {
			red := color.Set(color.FgHiRed)
			_, _ = red.Fprintf(os.Stderr, "Cannot connect to redis.\nError message: (%v)\nPlease check and reinput.\n", err)
			return hooks.BeginSection
		}
		defer r.Close()
		color.Green("Successfully connected to redis\n")
		return hooks.Normal
	}
}
