package time

import "github.com/gin-gonic/gin"

func RegisterRouter(g *gin.RouterGroup) {
	h := Handler{}
	g.GET("time", h.Show())
}
