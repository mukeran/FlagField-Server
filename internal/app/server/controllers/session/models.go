package session

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"time"
)

type ReqCreate struct {
	Username string `json:"username" validate:"max=20,min=1,keyword" binding:"max=20,min=1,keyword"`
	Password string `json:"password" validate:"max=50,min=6" binding:"max=50,min=6"`
}

type RSession struct {
	SessionID string    `json:"session_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	IsAdmin   bool      `json:"is_admin,omitempty"`
	ExpireAt  time.Time `json:"expire_at"`
}

func BindSession(session *cpt.Session) *RSession {
	return &RSession{
		SessionID: session.SessionID,
		UserID:    session.UserID,
		Username:  session.User.Username,
		IsAdmin:   session.User.IsAdmin,
		ExpireAt:  session.ExpireAt,
	}
}

type RespList struct {
	Response
	Sessions []*RSession `json:"sessions"`
}

func BindList(sessions []cpt.Session) []*RSession {
	var out []*RSession
	for _, session := range sessions {
		out = append(out, BindSession(&session))
	}
	return out
}

type RespCreate struct {
	Response
	Session *RSession `json:"session,omitempty"`
}

type RespViewCurrent struct {
	Response
	Session *RSession `json:"session,omitempty"`
}
