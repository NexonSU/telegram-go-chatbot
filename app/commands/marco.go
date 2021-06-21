package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Reply "Polo!" on "marco"
func Marco(m *tb.Message) {
	_, err := utils.Bot.Reply(m, "Polo!")
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
