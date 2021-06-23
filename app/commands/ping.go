package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Reply "Pong!" on "ping"
func Ping(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	_, err := utils.Bot.Reply(m, "Pong!")
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
