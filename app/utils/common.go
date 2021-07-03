package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm/clause"
	"html"
	"log"
	"math/big"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func UserFullName(user *tb.User) string {
	fullname := user.FirstName
	if user.LastName != "" {
		fullname = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	}
	return fullname
}

func UserName(user *tb.User) string {
	username := user.Username
	if user.Username == "" {
		username = UserFullName(user)
	}
	return username
}

func MentionUser(user *tb.User) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a>", user.ID, UserFullName(user))
}

func RandInt(min int, max int) int {
	b, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return min + int(b.Int64())
}

func IsAdmin(username string) bool {
	for _, b := range Config.Telegram.Admins {
		if b == username {
			return true
		}
	}
	return false
}

func IsAdminOrModer(username string) bool {
	for _, b := range Config.Telegram.Admins {
		if b == username {
			return true
		}
	}
	for _, b := range Config.Telegram.Moders {
		if b == username {
			return true
		}
	}
	return false
}

func RestrictionTimeMessage(seconds int64) string {
	var message = ""
	if seconds-30 > time.Now().Unix() {
		message = fmt.Sprintf(" до %v", time.Unix(seconds, 0).Format("02.01.2006 15:04:05"))
	}
	return message
}

func ErrorReporting(err error, message *tb.Message) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[%s:%d] %v at MessageID \"%v\" in Chat \"%v\"", fn, line, err, message.ID, message.Chat.Username)
	MarshalledMessage, _ := json.MarshalIndent(message, "", "    ")
	JsonMessage := html.EscapeString(string(MarshalledMessage))
	text := fmt.Sprintf("An exception was raised while handling an update\n<pre>%v</pre>\n\nMessage:\n<pre>%v</pre>", err, JsonMessage)
	_, err = Bot.Send(tb.ChatID(Config.Telegram.SysAdmin), text)
	if err != nil {
		return
	}
}

func FindUserInMessage(m tb.Message) (tb.User, int64, error) {
	var user tb.User
	var err error = nil
	var untildate = time.Now().Unix()
	var text = strings.Split(m.Text, " ")
	if m.ReplyTo != nil {
		user = *m.ReplyTo.Sender
		if len(text) == 2 {
			addtime, err := strconv.ParseInt(text[1], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	} else {
		if len(text) == 1 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		user, err = GetUserFromDB(text[1])
		if err != nil {
			return user, untildate, err
		}
		if len(text) == 3 {
			addtime, err := strconv.ParseInt(text[2], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	}
	return user, untildate, err
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
