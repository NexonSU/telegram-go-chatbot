package commands

import (
	tele "gopkg.in/telebot.v3"
)

// Reply with stats link
func Stats(context tele.Context) error {
	return context.Reply("http://t.me/zavtrachat_bot/stats")
}
