package middlewares

import (
	"fmt"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				obj, ok := c.Get("tx")
				if ok && obj != nil {
					tx := obj.(*gorm.DB)
					tx.Rollback()
				}
				c.AbortWithStatusJSON(controllers.RespServerError.GetHttpStatus(), controllers.RespServerError)
				fmt.Println(err)
				// log here
			}
		}()
		c.Next()
	}
}
