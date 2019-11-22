package middlewares

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	cfg "github.com/FlagField/FlagField-Server/internal/pkg/config"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

func resumeSession(c *gin.Context, config cfg.SessionConfig, db *gorm.DB, re *redis.Pool) (session *cpt.Session, err error) {
	sessionID, err := c.Cookie(config.Name)
	if err != nil {
		return
	}
	session, err = cpt.GetSessionBySessionID(db, sessionID)
	if err != nil {
		return
	}
	if session.IsExpired() {
		session.Destroy(db, re)
		err = errors.ErrSessionExpired
	}
	return
}

func ResumeSession(config cfg.SessionConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err     error
			session *cpt.Session
		)
		tx := c.MustGet("db").(*gorm.DB).Begin()
		c.Set("tx", tx)
		re := c.MustGet("redis").(*redis.Pool)
		session, err = resumeSession(c, config, tx, re)
		if err != nil {
			session = cpt.NewSession(tx, &cpt.DefaultUser, config.MaxAge)
		}
		c.SetCookie(config.Name, session.SessionID, config.MaxAge, config.Path, config.Domain, config.Secure, config.HttpOnly)
		session.ExpireAt = time.Now().Add(time.Duration(config.MaxAge) * time.Second)
		session.Update(tx)
		if v := tx.Commit(); v.Error != nil {
			panic(v.Error)
		}
		c.Set("tx", nil)
		c.Set("session", session)
		c.Next()
	}
}
