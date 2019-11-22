package components

import (
	"github.com/FlagField/FlagField-Server/internal/pkg/constants"
	"github.com/FlagField/FlagField-Server/internal/pkg/random"
	"github.com/jinzhu/gorm"
	"time"

	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	TableTeam                  = "team"
	TableRelTeamAdmin          = "rel_team_admin"
	TableTeamSnapshot          = "team_snapshot"
	TableRelTeamSnapshotMember = "rel_team_snapshot_member"
	TableTeamInvitation        = "team_invitation"
	TableTeamApplication       = "team_application"
)

type Team struct {
	gorm.Model
	Name            string          `gorm:"size:50;unique;not null"`
	Description     string          `gorm:"type:longtext"`
	Members         []*User         `gorm:"many2many:rel_user_team"`
	Admins          []*User         `gorm:"many2many:rel_team_admin"`
	Snapshots       []*TeamSnapshot `gorm:"foreignkey:TeamID"`
	CreatorID       uint
	IsDefault       bool
	InvitationToken string `gorm:"size:32"`
}

func (t *Team) Update(db *gorm.DB) {
	if v := db.Save(t); v.Error != nil {
		panic(v.Error)
	}
}

func (t *Team) Delete(db *gorm.DB) {
	if db.NewRecord(t) {
		return
	}
	if v := db.Unscoped().Model(t).Association("Members").Clear(); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Unscoped().Model(t).Association("Admins").Clear(); v.Error != nil {
		panic(v.Error)
	}
	snapshot := t.GetSnapshots(db)
	for _, s := range snapshot {
		s.Delete(db)
	}
	if v := db.Unscoped().Delete(t); v.Error != nil {
		panic(v.Error)
	}
}

// Get associations
func (t *Team) GetMembers(db *gorm.DB, query string) []User {
	var users []User
	if v := db.Table(TableUser).Where("id in (?)",
		db.Table(TableRelUserTeam).Select("user_id").Where("team_id = ?", t.ID).SubQuery(),
	).Where("id like ? or username like ?", "%"+query+"%", "%"+query+"%").Preload("Profile").Find(&users); v.Error != nil {
		panic(v.Error)
	}
	return users
}

