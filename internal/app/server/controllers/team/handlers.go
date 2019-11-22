package team

import (
	"github.com/gin-gonic/gin"

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
		)
		defer Pack(c, &err, &response)
		query := c.DefaultQuery("query", "")
		tx := GetDB(c)
		teams := cpt.GetTeams(tx, query)
		response.Teams = BindList(tx, teams)
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
		if cpt.HasTeamName(tx, request.Name) {
			err = errors.ErrDuplicatedName
			return
		}
		team := cpt.NewTeam(tx, request.Name, request.Description, session.UserID)
		team.AddMember(tx, session.User)
		team.AddAdmin(tx, session.User)
		response.TeamID = team.ID
	}
}

func (*Handler) Show() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespShow
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		response.Team = BindTeam(tx, team)
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
		team := GetTeam(c)
		request.Bind(team)
		team.Update(tx)
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
		team := GetTeam(c)
		team.Delete(tx)
	}
}

func (*Handler) UserList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespUserList
			err      error
		)
		defer Pack(c, &err, &response)
		query := c.DefaultQuery("query", "")
		tx := GetDB(c)
		team := GetTeam(c)
		users := team.GetMembers(tx, query)
		response.Members = BindUserList(users)
	}
}

func (*Handler) UserAdd() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqUserAdd
			response RespUserAdd
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		team := GetTeam(c)
		for _, userID := range request.UsersID {
			var (
				r2   RespUserAddDeleteDetail
				err2 error
			)
			user, err2 := cpt.GetUserByID(tx, userID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
				goto end
			}
			team.AddMember(tx, user)
		end:
			r2.UserID = userID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
	}
}

func (*Handler) UserDelete() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqUserDelete
			response RespUserDelete
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		team := GetTeam(c)
		for _, userID := range request.UsersID {
			var (
				r2   RespUserAddDeleteDetail
				err2 error
			)
			user, err2 := cpt.GetUserByID(tx, userID)
			if err2 != nil {
				err = errors.ErrPartlyFailed
				goto end
			}
			if user.ID == session.UserID {
				err = errors.ErrPartlyFailed
				err2 = errors.ErrInvalidRequest
				goto end
			}
			team.DeleteMember(tx, user)
		end:
			r2.UserID = userID
			r2.Pack(&err2)
			response.Details = append(response.Details, &r2)
		}
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
		team := GetTeam(c)
		admins := team.GetAdmins(tx)
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
		team := GetTeam(c)
		session := GetSession(c)
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
			if !session.User.IsAdmin && !cpt.IsTeamMemberByUserID(tx, team.ID, user.ID) {
				err2 = errors.ErrNotMemberOfTeam
			} else {
				team.AddAdmin(tx, user)
			}
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
		session := GetSession(c)
		team := GetTeam(c)
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
			if user.ID == session.UserID {
				err = errors.ErrPartlyFailed
				err2 = errors.ErrInvalidRequest
				goto end
			}
			team.DeleteAdmin(tx, user)
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
		team := GetTeam(c)
		contests := team.GetContests(tx)
		for _, contest := range contests {
			rank, err := team.GetContestRank(tx, &contest)
			if err != nil {
				continue
			}
			response.Contests = append(response.Contests, &RStatisticContest{
				ID:     contest.ID,
				Name:   contest.Name,
				Points: rank.Points,
				Rank:   rank.Rank,
			})
		}
	}
}

func (*Handler) InvitationList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespInvitationList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		invitations := cpt.GetTeamInvitations(tx, team.ID)
		response.Invitations = BindInvitationList(invitations)
	}
}

func (*Handler) InvitationNew() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqInvitationNewAndCancel
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		team := GetTeam(c)
		session := GetSession(c)
		if cpt.IsTeamMemberByUserID(tx, team.ID, request.UserID) {
			err = errors.ErrAlreadyJoinedInTeam
			return
		}
		if cpt.HasPendingInvitation(tx, team.ID, request.UserID) {
			err = errors.ErrAlreadyInvited
			return
		}
		_ = cpt.NewTeamInvitation(tx, team, request.UserID, session.UserID)
	}
}

func (*Handler) InvitationCancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqInvitationNewAndCancel
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		team := GetTeam(c)
		invitation, err := cpt.GetPendingInvitation(tx, team.ID, request.UserID)
		if err != nil {
			err = errors.ErrNotInvited
			return
		}
		invitation.Delete(tx)
	}
}

func (*Handler) InvitationTokenShow() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespInvitationToken
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		if team.InvitationToken == "" {
			team.GenerateInvitationToken(tx)
		}
		response.Token = team.InvitationToken
	}
}

func (*Handler) InvitationTokenRefresh() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespInvitationToken
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		team.GenerateInvitationToken(tx)
		response.Token = team.InvitationToken
	}
}

func (*Handler) InvitationAcceptByToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqInvitationAcceptByToken
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		session := GetSession(c)
		team := GetTeam(c)
		if team.InvitationToken != "" && request.Token == team.InvitationToken {
			if cpt.IsTeamMemberByUserID(tx, team.ID, session.UserID) {
				err = errors.ErrAlreadyJoinedInTeam
				return
			}
			team.AddMember(tx, session.User)
		} else {
			err = errors.ErrInvalidInvitationToken
			return
		}
	}
}

func (*Handler) InvitationAccept() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		session := GetSession(c)
		team := GetTeam(c)
		invitation, err := cpt.GetPendingInvitation(tx, team.ID, session.UserID)
		if err != nil {
			err = errors.ErrNotInvited
			return
		}
		err = invitation.Accept(tx)
	}
}

func (*Handler) InvitationReject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		session := GetSession(c)
		invitation, err := cpt.GetPendingInvitation(tx, team.ID, session.UserID)
		if err != nil {
			err = errors.ErrNotInvited
			return
		}
		err = invitation.Reject(tx)
	}
}

func (*Handler) ApplicationList() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response RespApplicationList
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		applications := cpt.GetTeamApplications(tx, team.ID)
		response.Applications = BindApplicationList(applications)
	}
}

func (*Handler) ApplicationNew() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		session := GetSession(c)
		if cpt.IsTeamMemberByUserID(tx, team.ID, session.UserID) {
			err = errors.ErrAlreadyJoinedInTeam
			return
		}
		if cpt.HasPendingApplication(tx, team.ID, session.UserID) {
			err = errors.ErrAlreadyApplied
			return
		}
		_ = cpt.NewTeamApplication(tx, team, session.UserID)
	}
}

func (*Handler) ApplicationCancel() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		tx := GetDB(c)
		team := GetTeam(c)
		session := GetSession(c)
		application, err := cpt.GetPendingApplication(tx, team.ID, session.UserID)
		if err != nil {
			err = errors.ErrNoSuchApplication
			return
		}
		application.Delete(tx)
	}
}

func (*Handler) ApplicationAccept() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqApplicationAcceptAndReject
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		team := GetTeam(c)
		application, err := cpt.GetPendingApplication(tx, team.ID, request.UserID)
		if err != nil {
			err = errors.ErrNoSuchApplication
			return
		}
		err = application.Accept(tx)
	}
}

func (*Handler) ApplicationReject() gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			request  ReqApplicationAcceptAndReject
			response Response
			err      error
		)
		defer Pack(c, &err, &response)
		err = BindJSON(c, &request)
		if err != nil {
			return
		}
		tx := GetDB(c)
		team := GetTeam(c)
		application, err := cpt.GetPendingApplication(tx, team.ID, request.UserID)
		if err != nil {
			err = errors.ErrNoSuchApplication
			return
		}
		err = application.Reject(tx)
	}
}
