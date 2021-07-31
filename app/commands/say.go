package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send text in chat on /say
func Say(context telebot.Context) error {
	if len(context.Args()) == 0 {
		return context.Reply("Укажите сообщение.")
	}
	utils.Bot.Delete(context.Message())
	return context.Send(context.Data())
}
