package utils

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"strconv"
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

var DB, _ = gorm.Open(sqlite.Open("../../bot.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})

func GetUserFromDB(findstring string) (tb.User, error) {
	var user tb.User
	var err error = nil
	if string(findstring[0]) == "@" {
		user.Username = findstring[1:]
	} else {
		user.ID, err = strconv.Atoi(findstring)
	}
	result := DB.Where(&user).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return user, err
}

func GatherData(user *tb.User) error {
	result := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
