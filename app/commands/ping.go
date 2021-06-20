package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Reply "Pong!" on "ping"
func Ping(m *tb.Message) {
	_, err := utils.Bot.Reply(m, "Pong!")
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
