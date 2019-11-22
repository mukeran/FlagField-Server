package middlewares

import (
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"github.com/gin-gonic/gin"
)

func Handler404() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.AbortWithStatusJSON(controllers.RespNotFound.GetHttpStatus(), controllers.RespNotFound)
	}
}
