package config

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (*Handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespAll
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		configs := cpt.GetConfigs(tx)
		response.Config = BindList(configs)
	}
}

func (*Handler) Get() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespGet
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		key := c.Param("configKey")
		response.Value = cpt.GetConfig(tx, key)
	}
}

func (*Handler) Set() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  map[string]string
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		for key, value := range request {
			cpt.SetConfig(tx, key, value)
		}
	}
}
