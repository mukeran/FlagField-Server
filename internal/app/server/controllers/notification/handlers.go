package notification

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/gin-gonic/gin"
	"strconv"
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
		if session.User.IsAdmin && (c.Query("all") == "1" || c.Query("all") == "true") {
			offset, err1 := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 32)
			limit, err2 := strconv.ParseUint(c.DefaultQuery("limit", "100"), 10, 32)
			if err1 != nil || err2 != nil {
				err = errors.ErrInvalidRequest
				return
			}
			notifications := cpt.GetNotifications(tx, uint(offset), uint(limit))
			response.Notifications = BindList(notifications)
		} else {
			offset, err := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 32)
			if err != nil {
				err = errors.ErrInvalidRequest
				return
			}
			notifications := cpt.GetUserNotifications(tx, session.User, uint(offset))
			response.Notifications = BindList(notifications)
		}
	}
}

func (*Handler) New() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqNew
			response RespNew
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		notification := cpt.NewNotification(tx, request.UserID, request.Content, session.UserID)
		response.NotificationID = notification.ID
	}
}

func (*Handler) Delete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  []uint
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		if v := tx.Exec("delete from `"+cpt.TableNotification+"` where `id` in (?)", request); v.Error != nil {
			panic(v.Error)
		}
	}
}

func (*Handler) MarkRead() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		notificationID, err := strconv.ParseUint(c.Param("notificationID"), 10, 32)
		if err != nil {
			err = errors.ErrInvalidRequest
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		notification, err := cpt.GetNotificationByID(tx, uint(notificationID))
		if err != nil {
			return
		}
		if !session.User.IsAdmin && notification.UserID != session.UserID {
			err = errors.ErrNotFound
			return
		}
		notification.IsRead = true
		notification.Update(tx)
	}
}
