package components

import (
	"sort"
	"time"

	"github.com/FlagField/FlagField-Server/internal/pkg/constants"

	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	ContestAccessRegister = iota
	ContestAccessPrivate
	TableContest             = "contest"
	TableContestNotification = "contest_notification"
	TableRelContestAdmin     = "rel_contest_admin"
	TableContestRank         = "contest_rank"
)

var (
	DefaultContest = Contest{
		Name:        "~FlagField__defaultContest",
		Description: "This is the default contest for Practice. Please do not delete this.",
		StartTime:   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		EndTime:     constants.MaxTime,
		Access:      0,
		IsHidden:    true,
		CreatorID:   0,
	}
)

type Contest struct {
	gorm.Model
	Name          string `gorm:"size:50;not null"`
	Description   string `gorm:"type:longtext"`
	StartTime     time.Time
	EndTime       time.Time
	Access        uint
	IsHidden      bool
	CreatorID     uint
	Problems      []*Problem
	Admins        []*User         `gorm:"many2many:rel_contest_admin"`
	TeamSnapshots []*TeamSnapshot `gorm:"foreignkey:ContestID"`
	Resources     []*Resource
	Notifications []*ContestNotification
}

func (c *Contest) Update(db *gorm.DB) {
	if v := db.Save(c); v.Error != nil {
		panic(v.Error)
	}
}

func (c *Contest) Delete(db *gorm.DB) {
	if db.NewRecord(c) {
		return
	}
	if v := db.Unscoped().Delete(c); v.Error != nil {
		panic(v.Error)
	}
}

// Get association
func (c *Contest) GetNotificationByOrder(db *gorm.DB, order int) (*ContestNotification, error) {
	if order <= 0 {
		return nil, errors.ErrOutOfRange
	}
	var notification ContestNotification
	if v := db.Table(TableContestNotification).Where("contest_id = ?", c.ID).Offset(order - 1).Limit(1).First(&notification); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrOutOfRange
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &notification, nil
}

func (c *Contest) GetNotifications(db *gorm.DB, from time.Time, offset uint) []ContestNotification {
	var notifications []ContestNotification
	if v := db.Table(TableContestNotification).Where("contest_id = ? and created_at > ?", c.ID, from).Order("created_at desc").Offset(offset).Find(&notifications); v.Error != nil {
		panic(v.Error)
	}
	return notifications
}

func (c *Contest) GetProblems(db *gorm.DB) []Problem {
	var problems []Problem
	if v := db.Model(c).Association("Problems").Find(&problems); v.Error != nil {
		panic(v.Error)
	}
	return problems
}

func (c *Contest) GetTeams(db *gorm.DB) []Team {
	var teams []Team
	if v := db.Table(TableTeam).Where("id in (?)",
		db.Table(TableTeamSnapshot).Select("team_id").Where("contest_id = ?", c.ID).SubQuery(),
	).Find(&teams); v.Error != nil {
		panic(v.Error)
	}
	return teams
}

func (c *Contest) GetTeamSnapshots(db *gorm.DB) []TeamSnapshot {
	var snapshots []TeamSnapshot
	if v := db.Model(c).Association("TeamSnapshots").Find(&snapshots); v.Error != nil {
		panic(v.Error)
	}
	return snapshots
}

func (c *Contest) GetUserTeam(db *gorm.DB, user *User) (*Team, error) {
	var team Team
	if v := db.Where("id = ?",
		db.Table(TableTeamSnapshot).Select("team_id").Where("contest_id = ? and id in (?)", c.ID,
			db.Table(TableRelTeamSnapshotMember).Select("team_snapshot_id").Where("user_id = ?", user.ID).SubQuery(),
		).SubQuery(),
	).First(&team); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &team, nil
}

func (c *Contest) GetAdmins(db *gorm.DB) []User {
	var admins []User
	if v := db.Table(TableUser).Where("id in (?)",
		db.Table(TableRelContestAdmin).Select("user_id").Where("contest_id = ?", c.ID).SubQuery(),
	).Preload("Profile").Find(&admins); v.Error != nil {
		panic(v.Error)
	}
	return admins
}

// Add association
func (c *Contest) AddProblem(db *gorm.DB, problem *Problem) {
	if v := db.Model(c).Association("Problems").Append(problem); v.Error != nil {
		panic(v.Error)
	}
}

func (c *Contest) AddAdmin(db *gorm.DB, user *User) {
	if v := db.Model(c).Association("Admins").Append(user); v.Error != nil {
		panic(v.Error)
	}
}

func (c *Contest) AddTeam(db *gorm.DB, team *Team, usersID []uint) error {
	snapshot, err := team.NewSnapshot(db, c, usersID)
	if err != nil {
		return err
	}
	if v := db.Model(c).Association("TeamSnapshots").Append(snapshot); v.Error != nil {
		panic(v.Error)
	}
	return nil
}

func (c *Contest) AddResource(db *gorm.DB, res *Resource) {
	if v := db.Model(c).Association("Resources").Append(res); v.Error != nil {
		panic(v.Error)
	}
}

func (c *Contest) AddNotification(db *gorm.DB, content string) int {
	if v := db.Model(c).Association("Notifications").Append(&ContestNotification{
		Content: content,
	}); v.Error != nil {
		panic(v.Error)
	}
	return db.Model(c).Association("Notifications").Count()
}

