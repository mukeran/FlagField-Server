package config

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("config", RequirePermission("", Admin))
	{
		h := Handler{}
		r.GET("", h.List())
		r.PATCH("", h.Set())
		r.GET(":configKey", h.Get())
	}
}
