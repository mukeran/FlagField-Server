package templates

import (
	"github.com/jinzhu/gorm"

	cpt "github.com/FlagField/FlagField-Server/internal/pkg/components"
)

type Initial struct {
}

func (t *Initial) Execute(db *gorm.DB) error {
	db.SingularTable(true)
	db = db.Begin()
	if v := db.AutoMigrate(
		&cpt.Config{},
		&cpt.Contest{}, &cpt.ContestNotification{}, &cpt.ContestRank{},
		&cpt.Flag{},
		&cpt.Hint{},
		&cpt.Notification{},
		&cpt.Problem{}, &cpt.ProblemTag{},
		&cpt.Resource{},
		&cpt.Session{},
		&cpt.Submission{},
		&cpt.Team{}, &cpt.TeamSnapshot{}, &cpt.TeamInvitation{}, &cpt.TeamApplication{},
		&cpt.User{}, &cpt.UserProfile{},
	); v.Error != nil {
		db.Rollback()
		return v.Error
	}
	/* Create default contest */
	if db.Where("name = ?", cpt.DefaultContest.Name).First(&cpt.DefaultContest).RecordNotFound() {
		if v := db.Create(&cpt.DefaultContest); v.Error != nil {
			panic(v.Error)
		}
	}
	/* Create default user */
	if db.Where("username = ?", cpt.DefaultUser.Username).First(&cpt.DefaultUser).RecordNotFound() {
		if v := db.Create(&cpt.DefaultUser); v.Error != nil {
			panic(v.Error)
		}
	}
	if v := db.Commit(); v.Error != nil {
		db.Rollback()
		return v.Error
	}
	return nil
}
