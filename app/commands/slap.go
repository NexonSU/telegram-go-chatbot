package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
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
		err := context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			return err
		}
		return err
	}
	_, err = utils.Bot.Send(context.Chat(), fmt.Sprintf("👋 <b>%v</b> %v %v", utils.UserFullName(context.Sender()), action, utils.MentionUser(&target)))
	if err != nil {
		return err
	}
	return err
}
