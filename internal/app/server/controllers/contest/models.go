package contest

import (
	"time"

	"github.com/jinzhu/gorm"

	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/team"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type ReqCreate struct {
	Name        string    `json:"name" validate:"max=50,min=1,required" binding:"max=50,min=1,required"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time" validate:"required" binding:"required"`
	EndTime     time.Time `json:"end_time" validate:"required" binding:"required"`
	Access      uint      `json:"access"`
}

type ReqModify struct {
	Name        *string    `json:"name" validate:"omitempty,max=50,min=1" binding:"omitempty,max=50,min=1"`
	Description *string    `json:"description"`
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	Access      *uint      `json:"access"`
	IsHidden    *bool      `json:"is_hidden"`
}

func (r *ReqModify) Bind(c *cpt.Contest) {
	if r.Name != nil {
		c.Name = *r.Name
	}
	if r.Description != nil {
		c.Description = *r.Description
	}
	if r.StartTime != nil {
		c.StartTime = *r.StartTime
	}
	if r.EndTime != nil {
		c.EndTime = *r.EndTime
	}
	if r.Access != nil {
		c.Access = *r.Access
	}
	if r.IsHidden != nil {
		c.IsHidden = *r.IsHidden
	}
}

type ReqTeamAdd struct {
	TeamID    uint   `json:"team_id"`
	MembersID []uint `json:"members_id"`
}

type ReqTeamDelete struct {
	TeamsID []uint `json:"teams_id"`
}

type ReqNotificationCreate struct {
	Content string `json:"content" validate:"required" binding:"required"`
}

type RespCreate struct {
	Response
	ContestID uint `json:"contest_id,omitempty"`
}

type RespContest struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Access      uint      `json:"access"`
	IsHidden    bool      `json:"is_hidden,omitempty"`
}

func BindContest(contest *cpt.Contest) *RespContest {
	return &RespContest{
		ID:          contest.ID,
		Name:        contest.Name,
		Description: contest.Description,
		StartTime:   contest.StartTime,
		EndTime:     contest.EndTime,
		Access:      contest.Access,
		IsHidden:    contest.IsHidden,
	}
}

type RespList struct {
	Response
	Contests []*RespContest `json:"contests"`
}

func BindContestList(db *gorm.DB, session *cpt.Session, contests []cpt.Contest) []*RespContest {
	var outputContests []*RespContest
	for _, contest := range contests {
		if contest.ID == cpt.DefaultContest.ID {
			continue
		}
		if contest.IsHidden && !session.User.HasContestAccess(db, &contest) {
			continue
		}
		outputContests = append(outputContests, BindContest(&contest))
	}
	return outputContests
}

type RespShow struct {
	Response
	Contest *RespContest `json:"contest,omitempty"`
}

type RTeam struct {
	TeamID  uint              `json:"team_id"`
	Name    string            `json:"name"`
	Members []*team.RTeamUser `json:"members"`
}

func BindTeam(db *gorm.DB, ts *cpt.TeamSnapshot) *RTeam {
	t, err := cpt.GetTeamByID(db, ts.TeamID)
	if err != nil {
		t = &cpt.Team{Name: "~error"}
	}
	return &RTeam{
		TeamID:  ts.TeamID,
		Name:    t.Name,
		Members: team.BindUserList(ts.GetMembers(db)),
	}
}

type RespTeamList struct {
	Response
	Teams []*RTeam `json:"teams"`
}

func BindTeamList(db *gorm.DB, snapshots []cpt.TeamSnapshot) []*RTeam {
	var out []*RTeam
	for _, snapshot := range snapshots {
		out = append(out, BindTeam(db, &snapshot))
	}
	return out
}

type RespTeamDeleteDetail struct {
	Response
	TeamID uint `json:"team_id"`
}

type RespTeamDelete struct {
	Response
	Details []*RespTeamDeleteDetail `json:"details"`
}

type RTeamSubmission struct {
	ProblemAlias string    `json:"problem_alias"`
	UserID       uint      `json:"user_id"`
	SolvedTime   time.Time `json:"solved_time"`
}

type RespTeamCurrentShow struct {
	Response
	Team        *RTeam             `json:"team,omitempty"`
	IsAdmin     bool               `json:"is_admin"`
	Submissions []*RTeamSubmission `json:"submissions,omitempty"`
}

type RStatisticTeam struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Points         uint   `json:"points"`
	SolvedProblems uint   `json:"solved_problems"`
}

type RStatisticSubmission struct {
	TeamID       uint      `json:"team_id"`
	ProblemAlias string    `json:"problem_alias"`
	Points       uint      `json:"points"`
	SolvedTime   time.Time `json:"solved_time"`
}

type RespStatisticShow struct {
	Response
	Teams       []*RStatisticTeam       `json:"teams"`
	Submissions []*RStatisticSubmission `json:"submissions"`
}

type RNotification struct {
	Content       string    `json:"content"`
	PublishedTime time.Time `json:"published_time"`
}

func BindNotification(n *cpt.ContestNotification) *RNotification {
	return &RNotification{
		Content:       n.Content,
		PublishedTime: n.CreatedAt,
	}
}

type RespNotificationList struct {
	Response
	RequestTime   time.Time        `json:"request_time,omitempty"`
	Notifications []*RNotification `json:"notifications"`
}

func BindNotificationList(ns []cpt.ContestNotification) []*RNotification {
	var out []*RNotification
	for _, n := range ns {
		out = append(out, BindNotification(&n))
	}
	return out
}

type RespNotificationCreate struct {
	Response
	NotificationOrder int `json:"notification_order"`
}
