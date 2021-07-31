package welcome

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/valyala/fastjson"
	"gopkg.in/tucnak/telebot.v3"
)

type BorderUser struct {
	User     *telebot.User
	Status   string
	Reason   string
	JoinedAt time.Time
}

type JoinBorder struct {
	Message    *telebot.Message
	Chat       *telebot.Chat
	Users      []BorderUser
	NeedUpdate bool
	NeedCreate bool
}

var Border JoinBorder

var arabicSymbols, _ = regexp.Compile("[\u0600-\u06ff]|[\u0750-\u077f]|[\ufb50-\ufbc1]|[\ufbd3-\ufd3f]|[\ufd50-\ufd8f]|[\ufd92-\ufdc7]|[\ufe70-\ufefc]|[\uFDF0-\uFDFD]")

func OnJoin(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat {
		return err
	}
	err = utils.Bot.Delete(context.Message())
	if err != nil {
		return err
	}
	if Border.Message == nil {
		Border.Message = context.Message()
	}
	log.Printf("New user detected in %v (%v)! ID: %v. Login: %v. Name: %v.", context.Chat().Title, context.Chat().ID, context.Sender().ID, utils.UserName(context.Sender()), utils.UserFullName(context.Sender()))
	Border.Chat = context.Chat()
	Border.Users = append(Border.Users, BorderUser{
		User:     context.Sender(),
		Status:   "pending",
		JoinedAt: time.Now(),
	})
	Border.NeedCreate = true
	ChatMember := &telebot.ChatMember{
		Rights: telebot.Rights{CanSendMessages: false},
		User:   context.Sender(),
	}
	err = utils.Bot.Restrict(context.Chat(), ChatMember)
	if err != nil {
		return err
	}
	if arabicSymbols.MatchString(utils.UserFullName(context.Sender())) {
		err = utils.Bot.Ban(context.Chat(), ChatMember)
		if err != nil {
			return err
		}
		for i, e := range Border.Users {
			if e.User.ID == context.Sender().ID {
				Border.Users[i].Status = "banned"
				Border.Users[i].Reason = "арабская вязь в имени"
				Border.NeedUpdate = true
			}
		}
		return err
	}
	if context.Sender().FirstName == "ICSM" {
		err = utils.Bot.Ban(context.Chat(), ChatMember)
		if err != nil {
			return err
		}
		for i, e := range Border.Users {
			if e.User.ID == context.Sender().ID {
				Border.Users[i].Status = "banned"
				Border.Users[i].Reason = "ICSM в имени"
				Border.NeedUpdate = true
			}
		}
		return err
	}
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, _ := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", context.Sender().ID))
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(httpResponse.Body)
	jsonBytes, _ := ioutil.ReadAll(httpResponse.Body)
	if fastjson.GetBool(jsonBytes, "ok") {
		err = utils.Bot.Ban(context.Chat(), ChatMember)
		if err != nil {
			return err
		}
		for i, e := range Border.Users {
			if e.User.ID == context.Sender().ID {
				Border.Users[i].Status = "banned"
				Border.Users[i].Reason = "Combot Anti-Spam"
				Border.NeedUpdate = true
			}
		}
		return err
	}
	return err
}
