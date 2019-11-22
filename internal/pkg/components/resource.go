package components

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/constants"
	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	TableResource = "resource"
)

type Resource struct {
	gorm.Model
	UUID        string `gorm:"unique;not null"`
	Name        string `gorm:"size:255;not null"`
	Path        string `gorm:"size:4096;not null"`
	ContentType string `gorm:"not null"`
	ExpiredAt   time.Time
	IsHidden    bool
	CreatorID   uint
	ContestID   uint
	ProblemID   uint
}

func (res *Resource) Update(db *gorm.DB) {
	if v := db.Save(res); v.Error != nil {
		panic(v.Error)
	}
}

func (res *Resource) Delete(db *gorm.DB) {
	if db.NewRecord(res) {
		return
	}
	err := os.Remove(res.Path)
	if err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if v := db.Unscoped().Delete(res); v.Error != nil {
		panic(v.Error)
	}
}

// type ResourceLog struct {
// 	gorm.Model
// 	Resource     Resource `gorm:"foreignkey:ResourceID"`
// 	ResourceID   uint
// 	Action       string  `gorm:"required"`
// 	Session      Session `gorm:"foreignkey:SessionNumID"`
// 	SessionNumID uint
// }

func GenerateResourceUUID() string {
	id, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	return id.String()
}

func newResourcePrepare(baseDir string, name string, contentType string, userID uint) (*Resource, *os.File) {
	res := &Resource{
		UUID:        GenerateResourceUUID(),
		Name:        name,
		ContentType: contentType,
		ExpiredAt:   constants.MaxTime,
		CreatorID:   userID,
	}
	res.Path = fmt.Sprintf("%s%s", baseDir, res.UUID)
	dest, err := os.Create(res.Path)
	if err != nil {
		panic(err)
	}
	return res, dest
}

func NewResourceWithReader(db *gorm.DB, baseDir string, name string, contentType string, src io.Reader, userID uint) *Resource {
	res, dest := newResourcePrepare(baseDir, name, contentType, userID)
	defer dest.Close()
	_, err := io.Copy(dest, src)
	if err != nil {
		panic(err)
	}
	if v := db.Create(res); v.Error != nil {
		panic(v.Error)
	}
	return res
}

func GetResourceByUUID(db *gorm.DB, resourceUUID string) (*Resource, error) {
	var res Resource
	if v := db.Where("uuid = ?", resourceUUID).First(&res); v.Error == gorm.ErrRecordNotFound {
		return nil, errors.ErrNotFound
	} else if v.Error != nil {
		panic(v.Error)
	}
	return &res, nil
}

func GetResources(db *gorm.DB) []Resource {
	var resources []Resource
	if v := db.Find(&resources); v.Error != nil {
		panic(v.Error)
	}
	return resources
}

func GetContestResources(db *gorm.DB, contestID uint) []Resource {
	var resources []Resource
	if v := db.Where("contest_id = ?", contestID).Find(&resources); v.Error != nil {
		panic(v.Error)
	}
	return resources
}

func GetProblemResources(db *gorm.DB, problemID uint) []Resource {
	var resources []Resource
	if v := db.Where("problem_id = ?", problemID).Find(&resources); v.Error != nil {
		panic(v.Error)
	}
	return resources
}
