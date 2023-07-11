package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Reply "Polo!" on "marco"
func Marco(context tele.Context) error {
	return utils.ReplyAndRemove("Polo!", context)
}
