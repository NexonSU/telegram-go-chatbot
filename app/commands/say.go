package commands

import (
	"gopkg.in/tucnak/telebot.v3"
)

//Send text in chat on /say
func Say(context telebot.Context) error {
	if len(context.Args()) == 0 {
		return context.Reply("Укажите сообщение.")
	}
	context.Delete()
	return context.Send(context.Data())
}
