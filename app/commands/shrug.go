package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send shrug in chat on /shrug
func Shrug(m *tb.Message) {
	_, err := utils.Bot.Send(m.Chat, "¯\\_(ツ)_/¯")
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
