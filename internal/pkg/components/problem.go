package components

import (
	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	TableProblem          = "problem"
	TableProblemTag       = "problem_tag"
	TableRelProblemTag    = "rel_problem_tag"
	TableRelProblemSolver = "rel_problem_solver"
)

type Problem struct {
	gorm.Model
	Name        string `gorm:"size:50;not null"`
	Description string `gorm:"type:longtext"`
	Alias       string `gorm:"size:50;not null"`
	ContestID   uint
	Points      uint
	Type        string `gorm:"size:50;not null"`
	IsHidden    bool
	CreatorID   uint
	Flags       []*Flag
	Hints       []*Hint
	Resources   []*Resource
	Submissions []*Submission
	Solvers     []*User       `gorm:"many2many:rel_problem_solver"`
	Tags        []*ProblemTag `gorm:"many2many:rel_problem_tag"`
}

func (p *Problem) Update(db *gorm.DB) {
	if v := db.Save(p); v.Error != nil {
		panic(v.Error)
	}
}

func (p *Problem) AddFlagWithSettingsMap(db *gorm.DB, ft FlagType, m map[string]interface{}) (int, error) {
	fs, err := NewFlagSettingsFromMap(ft, m)
	if err != nil {
		return 0, err
	}
	return p.AddFlag(db, ft, fs), nil
}

func (p *Problem) AddFlag(db *gorm.DB, ft FlagType, fs FlagSettings) int {
	if v := db.Model(p).Association("Flags").Append(&Flag{
		Type:     string(ft),
		Settings: fs,
	}); v.Error != nil {
		panic(v.Error)
	}
	return db.Model(p).Association("Flags").Count()
}

func (p *Problem) AddHint(db *gorm.DB, cost uint, content string) int {
	if v := db.Model(p).Association("Hints").Append(&Hint{
		Cost:    cost,
		Content: content,
	}); v.Error != nil {
		panic(v.Error)
	}
	return db.Model(p).Association("Hints").Count()
}

func (p *Problem) AddSolver(db *gorm.DB, u *User) {
	if v := db.Model(p).Association("Solvers").Append(u); v.Error != nil {
		panic(v.Error)
	}
}

func (p *Problem) AddResource(db *gorm.DB, res *Resource) {
	if v := db.Model(p).Association("Resources").Append(res); v.Error != nil {
		panic(v.Error)
	}
}

func (p *Problem) AddTag(db *gorm.DB, tag string) {
	var t ProblemTag
	t.Name = tag
	if v := db.Where(t).FirstOrCreate(&t); v.Error != nil {
		panic(v.Error)
	}
	if v := db.Model(p).Association("Tags").Append(t); v.Error != nil {
		panic(v.Error)
	}
}

func (p *Problem) Delete(db *gorm.DB) {
	if db.NewRecord(p) {
		return
	}
	if v := db.Unscoped().Delete(p); v.Error != nil {
		panic(v.Error)
	}
}

func (p *Problem) GetFlagByOrder(db *gorm.DB, order int) (*Flag, error) {
	if order <= 0 {
		return nil, errors.ErrOutOfRange
	}
	var flag Flag
	if v := db.Table(TableFlag).Where("problem_id = ?", p.ID).Offset(order - 1).Limit(1).First(&flag); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrOutOfRange
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &flag, nil
}

func (p *Problem) GetFlags(db *gorm.DB) []Flag {
	var flags []Flag
	if v := db.Model(p).Association("Flags").Find(&flags); v.Error != nil {
		panic(v.Error)
	}
	return flags
}

func (p *Problem) GetHintByOrder(db *gorm.DB, order int) (*Hint, error) {
	if order <= 0 {
		return nil, errors.ErrOutOfRange
	}
	var hint Hint
	if v := db.Table(TableHint).Where("problem_id = ?", p.ID).Offset(order - 1).Limit(1).First(&hint); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrOutOfRange
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &hint, nil
}

func (p *Problem) GetHints(db *gorm.DB) []Hint {
	var hints []Hint
	if v := db.Model(p).Association("Hints").Find(&hints); v.Error != nil {
		panic(v.Error)
	}
	return hints
}

func (p *Problem) GetSolvers(db *gorm.DB) []User {
	var solvers []User
	if v := db.Model(p).Association("Solvers").Find(&solvers); v.Error != nil {
		panic(v.Error)
	}
	return solvers
}

func (p *Problem) GetResources(db *gorm.DB) []Resource {
	var resources []Resource
	if v := db.Model(p).Association("Resources").Find(&resources); v.Error != nil {
		panic(v.Error)
	}
	return resources
}

func (p *Problem) GetTags(db *gorm.DB) []ProblemTag {
	var tags []ProblemTag
	if v := db.Model(p).Association("Tags").Find(&tags); v.Error != nil {
		panic(v.Error)
	}
	return tags
}

func (p *Problem) GetTagsName(db *gorm.DB) []string {
	var tagsName []string
	if v := db.Table(TableProblemTag).Where("id in (?)",
		db.Table(TableRelProblemTag).Select("problem_tag_id").Where("problem_id = ?", p.ID).SubQuery(),
	).Pluck("name", &tagsName); v.Error != nil {
		panic(v.Error)
	}
	return tagsName
}

func (p *Problem) DeleteTag(db *gorm.DB, tagName string) {
	var tag ProblemTag
	if v := db.Where("name = ?", tagName).First(&tag); v.Error == gorm.ErrRecordNotFound {
		return
	} else if v.Error != nil {
		panic(v.Error)
	}
	if v := db.Model(p).Association("Tags").Delete(tag); v.Error == gorm.ErrRecordNotFound {
		return
	} else if v.Error != nil {
		panic(v.Error)
	}
}

type ProblemTag struct {
	gorm.Model
	Name     string     `gorm:"not null"`
	Problems []*Problem `gorm:"many2many:rel_problem_tag"`
}

func NewProblem(db *gorm.DB, name string, description string, alias string, contestID uint, points uint, typ string, creatorID uint) *Problem {
	problem := &Problem{
		Name:        name,
		Description: description,
		Alias:       alias,
		ContestID:   contestID,
		Points:      points,
		Type:        typ,
		CreatorID:   creatorID,
		IsHidden:    true,
	}
	if v := db.Create(problem); v.Error != nil {
		panic(v.Error)
	}
	return problem
}

func GetProblemByAliasAndContestID(db *gorm.DB, alias string, contestID uint) (*Problem, error) {
	var problem Problem
	if v := db.Table(TableProblem).Where("alias = ? and contest_id = ?", alias, contestID).First(&problem); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &problem, nil
}

func HasProblemAlias(db *gorm.DB, contestID uint, alias string) bool {
	var count int
	if v := db.Table(TableProblem).Where("contest_id = ? and alias = ?", contestID, alias).Count(&count); v.Error != nil {
		panic(v.Error)
	}
	return count >= 1
}

func GetProblemID(db *gorm.DB, contestID uint, alias string) (uint, error) {
	var problem Problem
	if v := db.Table(TableProblem).Select("id").Where("alias = ? and contest_id = ?", alias, contestID).First(&problem); v.Error == gorm.ErrRecordNotFound {
		return 0, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return problem.ID, nil
}
