package resource

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	"time"

	"github.com/jinzhu/gorm"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

type ReqModify struct {
	ExpiredAt    *time.Time `json:"expired_at"`
	IsHidden     *bool      `json:"is_hidden"`
	ContestID    *uint      `json:"contest_id"`
	ProblemAlias *string    `json:"problem_alias"`
}

func (r *ReqModify) Bind(db *gorm.DB, session *cpt.Session, res *cpt.Resource) (err error) {
	if r.ExpiredAt != nil {
		res.ExpiredAt = *r.ExpiredAt
	}
	if r.IsHidden != nil {
		res.IsHidden = *r.IsHidden
	}
	var (
		contest *cpt.Contest
		problem *cpt.Problem
	)
	if r.ContestID != nil {
		contest, err = cpt.GetContestByID(db, *r.ContestID)
		if err != nil {
			return
		}
		if !session.User.IsAdmin && !session.User.IsContestAdmin(db, contest) {
			return errors.ErrNotFound
		}
		res.ContestID = *r.ContestID
	}
	if r.ProblemAlias != nil {
		if r.ContestID == nil {
			return errors.ErrInvalidRequest
		}
		problem, err = cpt.GetProblemByAliasAndContestID(db, *r.ProblemAlias, contest.ID)
		if err != nil {
			return
		}
		res.ProblemID = problem.ID
	}
	return
}

type RespResource struct {
	UUID        string    `json:"uuid"`
	Name        string    `json:"name"`
	ContentType string    `json:"content_type"`
	ExpiredAt   time.Time `json:"expired_at"`
	IsHidden    bool      `json:"is_hidden"`
	ContestID   uint      `json:"contest_id"`
	ProblemID   uint      `json:"problem_id"`
	CreatedAt   time.Time `json:"created_at"`
}

func BindResource(res *cpt.Resource) *RespResource {
	return &RespResource{
		UUID:        res.UUID,
		Name:        res.Name,
		ContentType: res.ContentType,
		ExpiredAt:   res.ExpiredAt,
		IsHidden:    res.IsHidden,
		ContestID:   res.ContestID,
		ProblemID:   res.ProblemID,
		CreatedAt:   res.CreatedAt,
	}
}

type RespList struct {
	Response
	Resources []*RespResource `json:"resources"`
}

func BindList(resources []cpt.Resource) []*RespResource {
	var out []*RespResource
	for _, res := range resources {
		out = append(out, BindResource(&res))
	}
	return out
}

type RespUpload struct {
	Response
	UUID string `json:"uuid"`
}
