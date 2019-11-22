package middlewares

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	User = iota
	Participate
	Admin
)

func checkPermission(c *gin.Context, resource string, access uint) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	session := c.MustGet("session").(*cpt.Session)
	if !session.IsLoggedIn() {
		return false
	}
	if access == User || session.User.IsAdmin {
		return true
	}
	switch access {
	case Participate:
		switch resource {
		case "contest":
			contest := c.MustGet("contest").(*cpt.Contest)
			if session.User.HasContestAccess(tx, contest) {
				return true
			}
		case "resource":
			resource := c.MustGet("resource").(*cpt.Resource)
			if session.User.HasResourceAccess(tx, resource) {
				return true
			}
		case "team":
			team := c.MustGet("team").(*cpt.Team)
			if session.User.HasTeamAccess(tx, team) {
				return true
			}
		}
	case Admin:
		switch resource {
		case "contest":
			contest := c.MustGet("contest").(*cpt.Contest)
			if session.User.IsContestAdmin(tx, contest) {
				return true
			}
		case "resource":
			resource := c.MustGet("resource").(*cpt.Resource)
			if session.User.IsResourceAdmin(resource) {
				return true
			}
		case "team":
			team := c.MustGet("team").(*cpt.Team)
			if session.User.IsTeamAdmin(tx, team) {
				return true
			}
		case "user":
			user := c.MustGet("user").(*cpt.User)
			if session.UserID == user.ID {
				return true
			}
		}
	}
	return false
}

func RequirePermission(resource string, access uint) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkPermission(c, resource, access) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
