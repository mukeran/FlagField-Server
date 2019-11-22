package hint

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
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
		session := GetSession(c)
		contest := GetContest(c)
		problem := GetProblem(c)
		hints := problem.GetHints(tx)
		response.Hints, err = BindList(tx, contest, problem, session, hints)
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
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		problem := GetProblem(c)
		response.HintOrder = problem.AddHint(tx, request.Cost, request.Content)
	}
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespShow
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		session := GetSession(c)
		contest := GetContest(c)
		problem := GetProblem(c)
		hint := GetHint(c)
		if err != nil {
			return
		}
		if !hint.IsUnlocked(tx, contest, problem, session.User) {
			err = hint.Unlock(tx, contest, problem, session.User)
			if err != nil {
				return
			}
		}
		response.Hint = BindHint(hint, true)
	}
}

func (*Handler) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqModify
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		hint := GetHint(c)
		request.Bind(hint)
		hint.Update(tx)
	}
}

func (*Handler) SingleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		hint := GetHint(c)
		hint.Delete(tx)
	}
}
