package contest

import (
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"strconv"
)

func loadContest(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	str := c.Param("contestID")
	var contest *cpt.Contest
	if str == "__practice__" {
		contest = cpt.GetDefaultContest(tx)
	} else {
		contestID, err := strconv.Atoi(str)
		if err != nil {
			return false
		}
		contest, err = cpt.GetContestByID(tx, uint(contestID))
		if err != nil {
			return false
		}
	}
	c.Set("contest", contest)
	return true
}

func LoadContest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadContest(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}

func loadNotification(c *gin.Context) bool {
	tx := c.MustGet("tx").(*gorm.DB)
	contest := c.MustGet("contest").(*cpt.Contest)
	notificationOrder, err := strconv.Atoi(c.Param("notificationOrder"))
	if err != nil {
		return false
	}
	notification, err := contest.GetNotificationByOrder(tx, notificationOrder)
	if err != nil {
		return false
	}
	c.Set("notification", notification)
	return true
}

func LoadNotification() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !loadNotification(c) {
			c.Set("type", "json")
			c.Set("error", &errors.ErrNotFound)
			c.Abort()
		}
		c.Next()
	}
}
