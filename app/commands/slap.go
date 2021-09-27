package commands

import (
	"fmt"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

//Send slap message on /slap
func Slap(context telebot.Context) error {
	var action = "дал леща"
	var target telebot.User
	if utils.IsAdminOrModer(context.Sender().Username) {
		action = "дал отцовского леща"
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Send(fmt.Sprintf("👋 <b>%v</b> %v %v", context.Sender().FullName(), action, target.MentionHTML()))
}
