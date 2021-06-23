package welcome

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func OnLeft(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat {
		return
	}
	err := utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
