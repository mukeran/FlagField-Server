package time

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"time"
)

type ReqShow struct {
	Response
	Time time.Time `json:"time"`
}
