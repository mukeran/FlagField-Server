package components

import (
	"encoding/json"
	"strings"

	"github.com/goinggo/mapstructure"
	"github.com/jinzhu/gorm"

	"github.com/FlagField/FlagField-Server/internal/pkg/errors"
)

const (
	FlagTypeStatic = "static"
	TableFlag      = "flag"
)

type FlagType string

type Flag struct {
	gorm.Model
	ProblemID     uint
	Type          string       `gorm:"not null"`
	SettingsPlain []byte       `gorm:"column:settings,type:longtext"`
	Settings      FlagSettings `gorm:"-"`
}

func (f *Flag) BeforeSave() (err error) {
	if f.Settings == nil || !f.Settings.IsValid() {
		err = errors.ErrInvalidFlagSettings
	} else {
		f.Type = f.Settings.Type()
		f.SettingsPlain = f.Settings.Plain()
	}
	return
}

func (f *Flag) AfterFind() (err error) {
	switch f.Type {
	case FlagTypeStatic:
		f.Settings = &FlagStatic{}
	}
	err = f.Settings.FromPlain(f.SettingsPlain)
	return
}

func (f *Flag) Check(s string) bool {
	if f.Settings == nil || !f.Settings.IsValid() {
		return false
	}
	return f.Settings.Check(s)
}

func (f *Flag) Update(db *gorm.DB) {
	if v := db.Save(f); v.Error != nil {
		panic(v.Error)
	}
}

func (f *Flag) Delete(db *gorm.DB) {
	if db.NewRecord(f) {
		return
	}
	if v := db.Unscoped().Delete(f); v.Error != nil {
		panic(v.Error)
	}
}

type FlagSettings interface {
	Type() string
	FromMap(*map[string][]string) error
	FromPlain([]byte) error
	Plain() []byte
	Check(string) bool
	IsValid() bool
}

type FlagStatic struct {
	Flag string `json:"flag"`
}

func (f *FlagStatic) Type() string {
	return FlagTypeStatic
}

func (f *FlagStatic) FromMap(m *map[string][]string) error {
	return mapstructure.Decode(m, f)
}

func (f *FlagStatic) FromPlain(b []byte) error {
	return json.Unmarshal(b, f)
}

func (f *FlagStatic) Plain() (j []byte) {
	j, _ = json.Marshal(f)
	return
}

func (f *FlagStatic) Check(s string) bool {
	return strings.Compare(f.Flag, s) == 0
}

func (f *FlagStatic) IsValid() bool {
	return len(f.Flag) > 0
}

func NewFlagSettingsFromMap(ft FlagType, m map[string]interface{}) (fs FlagSettings, err error) {
	switch ft {
	case FlagTypeStatic:
		var s FlagStatic
		err = mapstructure.Decode(m, &s)
		if err != nil {
			err = errors.ErrInvalidFlagSettings
			return
		}
		fs = &s
	default:
		err = errors.ErrInvalidFlagType
	}
	return
}
