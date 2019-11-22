package session

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("session")
	{
		h := Handler{}
		r.GET("", RequirePermission("", Admin), h.List())
		r.POST("", h.Create())
		r.DELETE("", h.Destroy())
		r.GET("__current__", h.ViewCurrent())
	}
}
