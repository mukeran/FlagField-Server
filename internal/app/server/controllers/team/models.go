package team

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
	"github.com/jinzhu/gorm"
	"time"
)

type ReqCreate struct {
	Name        string `json:"name" validate:"max=50,min=1,keyword" binding:"max=50,min=1,keyword"`
	Description string `json:"description"`
}

type ReqModify struct {
	Name        *string `json:"name" validate:"omitempty,max=50,min=1,keyword" binding:"omitempty,max=50,min=1,keyword"`
	Description *string `json:"description"`
}

func (r *ReqModify) Bind(t *cpt.Team) {
	if r.Name != nil {
		t.Name = *r.Name
	}
	if r.Description != nil {
		t.Description = *r.Description
	}
}

type ReqUserAdd struct {
	UsersID []uint `json:"users_id"`
}

type ReqUserDelete struct {
	UsersID []uint `json:"users_id"`
}

type ReqInvitationNewAndCancel struct {
	UserID uint `json:"user_id" validate:"required" binding:"required"`
}

type ReqInvitationAcceptByToken struct {
	Token string `json:"token"`
}

type ReqApplicationAcceptAndReject struct {
	UserID uint `json:"user_id" validate:"required" binding:"required"`
}

type RTeam struct {
	TeamID      uint      `json:"team_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MemberCount uint      `json:"member_count"`
	CreateAt    time.Time `json:"create_at"`
}

func BindTeam(db *gorm.DB, t *cpt.Team) *RTeam {
	return &RTeam{
		TeamID:      t.ID,
		Name:        t.Name,
		Description: t.Description,
		MemberCount: t.GetMembersCount(db),
		CreateAt:    t.CreatedAt,
	}
}

type RespList struct {
	Response
	Teams []*RTeam `json:"teams"`
}

func BindList(db *gorm.DB, teams []cpt.Team) []*RTeam {
	var out []*RTeam
	for _, team := range teams {
		out = append(out, BindTeam(db, &team))
	}
	return out
}

type RespCreate struct {
	Response
	TeamID uint `json:"team_id,omitempty"`
}

type RespShow struct {
	Response
	Team *RTeam `json:"team"`
}

type RTeamUser struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	EmailHash string `json:"email_hash"`
	Nickname  string `json:"nickname"`
}

func BindUser(u *cpt.User) *RTeamUser {
	return &RTeamUser{
		ID:        u.ID,
		Username:  u.Username,
		EmailHash: random.MD5(u.Email),
		Nickname:  u.Profile.Nickname,
	}
}

type RespUserList struct {
	Response
	Members []*RTeamUser `json:"members"`
}

func BindUserList(users []cpt.User) []*RTeamUser {
	var out []*RTeamUser
	for _, user := range users {
		out = append(out, BindUser(&user))
	}
	return out
}

type RespUserAddDeleteDetail struct {
	Response
	UserID uint `json:"user_id"`
}

type RespUserAdd struct {
	Response
	Details []*RespUserAddDeleteDetail `json:"details"`
}

type RespUserDelete struct {
	Response
	Details []*RespUserAddDeleteDetail `json:"details"`
}

type RStatisticContest struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Points uint   `json:"points"`
	Rank   uint   `json:"rank"`
}

type RespStatisticShow struct {
	Response
	Contests []*RStatisticContest `json:"contests"`
}

type RInvitation struct {
	ToUser     uint      `json:"to_user"`
	FromUser   uint      `json:"from_user"`
	InviteTime time.Time `json:"invite_time"`
	Status     uint      `json:"status"`
}

func BindInvitation(inv *cpt.TeamInvitation) *RInvitation {
	return &RInvitation{
		ToUser:     inv.ToUser,
		FromUser:   inv.FromUser,
		InviteTime: inv.CreatedAt,
		Status:     inv.Status,
	}
}

type RespInvitationList struct {
	Response
	Invitations []*RInvitation `json:"invitations"`
}

func BindInvitationList(invitations []cpt.TeamInvitation) []*RInvitation {
	var out []*RInvitation
	for _, invitation := range invitations {
		out = append(out, BindInvitation(&invitation))
	}
	return out
}

type RespInvitationToken struct {
	Response
	Token string `json:"token"`
}

type RApplication struct {
	UserID    uint      `json:"user_id"`
	ApplyTime time.Time `json:"apply_time"`
	Status    uint      `json:"status"`
}

func BindApplication(app *cpt.TeamApplication) *RApplication {
	return &RApplication{
		UserID:    app.UserID,
		ApplyTime: app.CreatedAt,
		Status:    app.Status,
	}
}

type RespApplicationList struct {
	Response
	Applications []*RApplication `json:"applications"`
}

func BindApplicationList(applications []cpt.TeamApplication) []*RApplication {
	var out []*RApplication
	for _, application := range applications {
		out = append(out, BindApplication(&application))
	}
	return out
}
