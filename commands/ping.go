package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Reply "Pong!" on /ping
func Ping(context tele.Context) error {
	return utils.ReplyAndRemove("Pong!", context)
}
