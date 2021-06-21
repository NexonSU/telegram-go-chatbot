package services

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Repost channel post to chat
func OnPost(m *tb.Message) {
	if m.Chat.Username == utils.Config.Telegram.Channel {
		chat, err := utils.Bot.ChatByID("@" + utils.Config.Telegram.Chat)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		_, err = utils.Bot.Forward(chat, m)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
