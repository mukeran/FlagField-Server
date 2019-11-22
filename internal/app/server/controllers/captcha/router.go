package captcha

import "github.com/gin-gonic/gin"

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("captcha")
	h := Handler{}
	{
		r.POST("email", h.Email())
	}
}
