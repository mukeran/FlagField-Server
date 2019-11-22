package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	cfg "github.com/FlagField/FlagField-Server/internal/pkg/config"
	. "github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/FlagField/FlagField-Server/internal/pkg/validator"
)

func GetSession(c *gin.Context) *cpt.Session {
	return c.MustGet("session").(*cpt.Session)
}

func GetDB(c *gin.Context) *gorm.DB {
	return c.MustGet("tx").(*gorm.DB)
}

func GetRedis(c *gin.Context) *redis.Pool {
	return c.MustGet("redis").(*redis.Pool)
}

func GetConfig(c *gin.Context) *cfg.Config {
	return c.MustGet("config").(*cfg.Config)
}

func GetContest(c *gin.Context) *cpt.Contest {
	return c.MustGet("contest").(*cpt.Contest)
}

func GetProblem(c *gin.Context) *cpt.Problem {
	return c.MustGet("problem").(*cpt.Problem)
}

func GetFlag(c *gin.Context) *cpt.Flag {
	return c.MustGet("flag").(*cpt.Flag)
}

func GetHint(c *gin.Context) *cpt.Hint {
	return c.MustGet("hint").(*cpt.Hint)
}

func GetNotification(c *gin.Context) *cpt.ContestNotification {
	return c.MustGet("notification").(*cpt.ContestNotification)
}

func GetResource(c *gin.Context) *cpt.Resource {
	return c.MustGet("resource").(*cpt.Resource)
}

func GetTeam(c *gin.Context) *cpt.Team {
	return c.MustGet("team").(*cpt.Team)
}

func GetUser(c *gin.Context) *cpt.User {
	return c.MustGet("user").(*cpt.User)
}

func Pack(c *gin.Context, err *error, response ResponseInterface) {
	c.Set("type", "json")
	c.Set("error", err)
	c.Set("response", response)
}

func BindJSON(c *gin.Context, r interface{}) (err error) {
	err = c.ShouldBindJSON(r)
	if err != nil && !validator.IsValidationError(err) {
		return ErrInvalidRequest
	}
	return
}
