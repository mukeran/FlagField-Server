package user

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
	"time"
)

type Handler struct {
}

// list godoc
// @Summary Show user list
// @Produce json
// @Param order_by query string false "OrderBy"
// @Param page query int false "Page"
// @Success 200 {object} response.List
// @Failure 400 {object} response.List
// @Failure 500 {object} response.List
// @Router /v1/user/ [GET]
func (*Handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		session := GetSession(c)
		users := cpt.GetUsers(tx)
		response.Users = BindList(session, users)
	}
}

// register godoc
// @Summary Register a account
// @Produce json
// @Param register body request.Register true "Register"
// @Success 200 {object} response.Register
// @Failure 400 {object} response.Register
// @Failure 500 {object} response.Register
// @Router /v1/user/ [POST]
func (*Handler) Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqRegister
			response RespRegister
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		redis := GetRedis(c)
		if cpt.GetConfig(tx, "user.register.enableEmailWhitelist") == "1" {
			pos := strings.IndexByte(request.Email, byte('@'))
			emailWhitelist := strings.Split(cpt.GetConfig(tx, "user.register.emailWhitelist"), ",")
			flag := false
			requestDomain := request.Email[pos+1:]
			for _, domain := range emailWhitelist {
				if requestDomain == domain {
					flag = true
					break
				}
			}
			if !flag {
				err = errors.ErrNotInWhitelist
				return
			}
		}
		flagEmailCaptcha := !session.User.IsAdmin && cpt.GetConfig(tx, "user.register.enableEmailCaptcha") == "1"
		if flagEmailCaptcha {
			var (
				code            string
				email           string
				codeRequestTime time.Time
			)
			session.Get(redis, "captcha.userRegister.code", &code)
			session.Get(redis, "captcha.userRegister.email", &email)
			session.Get(redis, "captcha.userRegister.time", &codeRequestTime)
			if codeRequestTime.Add(10*time.Minute).Before(time.Now()) || code != request.EmailCaptcha || email != request.Email {
				err = errors.ErrInvalidCaptcha
				return
			}
		}
		if cpt.HasEmail(tx, request.Email) {
			err = errors.ErrDuplicatedEmail
			return
		}
		user, err := cpt.NewUser(tx, request.Username, request.Password)
		if err != nil {
			return
		}
		user.Email = request.Email
		user.Update(tx)
		response.UserID = user.ID
		if flagEmailCaptcha {
			session.Del(redis, "captcha.userRegister.code")
			session.Del(redis, "captcha.userRegister.email")
		}
	}
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespShow
			err      error
		)
		defer Pack(c, &err, &response)
		session := GetSession(c)
		user := GetUser(c)
		response.User = BindUser(session, user)
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
		user := GetUser(c)
		user.Delete(tx)
	}
}

func (*Handler) PasswordModify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqPasswordModify
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		user := GetUser(c)
		if !session.User.IsAdmin && !user.MatchPassword(request.OldPassword) {
			err = errors.ErrWrongPassword
			return
		}
		user.ChangePassword(tx, request.NewPassword)
	}
}

func (*Handler) EmailModify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqEmailModify
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		if cpt.HasEmail(tx, request.Email) {
			err = errors.ErrDuplicatedEmail
			return
		}
		user := GetUser(c)
		user.Email = request.Email
		user.Update(tx)
	}
}

func (*Handler) ProfileModify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqProfileModify
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		user := GetUser(c)
		request.Bind(user.Profile)
		user.Update(tx)
	}
}

func (*Handler) TeamList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespTeamList
			err      error
		)
		defer Pack(c, &err, &response)
		admin, err := strconv.ParseBool(c.DefaultQuery("admin", "0"))
		query := c.DefaultQuery("query", "")
		if err != nil {
			err = errors.ErrInvalidRequest
			return
		}
		tx := GetDB(c)
		user := GetUser(c)
		teams := cpt.GetUserTeams(tx, user.ID, admin, query)
		response.Teams = BindTeamList(tx, teams, user)
	}
}

func (*Handler) TeamInvitationList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespTeamInvitationList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		user := GetUser(c)
		invitations := cpt.GetUserInvitations(tx, user.ID)
		response.Invitations = BindTeamInvitationList(invitations)
	}
}

func (*Handler) TeamApplicationList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespTeamApplicationList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		user := GetUser(c)
		applications := cpt.GetUserApplications(tx, user.ID)
		response.Applications = BindTeamApplicationList(applications)
	}
}
