package middlewares

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func ResumeDefault() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := c.MustGet("db").(*gorm.DB).Begin()
		c.Set("tx", tx)
		if v := tx.Where("name = ?", cpt.DefaultContest.Name).First(&cpt.DefaultContest); v.Error != nil {
			panic(v.Error)
		}
		if v := tx.Where("username = ?", cpt.DefaultUser.Username).First(&cpt.DefaultUser); v.Error != nil {
			panic(v.Error)
		}
		if v := tx.Commit(); v.Error != nil {
			panic(v.Error)
		}
		c.Set("tx", nil)
		c.Next()
	}
}
