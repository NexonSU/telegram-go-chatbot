package welcome

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

type BorderUser struct {
	User     *tb.User
	Status   string
	Reason   string
	JoinedAt time.Time
}

type JoinBorder struct {
	Message    *tb.Message
	Chat       *tb.Chat
	Users      []BorderUser
	NeedUpdate bool
	NeedCreate bool
}

var Border JoinBorder

var arabicSymbols, _ = regexp.Compile("[\u0600-\u06ff]|[\u0750-\u077f]|[\ufb50-\ufbc1]|[\ufbd3-\ufd3f]|[\ufd50-\ufd8f]|[\ufd92-\ufdc7]|[\ufe70-\ufefc]|[\uFDF0-\uFDFD]")

func OnJoin(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat {
		return
	}
	err := utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	if Border.Message == nil {
		Border.Message = m
	}
	log.Printf("New user detected in %v (%v)! ID: %v. Login: %v. Name: %v.", m.Chat.Title, m.Chat.ID, m.Sender.ID, utils.UserName(m.Sender), utils.UserFullName(m.Sender))
	Border.Chat = m.Chat
	Border.Users = append(Border.Users, BorderUser{
		User:     m.Sender,
		Status:   "pending",
		JoinedAt: time.Now(),
	})
	Border.NeedCreate = true
	ChatMember := &tb.ChatMember{
		Rights: tb.Rights{CanSendMessages: false},
		User:   m.Sender,
	}
	err = utils.Bot.Restrict(m.Chat, ChatMember)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	if arabicSymbols.MatchString(utils.UserFullName(m.Sender)) {
		err = utils.Bot.Ban(m.Chat, ChatMember)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		for i, e := range Border.Users {
			if e.User.ID == m.Sender.ID {
				Border.Users[i].Status = "banned"
				Border.Users[i].Reason = "арабская вязь в имени"
				Border.NeedUpdate = true
			}
		}
		return
	}
	if m.Sender.FirstName == "ICSM" {
		err = utils.Bot.Ban(m.Chat, ChatMember)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		for i, e := range Border.Users {
			if e.User.ID == m.Sender.ID {
				Border.Users[i].Status = "banned"
				Border.Users[i].Reason = "ICSM в имени"
				Border.NeedUpdate = true
			}
		}
		return
	}
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, _ := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", m.Sender.ID))
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(httpResponse.Body)
	jsonBytes, _ := ioutil.ReadAll(httpResponse.Body)
	if fastjson.GetBool(jsonBytes, "ok") {
		err = utils.Bot.Ban(m.Chat, ChatMember)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		for i, e := range Border.Users {
			if e.User.ID == m.Sender.ID {
				Border.Users[i].Status = "banned"
				Border.Users[i].Reason = "Combot Anti-Spam"
				Border.NeedUpdate = true
			}
		}
		return
	}
}
