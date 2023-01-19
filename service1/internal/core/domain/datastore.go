package domain

import (
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	if db == nil {
		panic("An error when connect database")
	}

	db.AutoMigrate(&Machine{})
}

type Machine struct {
	gorm.Model
	Machine   *string `gorm:"varchar(255);"`
	Stock     *int    `gorm:"int;"`
	CreatedBy *string `gorm:"varchar(255);"`
	UpdatedBy *string `gorm:"varchar(255);"`
	DeletedBy *string `gorm:"varchar(255);"`
}

func (t Machine) TableName() string {
	return "machine"
}
