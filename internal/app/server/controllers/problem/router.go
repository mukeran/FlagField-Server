package problem

import (
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/flag"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/hint"
	. "github.com/FlagField/FlagField-Server/internal/app/server/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRouter(g *gin.RouterGroup) {
	r := g.Group("problem", RequirePermission("contest", Participate))
	{
		h := Handler{}
		r.GET("", h.List())
		r.POST("", RequirePermission("contest", Admin), h.Create())
		i := r.Group(":problemAlias", LoadProblem())
		{
			i.GET("", h.Show())
			i.PATCH("", RequirePermission("contest", Admin), h.Modify())
			i.DELETE("", RequirePermission("contest", Admin), h.SingleDelete())
			flag.RegisterRouter(i)
			hint.RegisterRouter(i)
			i.POST("submission", h.SubmissionCreate())
			r := i.Group("tag")
			{
				r.POST("", RequirePermission("contest", Admin), h.TagAdd())
				r.DELETE("", RequirePermission("contest", Admin), h.TagDelete())
			}
		}
	}
}
