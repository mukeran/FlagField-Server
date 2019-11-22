package admin

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("admin", RequirePermission("", Admin))
	{
		h := Handler{}
		r.GET("", h.List())
		r.POST("", h.Add())
		r.DELETE("", h.Delete())
	}
}
