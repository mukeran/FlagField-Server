package components

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/constants"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
)

const (
	DefaultUserTeamPrefix      = "~FlagField__defaultTeam"
	DefaultUserTeamDescription = "This is the default team for a single user."
	TableUser                  = "user"
	TableUserProfile           = "user_profile"
	TableRelUserTeam           = "rel_user_team"
	TableRelUserUnlockedHint   = "rel_user_unlocked_hint"
)

var (
	DefaultUser = User{
		Username:     "~FlagField__guest",
		PasswordHash: "",
	}
)

type User struct {
	gorm.Model
	Username         string `gorm:"size:20;unique;not null"`
	PasswordHash     string `gorm:"not null"`
	Email            string `gorm:"size:100"`
	Profile          *UserProfile
	IsAdmin          bool
	IsHidden         bool
	IsBanned         bool
	DefaultTeam      uint
	AsContestAdmin   []*Contest      `gorm:"many2many:rel_contest_admin"`
	AsTeamAdmin      []*Team         `gorm:"many2many:rel_team_admin"`
	SolvedProblems   []*Problem      `gorm:"many2many:rel_problem_solver"`
	UnlockedHints    []*Hint         `gorm:"many2many:rel_user_unlocked_hint"`
	Submissions      []*Submission   `gorm:"foreignkey:CreatorID"`
	CreatedContests  []*Contest      `gorm:"foreignkey:CreatorID"`
	CreatedProblems  []*Problem      `gorm:"foreignkey:CreatorID"`
	CreatedTeams     []*Team         `gorm:"foreignkey:CreatorID"`
	CreatedResources []*Resource     `gorm:"foreignkey:CreatorID"`
	Teams            []*Team         `gorm:"many2many:rel_user_team"`
	TeamSnapshots    []*TeamSnapshot `gorm:"many2many:rel_team_snapshot_member"`
}

func (u *User) Update(db *gorm.DB) {
	if v := db.Save(&u.Profile); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Save(u); v.Error != nil {
		panic(v.Error)
	}
}

func (u *User) Delete(db *gorm.DB) {
	if db.NewRecord(u) {
		return
	}
	if v := db.Unscoped().Delete(Session{}, "user_id = ?", u.ID); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Unscoped().Delete(Team{}, "id = ?", u.DefaultTeam); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Unscoped().Delete(&u.Profile); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Unscoped().Delete(u); v.Error != nil {
		panic(v.Error)
	}
}

// Information control
func (u *User) MatchPassword(password string) bool {
	v := strings.Split(u.PasswordHash, "|")
	if len(v) != 2 {
		return false
	}
	salt, passwordHash := v[0], v[1]
	if passwordHash != GetPasswordHash(password, salt) {
		return false
	}
	return true
}

func (u *User) ChangePassword(db *gorm.DB, password string) {
	salt := GenerateSalt()
	u.PasswordHash = salt + "|" + GetPasswordHash(password, salt)
	u.Update(db)
}

// Get associations
func (u *User) GetSolvedProblems(db *gorm.DB, contest *Contest) []Problem {
	var problems []Problem
	if v := db.Table(TableProblem).Where("contest_id = ? and id in (?)", contest.ID,
		db.Table(TableRelProblemSolver).Select("problem_id").Where("user_id = ?", u.ID).SubQuery(),
	).Find(&problems); v.Error != nil {
		panic(v.Error)
	}
	return problems
}

func (u *User) CalcTeamContestPoints(db *gorm.DB, contest *Contest) (uint, error) {
	team, err := contest.GetUserTeam(db, u)
	if err != nil {
		return 0, err
	}
	return team.CalcContestPoints(db, contest), nil
}

func (u *User) CalcPersonalContestPoints(db *gorm.DB, contest *Contest) (points uint) {
	problems := u.GetSolvedProblems(db, contest)
	for _, problem := range problems {
		points += problem.Points
	}
	return
}

func (u *User) GetContestTeam(db *gorm.DB, contest *Contest) (*Team, error) {
	return contest.GetUserTeam(db, u)
}

// Permission check
func (u *User) IsContestAdmin(db *gorm.DB, contest *Contest) bool {
	return contest.IsAdmin(db, u)
}

func (u *User) IsResourceAdmin(res *Resource) bool {
	return res.CreatorID == u.ID
}

