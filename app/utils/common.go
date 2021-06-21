package utils

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"html"
	"log"
	"math/big"
	"runtime"
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
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
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
	chat, err := Bot.ChatByID("@" + Config.Telegram.SysAdmin)
	if err != nil {
		return
	}
	_, err = Bot.Send(chat, text)
	if err != nil {
		return
	}
}
