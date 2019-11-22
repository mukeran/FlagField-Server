package hint

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("hint")
	{
		h := Handler{}
		r.GET("", h.List())
		r.POST("", RequirePermission("contest", Admin), h.Create())
		i := r.Group(":hintOrder", LoadHint())
		{
			i.GET("", h.Show())
			i.PATCH("", RequirePermission("contest", Admin), h.Modify())
			i.DELETE("", RequirePermission("contest", Admin), h.SingleDelete())
		}
	}
}
