package flag

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("flag")
	{
		h := Handler{}
		r.GET("", RequirePermission("contest", Admin), h.List())
		r.POST("", RequirePermission("contest", Admin), h.Create())
		r.DELETE("", RequirePermission("contest", Admin), h.BatchDelete())
		i := r.Group(":flagOrder", LoadFlag(), RequirePermission("contest", Admin))
		{
			i.GET("", h.Show())
			i.PATCH("", h.Modify())
			i.DELETE("", h.SingleDelete())
		}
	}
}
