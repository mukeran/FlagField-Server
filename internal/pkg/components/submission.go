package components

import (
	"github.com/jinzhu/gorm"
)

const (
	Correct = iota + 1
	Wrong
	TableSubmission = "submission"
)

type Submission struct {
	gorm.Model
	Flag      string `gorm:"type:longtext"`
	Result    uint
	ProblemID uint
	CreatorID uint
}

func NewSubmission(db *gorm.DB, flag string, problem *Problem, user *User) *Submission {
	s := &Submission{
		Flag:      flag,
		ProblemID: problem.ID,
		CreatorID: user.ID,
	}
	flags := problem.GetFlags(db)
	s.Result = Wrong
	for _, f := range flags {
		if f.Check(flag) {
			s.Result = Correct
			break
		}
	}
	if v := db.Create(s); v.Error != nil {
		panic(v.Error)
	}
	if s.Result == Correct {
		problem.AddSolver(db, user)
	}
	return s
}

func GetSubmissions(db *gorm.DB) []Submission {
	var submissions []Submission
	if v := db.Find(&submissions); v.Error != nil {
		panic(v.Error)
	}
	return submissions
}

func GetSubmissionCount(db *gorm.DB) uint {
	var count uint
	if v := db.Table(TableSubmission).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count
}
