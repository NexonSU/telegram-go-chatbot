package commands

import (
	tele "gopkg.in/telebot.v3"
)

// Reply with stats link
func Stats(context tele.Context) error {
	return context.Reply("<a href='https://t.me/zavtrachat_bot/stats'>Webapp</a>\n<a href='https://grafana.nexon.su/d/aef7a25c-3824-4046-8ed3-53ccb5850c9d/zavtrachat?kiosk'>Прямая ссылка</a>")
}
