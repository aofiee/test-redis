package gorm

import (
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DB struct {
	Postgres *gorm.DB
}

var dbConnect = &DB{}

func Connect2Postgres(host, port, username, password, dbName string, sslMode bool) (*DB, error) {
	var (
		err           error
		connectionStr string
	)
	if host == "" && port == "" && dbName == "" {
		return nil, errors.New("cannot estabished the connection")
	}
	if sslMode {
		connectionStr = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=require", host, username, password, dbName, port)
	} else {
		connectionStr = fmt.Sprintf("host=%v user=%v password=%v dbname=%v port=%v sslmode=disable", host, username, password, dbName, port)
	}
	dial := postgres.Open(connectionStr)
	pg, err := gorm.Open(dial, &gorm.Config{
		DryRun: false,
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic(err)
	}
	dbConnect.Postgres = pg
	return dbConnect, nil
}

func DisconnectPostgres(db *gorm.DB) {
	sqlDb, err := db.DB()
	if err != nil {
		panic("close db")
	}
	sqlDb.Close()
	log.Println("Connected with postgres has closed")
}
