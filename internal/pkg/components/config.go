package components

import (
	"github.com/jinzhu/gorm"
)

const (
	TableConfig = "config"
)

type Config struct {
	gorm.Model
	Key   string `gorm:"size:100"`
	Value string `gorm:"type:longtext"`
}

func GetConfig(db *gorm.DB, key string) string {
	var config Config
	if v := db.Where(&Config{Key: key}).First(&config); v.Error == gorm.ErrRecordNotFound {
		return ""
	} else if v.Error != nil {
		panic(v.Error)
	}
	return config.Value
}

func SetConfig(db *gorm.DB, key string, value string) {
	var config Config
	if v := db.Where(&Config{Key: key}).FirstOrCreate(&config); v.Error != nil {
		panic(v.Error)
	}
	config.Value = value
	if v := db.Save(&config); v.Error != nil {
		panic(v.Error)
	}
}

func GetConfigs(db *gorm.DB) []Config {
	var configs []Config
	if v := db.Find(&configs); v.Error != nil {
		panic(v.Error)
	}
	return configs
}
