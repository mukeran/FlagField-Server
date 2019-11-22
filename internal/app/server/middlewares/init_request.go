package middlewares

import (
	"github.com/FlagField/FlagField-Server/internal/app/server/instance"
	"github.com/gin-gonic/gin"
)

func InitRequest(s *instance.Instance) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("config", s.GetConfig())
		c.Set("db", s.GetDB())
		c.Set("redis", s.GetRedis())
		c.Next()
	}
}
