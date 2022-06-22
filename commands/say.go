package commands

import (
	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Send text in chat on /say
func Say(context tele.Context) error {
	if len(context.Args()) == 0 {
		return context.Reply("Укажите сообщение.")
	}
	context.Delete()
	return context.Send(utils.GetHtmlText(*context.Message()))
}
