package utils

import (
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Get struct {
	Name    string `gorm:"primaryKey"`
	Type    string
	Data    string
	Caption string
	Creator int64
}

type AntiSpamLink struct {
	URL  string `gorm:"primaryKey"`
	Type string
}

type PidorStats struct {
	Date   time.Time `gorm:"primaryKey"`
	UserID int64
}

type PidorList telebot.User

type Duelist struct {
	UserID int64 `gorm:"primaryKey"`
	Deaths int
	Kills  int
}

type Warn struct {
	UserID   int64 `gorm:"primaryKey"`
	Amount   int
	LastWarn time.Time
}

type ZavtraStream struct {
	Service   string `gorm:"primaryKey"`
	LastCheck time.Time
	VideoID   string
}

type Nope struct {
	Text string `gorm:"primaryKey"`
}

func DataBaseInit(file string) gorm.DB {
	database, err := gorm.Open(
		sqlite.Open(file),
		&gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	//Create tables, if they not exists in DB
	err = database.AutoMigrate(telebot.User{}, Get{}, Warn{}, PidorStats{}, PidorList{}, Duelist{}, ZavtraStream{}, Nope{}, AntiSpamLink{})
	if err != nil {
		log.Println(err)
	}
	return *database
}

var DB = DataBaseInit("bot.db")
