package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
	"os"
)

var database *gorm.DB

func GetDatabase() (*gorm.DB, error) {
	if database != nil {
		if err := database.Debug().DB().Ping(); err == nil {
			return database, nil
		}
	}

	tp := os.Getenv("DISCORD_TASK_MANAGEMENT_DATABASE_TYPE")
	str := os.Getenv("DISCORD_TASK_MANAGEMENT_DATABASE_CONNECTION_STR")
	db, err := gorm.Open(tp, str)
	if err != nil {
		return nil, err
	}

	database = db

	return database, nil
}

func Migration() {
	db, err := GetDatabase()
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	db.AutoMigrate(&Creator{}, &Client{}, &Request{})
}
