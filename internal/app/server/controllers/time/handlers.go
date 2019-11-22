package time

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response ReqShow
			err      error
		)
		defer Pack(c, &err, &response)
		response.Time = time.Now()
	}
}
