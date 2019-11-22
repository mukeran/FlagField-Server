package user

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func loadUser(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	userID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		return false
	}
	user, err := cpt.GetUserByID(tx, uint(userID))
	if err != nil {
		return false
	}
	c.Set("user", user)
	return true
}

func LoadUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadUser(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
