package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send userid on /getid
func Getid(context telebot.Context) error {
	if context.Message().ReplyTo != nil && context.Message().ReplyTo.OriginalSender != nil {
		return context.Send(fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", context.Message().ReplyTo.OriginalSender.FirstName, context.Message().ReplyTo.OriginalSender.LastName, context.Message().ReplyTo.OriginalSender.Username, context.Message().ReplyTo.OriginalSender.ID))
	}
	if context.Message().ReplyTo != nil {
		return context.Send(fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", context.Message().ReplyTo.Sender.FirstName, context.Message().ReplyTo.Sender.LastName, context.Message().ReplyTo.Sender.Username, context.Message().ReplyTo.Sender.ID))
	}
	if len(context.Args()) == 1 {
		target, _, err := utils.FindUserInMessage(context)
		if err != nil {
			return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		}
		return context.Send(fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", target.FirstName, target.LastName, target.Username, target.ID))
	}
	return context.Send(fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", context.Sender().FirstName, context.Sender().LastName, context.Sender().Username, context.Sender().ID))
}
