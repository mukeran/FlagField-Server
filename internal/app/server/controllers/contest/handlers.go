package contest

import (
	teamCtl "github.com/FlagField/FlagField-Server/internal/app/server/controllers/team"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"

	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"github.com/FlagField/FlagField-Server/internal/app/server/controllers/admin"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type Handler struct {
}

func (*Handler) List() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespList
			err      error
			contests []cpt.Contest
		)
		defer Pack(c, &err, &response)
		_status, hasStatus := c.GetQuery("status")
		tx := GetDB(c)
		session := GetSession(c)
		if hasStatus {
			status, err := strconv.ParseUint(_status, 10, 64)
			if err != nil {
				err = errors.ErrInvalidRequest
				return
			}
			const (
				PENDING = iota
				RUNNING
				ENDED // 2
			)
			if !(status >= 0 && status <= 2) {
				err = errors.ErrInvalidRequest
				return
			}
			contests = cpt.GetContestsFilterByStatus(tx, uint(status))
		} else {
			contests = cpt.GetContests(tx)
		}
		response.Contests = BindContestList(tx, session, contests)
	}
}

func (*Handler) Create() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqCreate
			response RespCreate
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		contest := cpt.NewContest(tx, request.Name, request.Description, request.StartTime, request.EndTime, request.Access, session.UserID)
		contest.AddAdmin(tx, session.User)
		response.ContestID = contest.ID
	}
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespShow
			err      error
		)
		defer Pack(c, &err, &response)
		contest := GetContest(c)
		response.Contest = BindContest(contest)
	}
}

func (*Handler) Modify() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqModify
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		if contest.ID == cpt.DefaultContest.ID {
			err = errors.ErrDangerousOperation
			return
		}
		request.Bind(contest)
		contest.Update(tx)
	}
}

func (*Handler) SingleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		contest := GetContest(c)
		if contest.ID == cpt.DefaultContest.ID {
			err = errors.ErrDangerousOperation
			return
		}
		contest.Delete(tx)
	}
}

func (*Handler) TeamList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespTeamList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		contest := GetContest(c)
		response.Teams = BindTeamList(tx, contest.GetTeamSnapshots(tx))
	}
}

func (*Handler) TeamAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqTeamAdd
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		if contest.ID == cpt.DefaultContest.ID {
			err = errors.ErrDangerousOperation
			return
		}
		if contest.EndTime.Before(time.Now()) {
			err = errors.ErrContestEnded
			return
		}
		team, err := cpt.GetTeamByID(tx, request.TeamID)
		if err != nil {
			return
		}
		session := GetSession(c)
		if !session.User.IsAdmin && !contest.IsAdmin(tx, session.User) && !team.IsAdmin(tx, session.User) {
			err = errors.ErrNotAdminOfTeam
			return
		}
		err = contest.AddTeam(tx, team, request.MembersID)
	}
}

func (*Handler) TeamDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqTeamDelete
			response RespTeamDelete
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		if contest.ID == cpt.DefaultContest.ID {
			err = errors.ErrDangerousOperation
			return
		}
		for _, teamID := range request.TeamsID {
			var (
				r2   RespTeamDeleteDetail
				err2 error
			)
			err2 = contest.DeleteTeamByID(tx, teamID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
			}
			r2.TeamID = teamID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
	}
}

func (*Handler) TeamCurrentShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespTeamCurrentShow
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		contest := GetContest(c)
		session := GetSession(c)
		team, err := contest.GetUserTeam(tx, session.User)
		if err != nil {
			return
		}
		ts, _ := team.GetSnapshotByContestID(tx, contest.ID)
		problems := team.GetSolvedProblems(tx, contest)
		for _, problem := range problems {
			var submission cpt.Submission
			if v := tx.Table(cpt.TableSubmission).Where("problem_id = ? and creator_id in (?) and result = 1", problem.ID,
				tx.Table(cpt.TableRelTeamSnapshotMember).Select("user_id").Where("team_snapshot_id = ?", ts.ID).SubQuery(),
			).Order("created_at").First(&submission); v.Error != nil {
				panic(v.Error)
			}
			response.Submissions = append(response.Submissions, &RTeamSubmission{
				ProblemAlias: problem.Alias,
				UserID:       submission.CreatorID,
				SolvedTime:   submission.CreatedAt,
			})
		}
		response.IsAdmin = team.IsAdmin(tx, session.User)
		response.Team = &RTeam{
			TeamID:  team.ID,
			Name:    team.Name,
			Members: teamCtl.BindUserList(ts.GetMembers(tx)),
		}
	}
}

