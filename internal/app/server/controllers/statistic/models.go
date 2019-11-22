package statistic

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"time"
)

type RespShow struct {
	Response
	SetupTime       time.Time `json:"setup_time"`
	UserCount       uint      `json:"registered_users"`
	SubmissionCount uint      `json:"total_submissions"`
	ContestCount    uint      `json:"hosted_contests"`
	Notification    string    `json:"notification"`
}
