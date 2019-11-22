package middlewares

import (
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Transaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := c.MustGet("db").(*gorm.DB).Begin()
		c.Set("tx", tx)

		c.Next()

		if v := tx.Commit(); v.Error != nil {
			panic(v.Error)
		}
		if c.MustGet("type").(string) == "json" {
			var (
				obj interface{}
				ok  bool
				r   controllers.ResponseInterface
				err *error
			)
			obj, ok = c.Get("response")
			if ok {
				r = obj.(controllers.ResponseInterface)
			} else {
				r = &controllers.Response{}
			}
			obj, ok = c.Get("error")
			if ok {
				err = obj.(*error)
			} else {
				err = nil
			}
			r.Pack(err)
			c.JSON(r.GetHttpStatus(), r)
		}
	}
}
