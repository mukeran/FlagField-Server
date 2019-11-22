package problem

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	. "github.com/FlagField/FlagField-Server/internal/pkg/errors"
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
		if !session.User.IsAdmin && !contest.IsAdmin(tx, session.User) && contest.StartTime.After(time.Now()) {
			err = ErrContestPending
			return
		}
		problems := contest.GetProblems(tx)
		response.Problems = BindList(tx, session, session.User.IsContestAdmin(tx, contest), problems)
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
		session := GetSession(c)
		contest := GetContest(c)
		if cpt.HasProblemAlias(tx, contest.ID, request.Alias) {
			err = ErrDuplicatedAlias
			return
		}
		problem := cpt.NewProblem(tx, request.Name, request.Description, request.Alias, contest.ID, request.Points, strings.ToLower(request.Type), session.UserID)
		response.ProblemAlias = problem.Alias
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
		if !session.User.IsAdmin && !contest.IsAdmin(tx, session.User) && contest.StartTime.After(time.Now()) {
			err = ErrContestPending
			return
		}
		problem := GetProblem(c)
		response.Problem = &RProblem{
			Name:        problem.Name,
			Description: problem.Description,
			Alias:       problem.Alias,
			ContestID:   problem.ContestID,
			Points:      problem.Points,
			Type:        problem.Type,
			Tags:        problem.GetTagsName(tx),
			IsHidden:    problem.IsHidden,
		}
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
		problem := GetProblem(c)
		request.Bind(problem)
		problem.Update(tx)
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
		problem := GetProblem(c)
		problem.Delete(tx)
	}
}

func (*Handler) TagAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tags     []string
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = c.ShouldBindJSON(&tags)
		if err != nil {
			err = ErrInvalidRequest
			return
		}
		tx := GetDB(c)
		problem := GetProblem(c)
		for _, tag := range tags {
			problem.AddTag(tx, tag)
		}
	}
}

func (*Handler) TagDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			tags     []string
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = c.ShouldBindJSON(&tags)
		if err != nil {
			err = ErrInvalidRequest
			return
		}
		tx := GetDB(c)
		problem := GetProblem(c)
		for _, tag := range tags {
			problem.DeleteTag(tx, tag)
		}
	}
}

func (*Handler) SubmissionCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqSubmissionCreate
			response RespSubmissionCreate
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		contest := GetContest(c)
		if contest.StartTime.After(time.Now()) {
			err = ErrContestPending
			return
		}
		if contest.EndTime.Before(time.Now()) {
			err = ErrContestEnded
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		problem := GetProblem(c)
		submission := cpt.NewSubmission(tx, request.Flag, problem, session.User)
		response.Result = submission.Result
	}
}