func (t *Team) GetMembersCount(db *gorm.DB) uint {
	var count uint
	if v := db.Table(TableRelUserTeam).Where("team_id = ?", t.ID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count
}

func (t *Team) GetAdmins(db *gorm.DB) []User {
	var admins []User
	if v := db.Table(TableUser).Where("id in (?)",
		db.Table(TableRelTeamAdmin).Select("user_id").Where("team_id = ?", t.ID).SubQuery(),
	).Preload("Profile").Find(&admins); v.Error != nil {
		panic(v.Error)
	}
	return admins
}

func (t *Team) GetContests(db *gorm.DB) []Contest {
	var contests []Contest
	if v := db.Table(TableContest).Where("id in (?)",
		db.Table(TableTeamSnapshot).Select("contest_id").Where("team_id = ?", t.ID).SubQuery(),
	).Find(&contests); v.Error != nil {
		panic(v.Error)
	}
	return contests
}

func (t *Team) GetSnapshots(db *gorm.DB) []TeamSnapshot {
	var snapshots []TeamSnapshot
	if v := db.Model(t).Association("Snapshots").Find(&snapshots); v.Error != nil {
		panic(v.Error)
	}
	return snapshots
}

func (t *Team) GetSolvedProblems(db *gorm.DB, contest *Contest) []Problem {
	var problems []Problem
	if v := db.Table(TableProblem).Where("contest_id = ? and id in (?)", contest.ID,
		db.Table(TableRelProblemSolver).Select("problem_id").Where("user_id in (?)",
			db.Table(TableRelTeamSnapshotMember).Select("user_id").Where("team_snapshot_id = ?",
				db.Table(TableTeamSnapshot).Select("id").Where("team_id = ? and contest_id = ?", t.ID, contest.ID).SubQuery(),
			).SubQuery(),
		).SubQuery(),
	).Find(&problems); v.Error != nil {
		panic(v.Error)
	}
	return problems
}

func (t *Team) GetCorrectSubmissions(db *gorm.DB, contest *Contest) []Submission {
	var submissions []Submission
	problems := t.GetSolvedProblems(db, contest)
	for _, problem := range problems {
		var submission Submission
		if v := db.Table(TableSubmission).Where("problem_id = ? and creator_id in (?)", problem.ID,
			db.Table(TableRelTeamSnapshotMember).Select("user_id").Where("team_id = ? and contest_id = ?", t.ID, contest.ID).SubQuery(),
		).Order("created_at").First(&submission); v.Error != nil {
			panic(v.Error)
		}
		submissions = append(submissions, submission)
	}
	return submissions
}

func (t *Team) GetUnlockedHints(db *gorm.DB, contest *Contest) []Hint {
	var hints []Hint
	if v := db.Table(TableHint).Where("id in (?) and problem_id in (?)",
		db.Table(TableRelUserUnlockedHint).Select("hint_id").Where("user_id in (?)",
			db.Table(TableRelTeamSnapshotMember).Select("user_id").Where("team_snapshot_id = ?",
				db.Table(TableTeamSnapshot).Select("id").Where("team_id = ? and contest_id = ?", t.ID, contest.ID).SubQuery(),
			).SubQuery(),
		).SubQuery(),
		db.Table(TableProblem).Select("id").Where("contest_id = ?", contest.ID).SubQuery(),
	).Find(&hints); v.Error != nil {
		panic(v.Error)
	}
	return hints
}

func (t *Team) GetContestRank(db *gorm.DB, contest *Contest) (*ContestRank, error) {
	var rank ContestRank
	if v := db.Where("contest_id = ? and team_id = ?", contest.ID, t.ID).First(&rank); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &rank, nil
}

func (t *Team) GetSnapshotByContestID(db *gorm.DB, contestID uint) (*TeamSnapshot, error) {
	var snapshot TeamSnapshot
	if v := db.Table(TableTeamSnapshot).Select("id").Where("team_id = ? and contest_id = ?", t.ID, contestID).First(&snapshot); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &snapshot, nil
}

// Add associations
func (t *Team) AddMember(db *gorm.DB, user *User) {
	if v := db.Model(t).Association("Members").Append(user); v.Error != nil {
		panic(v.Error)
	}
}

func (t *Team) AddAdmin(db *gorm.DB, user *User) {
	if v := db.Model(t).Association("Admins").Append(user); v.Error != nil {
		panic(v.Error)
	}
}

func (t *Team) NewSnapshot(db *gorm.DB, contest *Contest, usersID []uint) (*TeamSnapshot, error) {
	if !t.IsAllMemberByUserID(db, usersID) {
		return nil, errors.ErrNotMemberOfTeam
	}
	if contest.HasTeamByTeamID(db, t.ID) {
		return nil, errors.ErrTeamParticipated
	}
	if contest.HasPlayerByUserID(db, usersID) {
		return nil, errors.ErrMemberParticipated
	}
	snapshot := TeamSnapshot{
		TeamID:    t.ID,
		ContestID: contest.ID,
	}
	if v := db.Create(&snapshot); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Model(snapshot).Association("Members").Append(GetUsersByID(db, usersID)); v.Error != nil {
		panic(v.Error)
	}
	return &snapshot, nil
}

// TeamDelete associations
func (t *Team) DeleteMember(db *gorm.DB, user *User) {
	if v := db.Model(t).Association("Members").Delete(user); v.Error == gorm.ErrRecordNotFound {
		return
	} else if v.Error != nil {
		panic(v.Error)
	}
}

func (t *Team) DeleteAdmin(db *gorm.DB, user *User) {
	if v := db.Model(t).Association("Admins").Delete(user); v.Error == gorm.ErrRecordNotFound {
		return
	} else if v.Error != nil {
		panic(v.Error)
	}
}

// Calculate points in contest
func (t *Team) CalcContestPoints(db *gorm.DB, contest *Contest) (points uint) {
	problems := t.GetSolvedProblems(db, contest)
	for _, problem := range problems {
		points += problem.Points
	}
	hints := t.GetUnlockedHints(db, contest)
	for _, hint := range hints {
		points -= hint.Cost
	}
	return points
}

// Permission check
func (t *Team) IsAllMemberByUserID(db *gorm.DB, userID []uint) bool {
	var count uint
	if v := db.Table(TableRelUserTeam).Where("team_id = ? and user_id in (?)", t.ID, userID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count == uint(len(userID))
}

func (t *Team) IsMemberByUserID(db *gorm.DB, userID uint) bool {
	var count uint
	if v := db.Table(TableRelUserTeam).Where("team_id = ? and user_id = ?", t.ID, userID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func (t *Team) IsMember(db *gorm.DB, user *User) bool {
	return t.IsMemberByUserID(db, user.ID)
}

func (t *Team) IsAdminByUserID(db *gorm.DB, userID uint) bool {
	var count uint
	if v := db.Table(TableRelTeamAdmin).Where("team_id = ? and user_id = ?", t.ID, userID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func (t *Team) IsAdmin(db *gorm.DB, user *User) bool {
	return t.IsAdminByUserID(db, user.ID)
}

// Invitation token
func (t *Team) GenerateInvitationToken(db *gorm.DB) {
	t.InvitationToken = random.RandomString(32, constants.DicLetterNumeric)
	t.Update(db)
}

type TeamSnapshot struct {
	gorm.Model
	TeamID    uint
	ContestID uint
	Members   []*User `gorm:"many2many:rel_team_snapshot_member"`
}

func (ts *TeamSnapshot) Delete(db *gorm.DB) {
	if db.NewRecord(ts) {
		return
	}
	if v := db.Model(ts).Association("Members").Clear(); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Unscoped().Delete(ts); v.Error != nil {
		panic(v.Error)
	}
}

func (ts *TeamSnapshot) GetMembers(db *gorm.DB) []User {
	var users []User
	if v := db.Where("id in (?)",
		db.Table(TableRelTeamSnapshotMember).Select("user_id").Where("team_snapshot_id = ?", ts.ID).SubQuery(),
	).Preload("Profile").Find(&users); v.Error != nil {
		panic(v.Error)
	}
	return users
}

const (
	InvitationPending = iota + 1
	InvitationAccepted
	InvitationRejected
	InvitationFailed
)

type TeamInvitation struct {
	gorm.Model
	TeamID   uint
	ToUser   uint
	FromUser uint
	Status   uint
}

func (inv *TeamInvitation) Update(db *gorm.DB) {
	if v := db.Save(inv); v.Error != nil {
		panic(v.Error)
	}
}

func (inv *TeamInvitation) Delete(db *gorm.DB) {
	if db.NewRecord(inv) {
		return
	}
	if v := db.Unscoped().Delete(inv); v.Error != nil {
		panic(v.Error)
	}
}

func (inv *TeamInvitation) setFailed(db *gorm.DB) {
	inv.Status = InvitationFailed
	inv.Update(db)
}

func (inv *TeamInvitation) Accept(db *gorm.DB) error {
	if inv.Status != InvitationPending {
		inv.setFailed(db)
		return errors.ErrProcessed
	}
	if inv.CreatedAt.Add(24 * time.Hour).Before(time.Now()) {
		inv.setFailed(db)
		return errors.ErrTimeout
	}
	if IsTeamMemberByUserID(db, inv.TeamID, inv.ToUser) {
		inv.setFailed(db)
		return errors.ErrAlreadyJoinedInTeam
	}
	addUserToTeam(db, inv.ToUser, inv.TeamID)
	inv.Status = InvitationAccepted
	inv.Update(db)
	return nil
}

func (inv *TeamInvitation) Reject(db *gorm.DB) error {
	if inv.Status != InvitationPending {
		inv.setFailed(db)
		return errors.ErrProcessed
	}
	if inv.CreatedAt.Add(24 * time.Hour).Before(time.Now()) {
		inv.setFailed(db)
		return errors.ErrTimeout
	}
	if IsTeamMemberByUserID(db, inv.TeamID, inv.ToUser) {
		inv.setFailed(db)
		return errors.ErrAlreadyJoinedInTeam
	}
	inv.Status = InvitationRejected
	inv.Update(db)
	return nil
}

const (
	ApplicationPending = iota + 1
	ApplicationAccepted
	ApplicationRejected
	ApplicationFailed
)

type TeamApplication struct {
	gorm.Model
	TeamID uint
	UserID uint
	Status uint
}

func (app *TeamApplication) Update(db *gorm.DB) {
	if v := db.Save(app); v.Error != nil {
		panic(v.Error)
	}
}

func (app *TeamApplication) Delete(db *gorm.DB) {
	if v := db.Unscoped().Delete(app); v.Error != nil {
		panic(v.Error)
	}
}

func (app *TeamApplication) setFailed(db *gorm.DB) {
	app.Status = ApplicationFailed
	app.Update(db)
}

func (app *TeamApplication) Accept(db *gorm.DB) error {
	if app.Status != ApplicationPending {
		app.setFailed(db)
		return errors.ErrProcessed
	}
	if app.CreatedAt.Add(24 * time.Hour).Before(time.Now()) {
		app.setFailed(db)
		return errors.ErrTimeout
	}
	if IsTeamMemberByUserID(db, app.TeamID, app.UserID) {
		app.setFailed(db)
		return errors.ErrAlreadyJoinedInTeam
	}
	addUserToTeam(db, app.UserID, app.TeamID)
	app.Status = ApplicationAccepted
	app.Update(db)
	return nil
}

func (app *TeamApplication) Reject(db *gorm.DB) error {
	if app.Status != ApplicationPending {
		app.setFailed(db)
		return errors.ErrProcessed
	}
	if app.CreatedAt.Add(24 * time.Hour).Before(time.Now()) {
		app.setFailed(db)
		return errors.ErrTimeout
	}
	if IsTeamMemberByUserID(db, app.TeamID, app.UserID) {
		app.setFailed(db)
		return errors.ErrAlreadyJoinedInTeam
	}
	app.Status = ApplicationRejected
	app.Update(db)
	return nil
}

func HasTeamName(db *gorm.DB, name string) bool {
	var count uint
	if v := db.Table(TableTeam).Where("name = ?", name).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func NewTeam(db *gorm.DB, name string, description string, creatorID uint) *Team {
	team := Team{
		Name:        name,
		Description: description,
		CreatorID:   creatorID,
	}
	if v := db.Create(&team); v.Error != nil {
		panic(v.Error)
	}
	return &team
}

func GetTeamByID(db *gorm.DB, id uint) (*Team, error) {
	var team Team
	if v := db.Not("is_default = 1").First(&team, id); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &team, nil
}

func GetTeams(db *gorm.DB, query string) []Team {
	var teams []Team
	if v := db.Where("id like ? or name like ?", "%"+query+"%", "%"+query+"%").Not("is_default = 1").Find(&teams); v.Error != nil {
		panic(v.Error)
	}
	return teams
}

func IsTeamMemberByUserID(db *gorm.DB, teamID uint, userID uint) bool {
	var count uint
	if v := db.Table(TableRelUserTeam).Where("team_id = ? and user_id = ?", teamID, userID).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func addUserToTeam(db *gorm.DB, userID, teamID uint) {
	if v := db.Exec("insert into `"+TableRelUserTeam+"` (`team_id`, `user_id`) values (?, ?)", teamID, userID); v.Error != nil {
		panic(v.Error)
	}
}

func pushTeamInvitationNotification(db *gorm.DB, userID uint) {
	_ = NewNotification(db, userID, "Someone wants to invite you to their team. Go to team invitation page and check.", 0)
}

func NewTeamInvitation(db *gorm.DB, team *Team, toUserID, fromUserID uint) *TeamInvitation {
	invitation := TeamInvitation{
		TeamID:   team.ID,
		ToUser:   toUserID,
		FromUser: fromUserID,
		Status:   InvitationPending,
	}
	if v := db.Save(&invitation); v.Error != nil {
		panic(v.Error)
	}
	pushTeamInvitationNotification(db, toUserID)
	return &invitation
}

func GetPendingInvitation(db *gorm.DB, teamID, userID uint) (*TeamInvitation, error) {
	var invitation TeamInvitation
	if v := db.Where(&TeamInvitation{TeamID: teamID, ToUser: userID, Status: InvitationPending}).First(&invitation); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &invitation, nil
}

func HasPendingInvitation(db *gorm.DB, teamID, userID uint) bool {
	var count uint
	if v := db.Table(TableTeamInvitation).Where(&TeamInvitation{TeamID: teamID, ToUser: userID, Status: InvitationPending}).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func GetTeamInvitations(db *gorm.DB, teamID uint) []TeamInvitation {
	var invitations []TeamInvitation
	if v := db.Where("team_id = ?", teamID).Find(&invitations); v.Error != nil {
		panic(v.Error)
	}
	return invitations
}

func pushTeamApplicationNotification(db *gorm.DB, team *Team) {
	admins := team.GetAdmins(db)
	for _, admin := range admins {
		_ = NewNotification(db, admin.ID, "Someone wants to join your team. Go to team application page and check.", 0)
	}
}

func NewTeamApplication(db *gorm.DB, team *Team, userID uint) *TeamApplication {
	application := TeamApplication{
		TeamID: team.ID,
		UserID: userID,
		Status: ApplicationPending,
	}
	if v := db.Save(&application); v.Error != nil {
		panic(v.Error)
	}
	pushTeamApplicationNotification(db, team)
	return &application
}

func GetPendingApplication(db *gorm.DB, teamID, userID uint) (*TeamApplication, error) {
	var application TeamApplication
	if v := db.Where(&TeamApplication{TeamID: teamID, UserID: userID, Status: ApplicationPending}).First(&application); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &application, nil
}

func HasPendingApplication(db *gorm.DB, teamID, userID uint) bool {
	var count uint
	if v := db.Table(TableTeamApplication).Where(&TeamApplication{TeamID: teamID, UserID: userID, Status: ApplicationPending}).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func GetTeamApplications(db *gorm.DB, teamID uint) []TeamApplication {
	var applications []TeamApplication
	if v := db.Where("team_id = ?", teamID).Find(&applications); v.Error != nil {
		panic(v.Error)
	}
	return applications
}

func GetUserInvitations(db *gorm.DB, userID uint) []TeamInvitation {
	var out []TeamInvitation
	if v := db.Where("to_user = ?", userID).Find(&out); v.Error != nil {
		panic(v.Error)
	}
	return out
}

func GetUserApplications(db *gorm.DB, userID uint) []TeamApplication {
	var out []TeamApplication
	if v := db.Where("user_id = ?", userID).Find(&out); v.Error != nil {
		panic(v.Error)
	}
	return out
}

func GetUserTeams(db *gorm.DB, userID uint, admin bool, query string) []Team {
	var teams []Team
	if admin {
		if v := db.Where("id in (?)",
			db.Table(TableRelTeamAdmin).Select("team_id").Where("user_id = ?", userID).SubQuery(),
		).Where("id like ? or name like ?", "%"+query+"%", "%"+query+"%").Not("is_default = 1").Find(&teams); v.Error != nil {
			panic(v.Error)
		}
	} else {
		if v := db.Where("id in (?)",
			db.Table(TableRelUserTeam).Select("team_id").Where("user_id = ?", userID).SubQuery(),
		).Where("id like ? or name like ?", "%"+query+"%", "%"+query+"%").Not("is_default = 1").Find(&teams); v.Error != nil {
			panic(v.Error)
		}
	}
	return teams
}
