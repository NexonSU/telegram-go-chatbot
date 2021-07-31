package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"math/big"
	"runtime"
	"strconv"
	"time"

	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

func UserFullName(user *telebot.User) string {
	fullname := user.FirstName
	if user.LastName != "" {
		fullname = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	}
	return fullname
}

func UserName(user *telebot.User) string {
	username := user.Username
	if user.Username == "" {
		username = UserFullName(user)
	}
	return username
}

func MentionUser(user *telebot.User) string {
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

func ErrorReporting(err error, context telebot.Context) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[%s:%d] %v at MessageID \"%v\" in Chat \"%v\"", fn, line, err, context.Message().ID, context.Chat().Username)
	MarshalledMessage, _ := json.MarshalIndent(context.Message(), "", "    ")
	JsonMessage := html.EscapeString(string(MarshalledMessage))
	text := fmt.Sprintf("An exception was raised while handling an update\n<pre>%v</pre>\n\nMessage:\n<pre>%v</pre>", err, JsonMessage)
	Bot.Send(telebot.ChatID(Config.Telegram.SysAdmin), text)
}

func FindUserInMessage(context telebot.Context) (telebot.User, int64, error) {
	var user telebot.User
	var err error = nil
	var untildate = time.Now().Unix()
	if context.Message().ReplyTo != nil {
		user = *context.Message().ReplyTo.Sender
		if len(context.Args()) == 1 {
			addtime, err := strconv.ParseInt(context.Args()[0], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	} else {
		if len(context.Args()) == 0 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		user, err = GetUserFromDB(context.Args()[0])
		if err != nil {
			return user, untildate, err
		}
		if len(context.Args()) == 2 {
			addtime, err := strconv.ParseInt(context.Args()[1], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	}
	return user, untildate, err
}

func GatherData(user *telebot.User) error {
	result := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserFromDB(findstring string) (telebot.User, error) {
	var user telebot.User
	var err error = nil
	if string(findstring[0]) == "@" {
		user.Username = findstring[1:]
	} else {
		user.ID, err = strconv.ParseInt(findstring, 10, 64)
	}
	result := DB.Where(&user).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return user, err
}

//Repost channel post to chat
func Repost(context telebot.Context) error {
	chat, err := Bot.ChatByID("@" + Config.Telegram.Chat)
	if err != nil {
		return err
	}
	return context.ForwardTo(chat)
}

//Remove message
func Remove(context telebot.Context) error {
	return context.Delete()
}
