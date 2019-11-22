package components

import (
	"encoding/json"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/constants"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
)

const (
	TableSession = "session"
)

type Session struct {
	gorm.Model
	SessionID string `gorm:"size:32;unique;not null"`
	User      *User
	UserID    uint
	ExpireAt  time.Time
	// Logs      []SessionLog `gorm:"foreignkey:SessionNumID"`
}

func (s *Session) IsExpired() bool {
	return s.ExpireAt.Before(time.Now())
}

func (s *Session) Renew(db *gorm.DB, d time.Duration) {
	s.ExpireAt = s.ExpireAt.Add(d)
	s.Update(db)
}

func (s *Session) Update(db *gorm.DB) {
	if v := db.Save(s); v.Error != nil {
		panic(v.Error)
	}
}

func (s *Session) IsLoggedIn() bool {
	return s.UserID != DefaultUser.ID
}

func (s *Session) Destroy(db *gorm.DB, rp *redis.Pool) {
	if db.NewRecord(s) {
		return
	}
	conn := rp.Get()
	if conn.Err() != nil {
		panic(conn.Err())
	}
	_, err := conn.Do("HDEL", s.SessionID, "*")
	if err != nil {
		panic(err)
	}
	if v := db.Unscoped().Delete(s); v.Error != nil {
		panic(v.Error)
	}
}

func (s *Session) Get(rp *redis.Pool, key string, value interface{}) {
	conn := rp.Get()
	if conn.Err() != nil {
		panic(conn.Err())
	}
	res, err := redis.Bytes(conn.Do("HGET", s.SessionID, key))
	if err == redis.ErrNil {
		return
	}
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(res, &value)
	if err != nil {
		panic(err)
	}
}

// value must have JSON tag in order to marshal into a JSON string
func (s *Session) Set(rp *redis.Pool, key string, value interface{}) {
	conn := rp.Get()
	if conn.Err() != nil {
		panic(conn.Err())
	}
	data, err := json.Marshal(value)
	if err != nil {
		panic(err)
	}
	_, err = conn.Do("HSET", s.SessionID, key, data)
	if err != nil {
		panic(err)
	}
}

func (s *Session) Del(rp *redis.Pool, key string) {
	conn := rp.Get()
	if conn.Err() != nil {
		panic(conn.Err())
	}
	_, err := conn.Do("HDEL", s.SessionID, key)
	if err != nil {
		panic(err)
	}
}

func generateSession() string {
	return random.RandomString(32, constants.DicLetterNumeric)
}

func NewSession(db *gorm.DB, u *User, maxAge int) *Session {
	s := &Session{
		SessionID: generateSession(),
		User:      u,
		UserID:    u.ID,
		ExpireAt:  time.Now().Add(time.Duration(maxAge) * time.Second),
	}
	if v := db.Create(s); v.Error != nil {
		panic(v.Error)
	}
	return s
}

func GetSessionBySessionID(db *gorm.DB, sessionID string) (*Session, error) {
	var s Session
	if v := db.Where("session_id = ?", sessionID).Preload("User").First(&s); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &s, nil
}

func GetSessions(db *gorm.DB) []Session {
	var sessions []Session
	if v := db.Preload("User").Find(&sessions); v.Error != nil {
		panic(v.Error)
	}
	return sessions
}

// type SessionLog struct {
// 	gorm.Model
// 	IP           string `gorm:"size:128"`
// 	UserAgent    string
// 	RequestPath  string
// 	RequestBody  string
// 	SessionNumID uint
// }
