package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Send text in chat on /say
func Say(context tele.Context) error {
	if len(context.Args()) == 0 {
		return utils.ReplyAndRemove("Укажите сообщение.", context)
	}
	context.Delete()
	return context.Send(utils.GetHtmlText(*context.Message()))
}
