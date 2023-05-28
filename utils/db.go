package utils

import (
	"log"
	"time"

	tele "gopkg.in/telebot.v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Get struct {
	Name    string `gorm:"primaryKey"`
	Title   string
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

type PidorList tele.User

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

type Bless struct {
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
	err = database.AutoMigrate(tele.User{}, Get{}, Warn{}, PidorStats{}, PidorList{}, Duelist{}, Bless{}, Nope{})
	if err != nil {
		log.Println(err)
	}
	database.Exec("DELETE FROM anti_spams")
	database.Exec("DELETE FROM check_point_restricts")
	database.Exec("DELETE FROM messages")
	database.Exec("DELETE FROM word_stats_excludes")
	database.Exec("DELETE FROM words")
	return *database
}

var DB = DataBaseInit("bot.db")
