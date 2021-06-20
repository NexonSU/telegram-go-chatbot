package utils

import "time"

type Warn struct {
	UserID   int `gorm:"primaryKey"`
	Amount   int
	LastWarn time.Time
}
