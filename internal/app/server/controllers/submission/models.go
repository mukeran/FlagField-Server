package submission

import (
	. "github.com/FlagField/FlagField-Server/internal/app/server/controllers"
	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type RSubmission struct {
	ID        uint   `json:"id"`
	Flag      string `json:"flag"`
	Result    uint   `json:"result"`
	ProblemID uint   `json:"problem_id"`
	UserID    uint   `json:"user_id"`
}

func BindSubmission(submission *cpt.Submission) *RSubmission {
	return &RSubmission{
		ID:        submission.ID,
		Flag:      submission.Flag,
		Result:    submission.Result,
		ProblemID: submission.ProblemID,
		UserID:    submission.CreatorID,
	}
}

type RespList struct {
	Response
	Submissions []*RSubmission `json:"submissions"`
}

func BindList(submissions []cpt.Submission) []*RSubmission {
	var out []*RSubmission
	for _, submission := range submissions {
		out = append(out, BindSubmission(&submission))
	}
	return out
}
