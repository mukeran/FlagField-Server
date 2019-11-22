package admin

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
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
		admins := cpt.GetAdmins(tx)
		response.Admins = BindList(admins)
	}
}

func (*Handler) Add() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqAdd
			response RespAdd
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		for _, userID := range request.UsersID {
			var (
				r2   RespAddDeleteDetail
				err2 error
			)
			user, err2 := cpt.GetUserByID(tx, userID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
				goto end
			}
			user.IsAdmin = true
			user.Update(tx)
		end:
			r2.UserID = userID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
	}
}

func (*Handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqDelete
			response RespDelete
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		for _, userID := range request.UsersID {
			var (
				r2   RespAddDeleteDetail
				err2 error
			)
			user, err2 := cpt.GetUserByID(tx, userID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
				goto end
			}
			user.IsAdmin = false
			user.Update(tx)
		end:
			r2.UserID = userID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
	}
}
