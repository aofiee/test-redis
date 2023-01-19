package domain

import "gorm.io/gorm"

func SeedingStock(tx *gorm.DB) {
	if tx == nil {
		panic("An error when connect database")
	}
	tx.Exec(`INSERT INTO "stock" ("stock") VALUES
	(100000);`)
	if tx.Error != nil {
		tx.Rollback()
		panic("An error when seeding geography data")
	}
}
