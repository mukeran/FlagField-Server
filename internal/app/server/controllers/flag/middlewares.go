package flag

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func loadFlag(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	problem := c.MustGet("problem").(*cpt.Problem)
	flagOrder, err := strconv.Atoi(c.Param("flagOrder"))
	if err != nil {
		return false
	}
	flag, err := problem.GetFlagByOrder(tx, flagOrder)
	if err != nil {
		return false
	}
	c.Set("flag", flag)
	return true
}

func LoadFlag() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadFlag(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
