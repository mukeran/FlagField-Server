package captcha

import (
	"encoding/json"
	"fmt"
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/config"
	"github.com/FlagField/FlagField-Server/internal/pkg/constants"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/FlagField/FlagField-Server/internal/pkg/mail"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

type Handler struct {
}

func (*Handler) Email() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqEmail
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
		redis := GetRedis(c)
		if cpt.GetConfig(tx, "user.register.enableEmailCaptcha") != "1" {
			err = errors.ErrNotFound
			return
		}
		if request.For != "userRegister" { // to be added
			err = errors.ErrInvalidCaptchaFor
			return
		}
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
		if cpt.HasEmail(tx, request.Email) {
			err = errors.ErrDuplicatedEmail
			return
		}
		var lastTime time.Time
		session.Get(redis, fmt.Sprintf("captcha.%s.time", request.For), &lastTime)
		if lastTime.Add(60 * time.Second).After(time.Now()) {
			err = errors.ErrCoolDown
			return
		}
		code := random.RandomString(8, constants.DicLetterNumeric)
		var mailConfig config.MailConfig
		err = json.Unmarshal([]byte(cpt.GetConfig(tx, "system.mail")), &mailConfig)
		if err != nil {
			panic(err)
		}
		m, err := mail.New(mailConfig)
		if err != nil {
			return
		}
		e := m.Mail([]string{request.Email}, "FlagField User Registration", nil, []byte(fmt.Sprintf("Hi!<br>Welcome to FlagField!<br>Here is your captcha: <strong>%s</strong>.<br>This code is valid for 10 minutes.<br>Enjoy hacking the world!<br>", code)))
		err = m.Send(e)
		if err == nil {
			session.Set(redis, fmt.Sprintf("captcha.%s.code", request.For), code)
			session.Set(redis, fmt.Sprintf("captcha.%s.email", request.For), request.Email)
			session.Set(redis, fmt.Sprintf("captcha.%s.time", request.For), time.Now())
		}
	}
}
