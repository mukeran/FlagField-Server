package statistic

import (
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("statistic")
	{
		h := Handler{}
		r.GET("", h.Show())
	}
}
