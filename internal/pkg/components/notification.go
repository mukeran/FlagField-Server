package components

import (
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/jinzhu/gorm"
)

const (
	TableNotification = "notification"
)

type Notification struct {
	gorm.Model
	UserID    uint
	Content   string `gorm:"type:longtext"`
	IsRead    bool
	CreatorID uint
}

func (n *Notification) Update(db *gorm.DB) {
	if v := db.Save(n); v.Error != nil {
		panic(v.Error)
	}
}

func (n *Notification) Delete(db *gorm.DB) {
	if v := db.Unscoped().Delete(n); v.Error != nil {
		panic(v.Error)
	}
}

func NewNotification(db *gorm.DB, userID uint, content string, creatorID uint) *Notification {
	notification := Notification{
		UserID:    userID,
		Content:   content,
		IsRead:    false,
		CreatorID: creatorID,
	}
	if v := db.Save(&notification); v.Error != nil {
		panic(v.Error)
	}
	return &notification
}

func GetNotifications(db *gorm.DB, offset uint, limit uint) []Notification {
	var out []Notification
	if v := db.Offset(offset).Limit(limit).Find(&out); v.Error != nil {
		panic(v.Error)
	}
	return out
}

func GetUserNotifications(db *gorm.DB, user *User, offset uint) []Notification {
	var out []Notification
	if v := db.Where("user_id = ? and created_at >= ?", user.ID, user.CreatedAt).Offset(offset).Find(&out); v.Error != nil {
		panic(v.Error)
	}
	return out
}

func GetNotificationByID(db *gorm.DB, id uint) (*Notification, error) {
	var out Notification
	if v := db.First(&out, id); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &out, nil
}
