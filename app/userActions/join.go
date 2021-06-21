package userActions

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/valyala/fastjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"
)

func OnJoin(welcomeMessage *tb.Message, welcomeSelector tb.ReplyMarkup, arab *regexp.Regexp) func(*tb.Message) {
	return func(m *tb.Message) {
		if welcomeMessage == nil {
			welcomeMessage = m
			welcomeMessage.Unixtime = 0
		}
		if m.Chat.Username != utils.Config.Telegram.Chat {
			return
		}
		err := utils.Bot.Delete(m)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		log.Printf("New user detected in %v (%v)! ID: %v. Login: %v. Name: %v.", m.Chat.Title, m.Chat.ID, m.Sender.ID, utils.UserName(m.Sender), utils.UserFullName(m.Sender))
		User := m.Sender
		Chat := m.Chat
		ChatMember, err := utils.Bot.ChatMemberOf(Chat, User)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		ChatMember.CanSendMessages = false
		err = utils.Bot.Restrict(Chat, ChatMember)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		if arab.MatchString(utils.UserFullName(User)) || User.FirstName == "ICSM" {
			err = utils.Bot.Ban(Chat, ChatMember)
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		}
		var httpClient = &http.Client{Timeout: 10 * time.Second}
		httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", User.ID))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(httpResponse.Body)
		jsonBytes, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		if fastjson.GetBool(jsonBytes, "ok") {
			err = utils.Bot.Ban(Chat, ChatMember)
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		}
		if time.Now().Unix()-welcomeMessage.Time().Unix() > 10 {
			welcomeMessage, err = utils.Bot.Send(Chat, fmt.Sprintf("Добро пожаловать, %v!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста?", utils.MentionUser(User)), &welcomeSelector)
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		} else {
			text := "Добро пожаловать"
			for _, element := range welcomeMessage.Entities {
				text += ", " + utils.MentionUser(element.User)
			}
			text += ", " + utils.MentionUser(m.Sender) + "!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста? У тебя 2 минуты."
			_, err = utils.Bot.Edit(welcomeMessage, text, &welcomeSelector)
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		}
		go func() {
			message := welcomeMessage
			time.Sleep(time.Second * 120)
			ChatMember, err := utils.Bot.ChatMemberOf(m.Chat, m.Sender)
			if err != nil {
				return
			}
			if ChatMember.Role != "member" {
				err := utils.Bot.Ban(m.Chat, &tb.ChatMember{User: m.Sender})
				if err != nil {
					utils.ErrorReporting(err, m)
					return
				}
			}
			err = utils.Bot.Delete(message)
			if err != nil {
				return
			}
		}()
	}
}
