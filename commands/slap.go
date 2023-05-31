package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Send slap message on /slap
func Slap(context tele.Context) error {
	var action = "дал леща"
	var target tele.User
	if utils.IsAdminOrModer(context.Sender().ID) {
		action = "дал отцовского леща"
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	return context.Send(fmt.Sprintf("👋 <b>%v</b> %v %v", utils.UserFullName(context.Sender()), action, utils.MentionUser(&target)))
}
