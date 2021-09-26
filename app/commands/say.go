package commands

import (
	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

//Send text in chat on /say
func Say(context telebot.Context) error {
	if len(context.Args()) == 0 {
		return context.Reply("Укажите сообщение.")
	}
	context.Delete()
	return context.Send(utils.GetHtmlText(*context.Message()))
}
