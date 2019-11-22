package notification

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("notification")
	h := Handler{}
	{
		r.GET("", RequirePermission("", User), h.List())
		r.POST("", RequirePermission("", Admin), h.New())
		r.DELETE("", RequirePermission("", Admin), h.Delete())
		r.PATCH(":notificationID", RequirePermission("", User), h.MarkRead())
	}
}
