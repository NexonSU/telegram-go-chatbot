package utils

import (
	"log"
	"time"

	"github.com/NexonSU/telebot"
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

type AntiSpam struct {
	Text string `gorm:"primaryKey"`
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

type Nope struct {
	Text string `gorm:"primaryKey"`
}

type Message struct {
	ID       int   `gorm:"primaryKey"`
	ChatID   int64 `gorm:"primaryKey"`
	UserID   int64
	Date     time.Time
	ReplyTo  int
	Text     string
	FileType string
	FileID   string
}

type Word struct {
	ChatID int64
	UserID int64
	Date   time.Time
	Text   string
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
	err = database.AutoMigrate(telebot.User{}, Message{}, Word{}, Get{}, Warn{}, PidorStats{}, PidorList{}, Duelist{}, Nope{}, AntiSpam{})
	if err != nil {
		log.Println(err)
	}
	return *database
}

var DB = DataBaseInit("bot.db")
