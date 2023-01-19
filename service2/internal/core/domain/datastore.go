package domain

import (
	"errors"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func MigrateDatabase(db *gorm.DB) {
	if db == nil {
		panic("An error when connect database")
	}

	db.AutoMigrate(&Stock{})

	tx := db.Begin()
	if db.Migrator().HasTable(&Stock{}) {
		if err := db.First(&Stock{}).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			logrus.Info("Inserting stock data")
			SeedingStock(tx)
		}
	}
	tx.Commit()
}

type Stock struct {
	gorm.Model
	Stock     *int    `gorm:"int;"`
	CreatedBy *string `gorm:"varchar(255);"`
	UpdatedBy *string `gorm:"varchar(255);"`
	DeletedBy *string `gorm:"varchar(255);"`
}

func (t Stock) TableName() string {
	return "stock"
}
