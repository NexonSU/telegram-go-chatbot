package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send userid on /getid
func Getid(m *tb.Message) {
	if !utils.IsAdminOrModer(m.Sender.Username) {
		if m.Chat.Username != utils.Config.Telegram.Chat {
			return
		}
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var text = strings.Split(m.Text, " ")
	if m.ReplyTo != nil && m.ReplyTo.OriginalSender != nil {
		_, err := utils.Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", m.ReplyTo.OriginalSender.FirstName, m.ReplyTo.OriginalSender.LastName, m.ReplyTo.OriginalSender.Username, m.ReplyTo.OriginalSender.ID))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else if m.ReplyTo != nil {
		_, err := utils.Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName, m.ReplyTo.Sender.Username, m.ReplyTo.Sender.ID))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else if len(text) == 2 {
		target, _, err := utils.FindUserInMessage(*m)
		if err != nil {
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = utils.Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", target.FirstName, target.LastName, target.Username, target.ID))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else {
		_, err := utils.Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", m.Sender.FirstName, m.Sender.LastName, m.Sender.Username, m.Sender.ID))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
