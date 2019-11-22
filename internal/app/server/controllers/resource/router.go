package resource

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("resource")
	{
		h := Handler{}
		r.GET("", h.List())    // Internal permission check
		r.POST("", h.Upload()) // Internal permission check
		rirg := r.Group(":resourceUUID", LoadResource())
		{
			rirg.GET("", RequirePermission("resource", Participate), h.Download())
			rirg.PATCH("", RequirePermission("resource", Admin), h.Modify())
			rirg.DELETE("", RequirePermission("resource", Admin), h.SingleDelete())
		}
	}
}
