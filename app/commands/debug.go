package commands

import (
	"encoding/json"
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Return message on /debug command
func Debug(m *tb.Message) {
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
	err := utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	var message = m
	if m.ReplyTo != nil {
		message = m.ReplyTo
	}
	MarshalledMessage, _ := json.MarshalIndent(message, "", "    ")
	_, err = utils.Bot.Send(m.Sender, fmt.Sprintf("<pre>%v</pre>", string(MarshalledMessage)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