// TeamDelete association
func (c *Contest) DeleteTeamByID(db *gorm.DB, teamID uint) error {
	var snapshot TeamSnapshot
	if v := db.Table(TableTeamSnapshot).Where("contest_id = ? and team_id = ?", c.ID, teamID).First(&snapshot); v.Error == gorm.ErrRecordNotFound {
		return errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	snapshot.Delete(db)
	return nil
}

func (c *Contest) DeleteAdmin(db *gorm.DB, user *User) {
	if v := db.Model(c).Association("Admins").Delete(user); v.Error == gorm.ErrRecordNotFound {
		return
	} else if v.Error != nil {
		panic(v.Error)
	}
}

// Permission check
func (c *Contest) IsAdmin(db *gorm.DB, user *User) bool {
	var count uint
	if v := db.Table(TableRelContestAdmin).Where("contest_id = ? and user_id = ?", c.ID, user.ID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func (c *Contest) IsPlayer(db *gorm.DB, user *User) bool {
	var count uint
	if v := db.Table(TableRelTeamSnapshotMember).Where("user_id = ? and team_snapshot_id in (?)", user.ID,
		db.Table(TableTeamSnapshot).Select("id").Where("contest_id = ?", c.ID).SubQuery(),
	).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func (c *Contest) HasPlayerByUserID(db *gorm.DB, usersID []uint) bool {
	var count uint
	if v := db.Table(TableRelTeamSnapshotMember).Where("user_id in (?) and team_snapshot_id in (?)", usersID,
		db.Table(TableTeamSnapshot).Select("id").Where("contest_id = ?", c.ID).SubQuery(),
	).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func (c *Contest) HasTeamByTeamID(db *gorm.DB, teamID uint) bool {
	var count uint
	if v := db.Table(TableTeamSnapshot).Where("team_id = ? and contest_id = ?", teamID, c.ID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

// Data proceed
func (c *Contest) GenerateRank(db *gorm.DB) {
	teams := c.GetTeams(db)
	ranks := make([]ContestRank, len(teams))
	for i, team := range teams {
		ranks[i] = ContestRank{
			ContestID: c.ID,
			TeamID:    team.ID,
			Points:    team.CalcContestPoints(db, c),
		}
	}
	sort.Slice(ranks, func(i, j int) bool {
		return ranks[i].Points > ranks[j].Points
	})
	for i := 0; i < len(teams); i++ {
		if i == 0 || ranks[i].Points != ranks[i-1].Points {
			ranks[i].Rank = uint(i + 1)
		} else {
			ranks[i].Rank = ranks[i-1].Rank
		}
		if v := db.Save(&ranks[i]); v.Error != nil {
			panic(v.Error)
		}
	}
}

type ContestNotification struct {
	gorm.Model
	ContestID uint
	Content   string
}

func (n *ContestNotification) Delete(db *gorm.DB) {
	if db.NewRecord(n) {
		return
	}
	if v := db.Unscoped().Delete(n); v.Error != nil {
		panic(v.Error)
	}
}

type ContestRank struct {
	gorm.Model
	ContestID uint
	TeamID    uint
	Points    uint
	Rank      uint
}

func NewContest(db *gorm.DB, name string, description string, startTime time.Time, endTime time.Time, access uint, creatorID uint) *Contest {
	contest := &Contest{
		Name:        name,
		Description: description,
		StartTime:   startTime,
		EndTime:     endTime,
		Access:      access,
		IsHidden:    true,
		CreatorID:   creatorID,
	}
	if v := db.Create(contest); v.Error != nil {
		panic(v.Error)
	}
	return contest
}

func GetContests(db *gorm.DB) []Contest {
	var contests []Contest
	if v := db.Find(&contests); v.Error != nil {
		panic(v.Error)
	}
	return contests
}

func GetContestsFilterByStatus(db *gorm.DB, status uint) []Contest {
	var contests []Contest
	if status == 0 {
		if v := db.Where("start_time > ?", time.Now()).Find(&contests); v.Error != nil {
			panic(v.Error)
		}
	} else if status == 1 {
		if v := db.Where("start_time <= ? and ? <= end_time", time.Now(), time.Now()).Find(&contests); v.Error != nil {
			panic(v.Error)
		}
	} else {
		if v := db.Where("end_time < ?", time.Now()).Find(&contests); v.Error != nil {
			panic(v.Error)
		}
	}
	return contests
}

func GetContestByID(db *gorm.DB, id uint) (*Contest, error) {
	var contest Contest
	if v := db.First(&contest, id); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &contest, nil
}

func GetDefaultContest(db *gorm.DB) *Contest {
	if v := db.Where("id = ?", DefaultContest.ID).First(&DefaultContest); v.Error != nil {
		panic(v.Error)
	}
	return &DefaultContest
}

func GetContestCount(db *gorm.DB) uint {
	var count uint
	if v := db.Table(TableContest).Where("id <> ?", DefaultContest.ID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count
}

func IsContestAdmin(db *gorm.DB, contestID, userID uint) bool {
	var count uint
	if v := db.Table(TableRelContestAdmin).Where("contest_id = ? and user_id = ?", contestID, userID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}
