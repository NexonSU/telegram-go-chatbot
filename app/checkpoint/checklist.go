package checkpoint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/valyala/fastjson"
)

func Check(User BorderUser) BorderUser {
	if User.Status != "pending" {
		return User
	}
	if time.Now().Unix()-User.JoinedAt > 120 {
		err := utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: User.User})
		if err != nil {
			return User
		}
		User.Status = "banned"
		User.Reason = "не ответил на вопрос"
		Border.NeedUpdate = true
		return User
	}
	if User.Checked {
		return User
	}
	if arabicSymbols.MatchString(User.User.FullName()) {
		err := utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: User.User})
		if err != nil {
			return User
		}
		User.Status = "banned"
		User.Reason = "арабская вязь в имени"
		Border.NeedUpdate = true
		return User
	}
	if User.User.FirstName == "ICSM" {
		err := utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: User.User})
		if err != nil {
			return User
		}
		User.Status = "banned"
		User.Reason = "ICSM в имени"
		Border.NeedUpdate = true
		return User
	}
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", User.User.ID))
	if err != nil {
		return User
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(httpResponse.Body)
	jsonBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return User
	}
	if fastjson.GetBool(jsonBytes, "ok") {
		err := utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: User.User})
		if err != nil {
			return User
		}
		User.Status = "banned"
		User.Reason = "Combot Anti-Spam"
		Border.NeedUpdate = true
		return User
	}
	return User
}
