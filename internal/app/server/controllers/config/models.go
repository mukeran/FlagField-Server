package config

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type RespAll struct {
	Response
	Config map[string]string `json:"config"`
}

func BindList(configs []cpt.Config) map[string]string {
	out := make(map[string]string)
	for _, config := range configs {
		out[config.Key] = config.Value
	}
	return out
}

type RespGet struct {
	Response
	Value string `json:"value"`
}