func (*Handler) TeamCurrentLeave() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		contest := GetContest(c)
		if contest.ID == cpt.DefaultContest.ID {
			err = errors.ErrDangerousOperation
			return
		}
		if contest.StartTime.Before(time.Now()) {
			err = errors.ErrContestStarted
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		team, err := contest.GetUserTeam(tx, session.User)
		if err != nil {
			return
		}
		if !team.IsAdmin(tx, session.User) {
			err = errors.ErrNotAdminOfTeam
			return
		}
		err = contest.DeleteTeamByID(tx, team.ID)
	}
}

func (*Handler) AdminList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response admin.RespList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		contest := GetContest(c)
		admins := contest.GetAdmins(tx)
		response.Admins = admin.BindList(admins)
	}
}

func (*Handler) AdminAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  admin.ReqAdd
			response admin.RespAdd
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		for _, userID := range request.UsersID {
			var (
				r2   admin.RespAddDeleteDetail
				err2 error
			)
			user, err2 := cpt.GetUserByID(tx, userID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
				goto end
			}
			contest.AddAdmin(tx, user)
		end:
			r2.UserID = userID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
	}
}

func (*Handler) AdminDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  admin.ReqDelete
			response admin.RespDelete
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		for _, userID := range request.UsersID {
			var (
				r2   admin.RespAddDeleteDetail
				err2 error
			)
			user, err2 := cpt.GetUserByID(tx, userID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
				goto end
			}
			contest.DeleteAdmin(tx, user)
		end:
			r2.UserID = userID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
	}
}

func (*Handler) StatisticShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespStatisticShow
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		contest := GetContest(c)
		teams := contest.GetTeams(tx)
		for _, team := range teams {
			statisticTeam := RStatisticTeam{
				ID:     team.ID,
				Name:   team.Name,
				Points: team.CalcContestPoints(tx, contest),
			}
			problems := team.GetSolvedProblems(tx, contest)
			statisticTeam.SolvedProblems = uint(len(problems))
			for _, problem := range problems {
				var submission cpt.Submission
				if v := tx.Table(cpt.TableSubmission).Select("created_at").Where("problem_id = ? and creator_id in (?) and result = 1", problem.ID,
					tx.Table(cpt.TableRelTeamSnapshotMember).Select("user_id").Where("team_snapshot_id = ?",
						tx.Table(cpt.TableTeamSnapshot).Select("id").Where("team_id = ? and contest_id = ?", team.ID, contest.ID).SubQuery(),
					).SubQuery(),
				).Order("created_at").First(&submission); v.Error != nil {
					panic(v.Error)
				}
				response.Submissions = append(response.Submissions, &RStatisticSubmission{
					TeamID:       team.ID,
					ProblemAlias: problem.Alias,
					Points:       problem.Points,
					SolvedTime:   submission.CreatedAt,
				})
			}
			response.Teams = append(response.Teams, &statisticTeam)
		}
	}
}

func (*Handler) NotificationList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespNotificationList
			err      error
		)
		defer Pack(c, &err, &response)
		offset, err1 := strconv.ParseUint(c.DefaultQuery("offset", "0"), 10, 32)
		from, err2 := time.Parse(time.RFC3339, c.DefaultQuery("from", time.Date(2000, 1, 1, 0, 0, 0, 0, time.Local).Format(time.RFC3339)))
		if err1 != nil || err2 != nil {
			err = errors.ErrInvalidRequest
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		notifications := contest.GetNotifications(tx, from, uint(offset))
		response.RequestTime = time.Now()
		response.Notifications = BindNotificationList(notifications)
	}
}

func (*Handler) NotificationCreate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqNotificationCreate
			response RespNotificationCreate
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		contest := GetContest(c)
		response.NotificationOrder = contest.AddNotification(tx, request.Content)
	}
}

func (*Handler) NotificationSingleDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		notification := GetNotification(c)
		notification.Delete(tx)
	}
}
