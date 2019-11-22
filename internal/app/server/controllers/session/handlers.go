package session

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	. "github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (*Handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		sessions := cpt.GetSessions(tx)
		response.Sessions = BindList(sessions)
	}
}

func (*Handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqCreate
			response RespCreate
			err      error
		)
		defer Pack(c, &err, &response)
		session := GetSession(c)
		if session.IsLoggedIn() {
			err = ErrLoggedIn
			return
		}
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		user, err := cpt.GetUserByUsername(tx, request.Username)
		if err != nil {
			err = ErrUserNotFound
			return
		}
		if !user.MatchPassword(request.Password) {
			err = ErrWrongPassword
			return
		}
		session.User = user
		session.Update(tx)
		response.Session = BindSession(session)
	}
}

func (*Handler) Destroy() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		re := GetRedis(c)
		session := GetSession(c)
		if !session.IsLoggedIn() {
			err = ErrNotLoggedIn
			return
		}
		session.Destroy(tx, re)
	}
}

func (*Handler) ViewCurrent() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespViewCurrent
			err      error
		)
		defer Pack(c, &err, &response)
		session := GetSession(c)
		if !session.IsLoggedIn() {
			err = ErrNotLoggedIn
			return
		}
		response.Session = BindSession(session)
	}
}
