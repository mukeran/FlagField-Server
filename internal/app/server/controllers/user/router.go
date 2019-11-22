package user

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("user")
	{
		h := Handler{}
		r.GET("", h.List())
		r.POST("", h.Register())
		i := r.Group(":userID", LoadUser())
		{
			i.GET("", h.Show())
			i.DELETE("", RequirePermission("user", Admin), h.SingleDelete())
			i.PUT("password", RequirePermission("user", Admin), h.PasswordModify())
			i.PUT("email", RequirePermission("", Admin), h.EmailModify())
			i.PATCH("profile", RequirePermission("user", Admin), h.ProfileModify())
			r = i.Group("team")
			{
				r.GET("", RequirePermission("user", Admin), h.TeamList())
				r.GET("invitation", RequirePermission("user", Admin), h.TeamInvitationList())
				r.GET("application", RequirePermission("user", Admin), h.TeamApplicationList())
			}
		}
	}
}
