package commands

import (
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send text in chat on /say
func Say(context telebot.Context) error {
	var text = strings.Split(context.Text(), " ")
	if len(text) > 1 {
		utils.Bot.Delete(context.Message())
		return context.Send(strings.Join(text[1:], " "))
	} else {
		return context.Reply("Укажите сообщение.")
	}
}
