package hint

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func loadHint(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	problem := c.MustGet("problem").(*cpt.Problem)
	hintOrder, err := strconv.Atoi(c.Param("hintOrder"))
	if err != nil {
		return false
	}
	hint, err := problem.GetHintByOrder(tx, hintOrder)
	if err != nil {
		return false
	}
	c.Set("hint", hint)
	return true
}

func LoadHint() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadHint(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
