package admin

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type ReqAdd struct {
	UsersID []uint `json:"users_id"`
}

type ReqDelete struct {
	UsersID []uint `json:"users_id"`
}

type RAdmin struct {
	ID       uint   `json:"user_id"`
	Username string `json:"username"`
}

func BindAdmin(u *cpt.User) *RAdmin {
	return &RAdmin{
		ID:       u.ID,
		Username: u.Username,
	}
}

type RespList struct {
	Response
	Admins []*RAdmin `json:"admins"`
}

func BindList(users []cpt.User) []*RAdmin {
	var out []*RAdmin
	for _, user := range users {
		out = append(out, BindAdmin(&user))
	}
	return out
}

type RespAddDeleteDetail struct {
	Response
	UserID uint `json:"user_id"`
}

type RespAdd struct {
	Response
	Details []*RespAddDeleteDetail `json:"details"`
}

type RespDelete struct {
	Response
	Details []*RespAddDeleteDetail `json:"details"`
}
