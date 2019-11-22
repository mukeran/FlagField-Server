package team

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("team")
	{
		h := Handler{}
		r.GET("", h.List())
		r.POST("", RequirePermission("", User), h.Create())
		i := r.Group(":teamID", LoadTeam())
		{
			i.GET("", h.Show())
			i.PATCH("", RequirePermission("team", Admin), h.Modify())
			i.DELETE("", RequirePermission("team", Admin), h.SingleDelete())
			r := i.Group("admin", RequirePermission("team", Admin))
			{
				r.GET("", h.AdminList())
				r.POST("", h.AdminAdd())
				r.DELETE("", h.AdminDelete())
			}
			r = i.Group("user")
			{
				r.GET("", h.UserList())
				r.POST("", RequirePermission("", Admin), h.UserAdd())
				r.DELETE("", RequirePermission("team", Admin), h.UserDelete())
			}
			i.GET("statistic", h.StatisticShow())
			r = i.Group("invitation")
			{
				r.GET("", RequirePermission("team", Admin), h.InvitationList())
				r.POST("", RequirePermission("team", Admin), h.InvitationNew())
				r.DELETE("", RequirePermission("team", Admin), h.InvitationCancel())
				r.GET("accept", RequirePermission("", User), h.InvitationAccept())
				r.POST("accept", RequirePermission("", User), h.InvitationAcceptByToken())
				r.GET("reject", RequirePermission("", User), h.InvitationReject())
				r.GET("token", RequirePermission("team", Admin), h.InvitationTokenShow())
				r.DELETE("token", RequirePermission("", Admin), h.InvitationTokenRefresh())
			}
			r = i.Group("application")
			{
				r.GET("", RequirePermission("team", Admin), h.ApplicationList())
				r.POST("", RequirePermission("", User), h.ApplicationNew())
				r.DELETE("", RequirePermission("", User), h.ApplicationCancel())
				r.POST("accept", RequirePermission("team", Admin), h.ApplicationAccept())
				r.POST("reject", RequirePermission("team", Admin), h.ApplicationReject())
			}
		}
	}
}
