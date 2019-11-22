package user

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
	"github.com/jinzhu/gorm"
	"time"
)

type ReqRegister struct {
	Username     string `json:"username" validate:"max=20,min=1,keyword" binding:"max=20,min=1,keyword"`
	Password     string `json:"password" validate:"max=50,min=6" binding:"max=50,min=6"`
	Email        string `json:"email" validate:"email" binding:"email"`
	EmailCaptcha string `json:"email_captcha,omitempty"`
}

type ReqPasswordModify struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password" validate:"max=50,min=6" binding:"max=50,min=6"`
}

type ReqEmailModify struct {
	Email string `json:"email" validate:"email" binding:"email"`
}

type ReqProfileModify struct {
	Nickname    *string `json:"nickname" validate:"omitempty,max=50,min=1" binding:"omitempty,max=50,min=1"`
	Page        *string `json:"page" validate:"omitempty,url" validate:"omitempty,url"`
	Description *string `json:"description"`
}

func (r *ReqProfileModify) Bind(p *cpt.UserProfile) {
	if r.Nickname != nil {
		p.Nickname = *r.Nickname
	}
	if r.Page != nil {
		p.Page = *r.Page
	}
	if r.Description != nil {
		p.Description = *r.Description
	}
}

type ReqPermissionAdd struct {
	Permissions []string `json:"permissions"`
}

type RUser struct {
	ID        uint          `json:"id"`
	Username  string        `json:"username"`
	Email     string        `json:"email,omitempty"`
	EmailHash string        `json:"email_hash"`
	Profile   *RUserProfile `json:"profile"`
	IsAdmin   bool          `json:"is_admin,omitempty"`
}

type RUserProfile struct {
	Nickname    string `json:"nickname"`
	Page        string `json:"page"`
	Description string `json:"description"`
}

func BindUser(session *cpt.Session, user *cpt.User) *RUser {
	u := &RUser{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		EmailHash: random.MD5(user.Email),
		IsAdmin:   user.IsAdmin,
	}
	if user.Profile != nil {
		u.Profile = &RUserProfile{
			Nickname:    user.Profile.Nickname,
			Page:        user.Profile.Page,
			Description: user.Profile.Description,
		}
	}
	if !session.User.IsAdmin {
		u.Email = ""
	}
	return u
}

type RespRegister struct {
	Response
	UserID uint `json:"user_id,omitempty"`
}

type RespList struct {
	Response
	Users []*RUser `json:"users"`
}

func BindList(session *cpt.Session, users []cpt.User) []*RUser {
	var out []*RUser
	for _, user := range users {
		out = append(out, BindUser(session, &user))
	}
	return out
}

type RespShow struct {
	Response
	User *RUser `json:"user"`
}

type RespEmailHashShow struct {
	Response
	EmailHash string `json:"email_hash"`
}

type RTeam struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	IsAdmin bool   `json:"is_admin,omitempty"`
}

func BindTeam(db *gorm.DB, team *cpt.Team, user *cpt.User) *RTeam {
	return &RTeam{
		ID:      team.ID,
		Name:    team.Name,
		IsAdmin: team.IsAdmin(db, user),
	}
}

type RespTeamList struct {
	Response
	Teams []*RTeam `json:"teams"`
}

func BindTeamList(db *gorm.DB, teams []cpt.Team, user *cpt.User) []*RTeam {
	var out []*RTeam
	for _, team := range teams {
		out = append(out, BindTeam(db, &team, user))
	}
	return out
}

type RTeamInvitationAndApplication struct {
	TeamID uint      `json:"team_id"`
	Time   time.Time `json:"time"`
	Status uint      `json:"status"`
}

func BindInvitation(inv *cpt.TeamInvitation) *RTeamInvitationAndApplication {
	return &RTeamInvitationAndApplication{
		TeamID: inv.TeamID,
		Time:   inv.CreatedAt,
		Status: inv.Status,
	}
}

func BindApplication(app *cpt.TeamApplication) *RTeamInvitationAndApplication {
	return &RTeamInvitationAndApplication{
		TeamID: app.TeamID,
		Time:   app.CreatedAt,
		Status: app.Status,
	}
}

type RespTeamInvitationList struct {
	Response
	Invitations []*RTeamInvitationAndApplication `json:"invitations"`
}

type RespTeamApplicationList struct {
	Response
	Applications []*RTeamInvitationAndApplication `json:"applications"`
}

func BindTeamInvitationList(invitations []cpt.TeamInvitation) []*RTeamInvitationAndApplication {
	var out []*RTeamInvitationAndApplication
	for _, invitation := range invitations {
		out = append(out, BindInvitation(&invitation))
	}
	return out
}

func BindTeamApplicationList(applications []cpt.TeamApplication) []*RTeamInvitationAndApplication {
	var out []*RTeamInvitationAndApplication
	for _, application := range applications {
		out = append(out, BindApplication(&application))
	}
	return out
}
