package components

import (
	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	TableHint = "hint"
)

type Hint struct {
	gorm.Model
	ProblemID     uint
	Cost          uint
	Content       string  `gorm:"type:longtext"`
	UnlockedUsers []*User `gorm:"many2many:rel_user_unlocked_hint"`
}

func (h *Hint) Unlock(db *gorm.DB, contest *Contest, problem *Problem, user *User) error {
	if user.IsAdmin {
		return nil
	}
	team, _ := contest.GetUserTeam(db, user)
	points := team.CalcContestPoints(db, contest)
	if points < h.Cost {
		return errors.ErrNotEnoughPoints
	}
	if v := db.Model(h).Association("UnlockedUsers").Append(user); v.Error != nil {
		panic(v.Error)
	}
	return nil
}

func (h *Hint) IsUnlocked(db *gorm.DB, contest *Contest, problem *Problem, user *User) bool {
	if user.IsAdmin {
		return true
	}
	team, _ := contest.GetUserTeam(db, user)
	var count uint
	if v := db.Table(TableRelUserUnlockedHint).Where("hint_id = ? and user_id in (?)", h.ID,
		db.Table(TableRelTeamSnapshotMember).Select("user_id").Where("team_id = ? and contest_id = ?", team.ID, contest.ID).SubQuery(),
	).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func (h *Hint) Update(db *gorm.DB) {
	if v := db.Save(h); v.Error != nil {
		panic(v.Error)
	}
}

func (h *Hint) Delete(db *gorm.DB) {
	if db.NewRecord(h) {
		return
	}
	if v := db.Unscoped().Delete(h); v.Error != nil {
		panic(v.Error)
	}
}
