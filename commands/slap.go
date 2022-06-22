package commands

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Send slap message on /slap
func Slap(context tele.Context) error {
	var action = "дал леща"
	var target tele.User
	if utils.IsAdminOrModer(context.Sender().ID) {
		action = "дал отцовского леща"
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Send(fmt.Sprintf("👋 <b>%v</b> %v %v", utils.UserFullName(context.Sender()), action, utils.MentionUser(&target)))
}
