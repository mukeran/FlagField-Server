package hint

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
	"github.com/jinzhu/gorm"
)

type ReqCreate struct {
	Cost    uint   `json:"cost"`
	Content string `json:"content" validate:"required" binding:"required"`
}

type ReqModify struct {
	Cost    *uint   `json:"cost"`
	Content *string `json:"content" validate:"omitempty,min=1" binding:"omitempty,min=1"`
}

func (r *ReqModify) Bind(h *cpt.Hint) {
	if r.Cost != nil {
		h.Cost = *r.Cost
	}
	if r.Content != nil {
		h.Content = *r.Content
	}
}

type RespCreate struct {
	Response
	HintOrder int `json:"hint_order,omitempty"`
}

type RHint struct {
	Cost       uint   `json:"cost"`
	Content    string `json:"content,omitempty"`
	IsUnlocked bool   `json:"is_unlocked"`
}

func BindHint(hint *cpt.Hint, isUnlocked bool) *RHint {
	out := &RHint{
		Cost:       hint.Cost,
		Content:    hint.Content,
		IsUnlocked: isUnlocked,
	}
	if !isUnlocked {
		out.Content = ""
	}
	return out
}

type RespList struct {
	Response
	Hints []*RHint `json:"hints"`
}

func BindList(db *gorm.DB, contest *cpt.Contest, problem *cpt.Problem, session *cpt.Session, hints []cpt.Hint) ([]*RHint, error) {
	var out []*RHint
	for _, hint := range hints {
		out = append(out, BindHint(&hint, hint.IsUnlocked(db, contest, problem, session.User)))
	}
	return out, nil
}

type RespShow struct {
	Response
	Hint *RHint `json:"hint"`
}
