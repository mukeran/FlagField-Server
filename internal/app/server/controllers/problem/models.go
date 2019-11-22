package problem

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/jinzhu/gorm"
)

type ReqCreate struct {
	Name        string `json:"name" validate:"max=50,min=1,required" binding:"max=50,min=1,required"`
	Description string `json:"description"`
	Alias       string `json:"alias" validate:"max=50,min=1,alphanum,required" binding:"max=50,min=1,alphanum,required"`
	Points      uint   `json:"points"`
	Type        string `json:"type"`
}

type ReqModify struct {
	Name        *string `json:"name" validate:"omitempty,max=50,min=1" binding:"omitempty,max=50,min=1"`
	Description *string `json:"description"`
	Points      *uint   `json:"points"`
	Type        *string `json:"type"`
	IsHidden    *bool   `json:"is_hidden"`
}

func (r *ReqModify) Bind(p *cpt.Problem) {
	if r.Name != nil {
		p.Name = *r.Name
	}
	if r.Description != nil {
		p.Description = *r.Description
	}
	if r.Points != nil {
		p.Points = *r.Points
	}
	if r.Type != nil {
		p.Type = *r.Type
	}
	if r.IsHidden != nil {
		p.IsHidden = *r.IsHidden
	}
}

type ReqSubmissionCreate struct {
	Flag string `json:"flag" validate:"required" binding:"required"`
}

type RespList struct {
	Response
	Problems []*RProblem `json:"problems"`
}

func BindList(db *gorm.DB, session *cpt.Session, isContestAdmin bool, problems []cpt.Problem) []*RProblem {
	var out []*RProblem
	for _, problem := range problems {
		if problem.IsHidden && !session.User.IsAdmin && !isContestAdmin {
			continue
		}
		out = append(out, &RProblem{
			Name:     problem.Name,
			Alias:    problem.Alias,
			Points:   problem.Points,
			Type:     problem.Type,
			Tags:     problem.GetTagsName(db),
			IsHidden: problem.IsHidden,
		})
	}
	return out
}

type RespCreate struct {
	Response
	ProblemAlias string `json:"problem_alias,omitempty"`
}

type RProblem struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Alias       string   `json:"alias"`
	ContestID   uint     `json:"contest_id"`
	Points      uint     `json:"points"`
	Type        string   `json:"type"`
	Tags        []string `json:"tags"`
	IsHidden    bool     `json:"is_hidden,omitempty"`
}

type RespShow struct {
	Response
	Problem *RProblem `json:"problem"`
}

type RespSubmissionCreate struct {
	Response
	Result uint `json:"result,omitempty"`
}
