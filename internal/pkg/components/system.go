package components

import "github.com/jinzhu/gorm"

type System struct {
	gorm.Model
	Key   string `json:"key"`
	Value string `json:"value"`
}

func MigrateSystem(db *gorm.DB) {
	db.AutoMigrate(&System{})
}

func SetItem() {

}
