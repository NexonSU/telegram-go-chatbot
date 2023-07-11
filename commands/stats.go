package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Reply with stats link
func Stats(context tele.Context) error {
	return utils.ReplyAndRemove("<a href='https://t.me/zavtrachat_bot/stats'>Webapp</a>\n<a href='https://grafana.nexon.su/d/aef7a25c-3824-4046-8ed3-53ccb5850c9d/zavtrachat?kiosk'>Прямая ссылка</a>", context)
}
