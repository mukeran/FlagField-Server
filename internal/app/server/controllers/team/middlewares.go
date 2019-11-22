package team

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func loadTeam(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	str := c.Param("teamID")
	teamID, err := strconv.Atoi(str)
	if err != nil {
		return false
	}
	team, err := cpt.GetTeamByID(tx, uint(teamID))
	if err != nil {
		return false
	}
	c.Set("team", team)
	return true
}

func LoadTeam() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadTeam(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
