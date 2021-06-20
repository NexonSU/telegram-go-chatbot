package commands

import (
	"encoding/json"
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Return message on /debug command
func Debug(m *tb.Message) {
	if !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Admins) && !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Moders) {
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var message = m
	if m.ReplyTo != nil {
		message = m.ReplyTo
	}
	MarshalledMessage, _ := json.MarshalIndent(message, "", "    ")
	_, err := utils.Bot.Reply(m, fmt.Sprintf("<pre>%v</pre>", string(MarshalledMessage)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
