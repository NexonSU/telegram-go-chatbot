package utils

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Get struct {
	Name    string `gorm:"primaryKey"`
	Type    string
	Data    string
	Caption string
}

type PidorStats struct {
	Date   time.Time `gorm:"primaryKey"`
	UserID int
}

type PidorList tb.User

type Duelist struct {
	UserID int `gorm:"primaryKey"`
	Deaths int
	Kills  int
}

type Warn struct {
	UserID   int `gorm:"primaryKey"`
	Amount   int
	LastWarn time.Time
}

type ZavtraStream struct {
	Service   string `gorm:"primaryKey"`
	LastCheck time.Time
	VideoID   string
}

func DataBaseInit(file string) gorm.DB {
	database, err := gorm.Open(
		sqlite.Open(file),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Info),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	//Create tables, if they not exists in DB
	err = database.AutoMigrate(tb.User{}, Get{}, Warn{}, PidorStats{}, PidorList{}, Duelist{}, ZavtraStream{})
	if err != nil {
		log.Println(err)
	}
	return *database
}

var DB = DataBaseInit("bot.db")