func (u *User) IsTeamAdmin(db *gorm.DB, team *Team) bool {
	return team.IsAdmin(db, u)
}

func (u *User) HasContestAccess(db *gorm.DB, contest *Contest) bool {
	return contest.IsPlayer(db, u) || contest.IsAdmin(db, u)
}

func (u *User) HasResourceAccess(db *gorm.DB, res *Resource) bool {
	if u.IsResourceAdmin(res) {
		return true
	}
	contest, err := GetContestByID(db, res.ContestID)
	if err != nil {
		return false
	}
	if u.HasContestAccess(db, contest) {
		return true
	}
	return false
}

func (u *User) HasTeamAccess(db *gorm.DB, team *Team) bool {
	return team.IsMember(db, u) || team.IsAdmin(db, u)
}

type UserProfile struct {
	gorm.Model
	Nickname    string `gorm:"size:20"`
	Page        string `gorm:"size:100"`
	Description string `gorm:"size:1000"`
	UserID      uint
}

func HasUsername(db *gorm.DB, username string) bool {
	var count uint
	if v := db.Table(TableUser).Where("username = ?", username).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func HasEmail(db *gorm.DB, email string) bool {
	var count uint
	if v := db.Table(TableUser).Where("email = ?", email).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func GenerateSalt() string {
	return random.RandomString(16, constants.DicLetterNumeric)
}

func GetPasswordHash(password string, salt string) string {
	b := sha256.Sum256([]byte(salt + password))
	return base64.StdEncoding.EncodeToString(b[:])
}

func NewUser(db *gorm.DB, username string, password string) (*User, error) {
	if HasUsername(db, username) {
		return nil, errors.ErrDuplicatedUsername
	}
	salt := GenerateSalt()
	u := &User{
		Username:     username,
		PasswordHash: salt + "|" + GetPasswordHash(password, salt),
	}
	if v := db.Create(&u); v.Error != nil {
		panic(v.Error)
	}
	u.Profile = &UserProfile{
		UserID: u.ID,
	}
	if v := db.Create(&u.Profile); v.Error != nil {
		panic(v.Error)
	}
	team := NewTeam(db, fmt.Sprintf("%s_%s", DefaultUserTeamPrefix, username), DefaultUserTeamDescription, 0)
	team.IsDefault = true
	team.Update(db)
	team.AddMember(db, u)
	_ = DefaultContest.AddTeam(db, team, []uint{u.ID})
	u.DefaultTeam = team.ID
	if v := db.Model(&u).Updates(User{DefaultTeam: team.ID}); v.Error != nil {
		panic(v.Error)
	}
	return u, nil
}

func GetUsers(db *gorm.DB) []User {
	var users []User
	if v := db.Not("id", DefaultUser.ID).Preload("Profile").Find(&users); v.Error != nil {
		panic(v.Error)
	}
	return users
}

func GetUsersWithQuery(db *gorm.DB, offset, limit uint, query string) []User {
	var users []User
	if v := db.Not("id", DefaultUser.ID).Where("id like ? or username like ? or email like ?",
		"%"+query+"%", "%"+query+"%", "%"+query+"%").Preload("Profile").Offset(offset).Limit(limit).Find(&users); v.Error != nil {
		panic(v.Error)
	}
	return users
}

func GetUserByUsername(db *gorm.DB, username string) (*User, error) {
	var user User
	if v := db.Where("username = ?", username).Preload("Profile").First(&user); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &user, nil
}

func GetUserByID(db *gorm.DB, id uint) (*User, error) {
	var user User
	if v := db.Preload("Profile").First(&user, id); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &user, nil
}

func GetUsersByID(db *gorm.DB, id []uint) []User {
	var users []User
	if v := db.Preload("Profile").Where("id in (?)", id).Find(&users); v.Error != nil {
		panic(v.Error)
	}
	return users
}

func GetAdmins(db *gorm.DB) []User {
	var admins []User
	if v := db.Not("id", DefaultUser.ID).Preload("Profile").Where("is_admin = 1").Find(&admins); v.Error != nil {
		panic(v.Error)
	}
	return admins
}

func GetUserCount(db *gorm.DB) uint {
	var count uint
	if v := db.Table(TableUser).Where("id <> ?", DefaultUser.ID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count
}
