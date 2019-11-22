package problem

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func loadProblem(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	session := c.MustGet("session").(*cpt.Session)
	contest := c.MustGet("contest").(*cpt.Contest)
	problemAlias := c.Param("problemAlias")
	problem, err := cpt.GetProblemByAliasAndContestID(tx, problemAlias, contest.ID)
	if err != nil {
		return false
	}
	if problem.IsHidden && !session.User.IsAdmin && !session.User.IsContestAdmin(tx, contest) {
		return false
	}
	c.Set("problem", problem)
	return true
}

func LoadProblem() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadProblem(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
