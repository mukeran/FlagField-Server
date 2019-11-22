package statistic

import (
	"time"

	"github.com/gin-gonic/gin"

	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type Handler struct {
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespShow
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		response.SetupTime, _ = time.Parse(time.RFC3339, cpt.GetConfig(tx, "system.setup_time"))
		response.UserCount = cpt.GetUserCount(tx)
		response.SubmissionCount = cpt.GetSubmissionCount(tx)
		response.ContestCount = cpt.GetContestCount(tx)
		response.Notification = cpt.GetConfig(tx, "index.notification")
	}
}
