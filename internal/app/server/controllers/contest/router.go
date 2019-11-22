package contest

import (
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/problem"
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("contest")
	{
		h := Handler{}
		r.GET("", h.List())
		r.POST("", RequirePermission("", Admin), h.Create())
		i := r.Group(":contestID", LoadContest())
		{
			i.GET("", h.Show())
			i.PATCH("", RequirePermission("contest", Admin), h.Modify())
			i.DELETE("", RequirePermission("contest", Admin), h.SingleDelete())
			r := i.Group("admin", RequirePermission("contest", Admin))
			{
				r.GET("", h.AdminList())
				r.POST("", h.AdminAdd())
				r.DELETE("", h.AdminDelete())
			}
			r = i.Group("team")
			{
				r.GET("", h.TeamList())
				r.POST("", h.TeamAdd())
				r.DELETE("", RequirePermission("contest", Admin), h.TeamDelete())
				r.GET("__current__", h.TeamCurrentShow())
				r.DELETE("__current__", h.TeamCurrentLeave())
			}
			problem.RegisterRouter(i)
			r = i.Group("notification", RequirePermission("contest", Participate))
			{
				r.GET("", h.NotificationList())
				r.POST("", RequirePermission("contest", Admin), h.NotificationCreate())
				r.DELETE(":notificationOrder", LoadNotification(), RequirePermission("contest", Admin), h.NotificationSingleDelete())
			}
			i.GET("statistic", h.StatisticShow())
		}
	}
}
