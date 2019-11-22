package submission

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("submission")
	{
		h := Handler{}
		r.GET("", RequirePermission("", Admin), h.List())
	}
}
