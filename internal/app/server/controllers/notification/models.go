package notification

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"time"
)

type ReqNew struct {
	UserID  uint   `json:"user_id" validate:"required" binding:"required"`
	Content string `json:"content" validate:"required" binding:"required"`
}

type RNotification struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreateTime time.Time `json:"create_time"`
}

func BindNotification(n *cpt.Notification) *RNotification {
	return &RNotification{
		ID:         n.ID,
		UserID:     n.UserID,
		Content:    n.Content,
		IsRead:     n.IsRead,
		CreateTime: n.CreatedAt,
	}
}

type RespList struct {
	Response
	Notifications []*RNotification `json:"notifications"`
}

func BindList(notifications []cpt.Notification) []*RNotification {
	var out []*RNotification
	for _, notification := range notifications {
		out = append(out, BindNotification(&notification))
	}
	return out
}

type RespNew struct {
	Response
	NotificationID uint `json:"notification_id,omitempty"`
}
