package welcome

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func OnClickCorrectButton(context telebot.Context) error {
	for i, e := range Border.Users {
		if e.User.ID == context.Sender().ID && e.Status == "pending" {
			var ChatMember telebot.ChatMember
			ChatMember.User = context.Sender()
			ChatMember.CanSendMessages = true
			ChatMember.CanSendMedia = true
			ChatMember.CanSendPolls = true
			ChatMember.CanSendOther = true
			ChatMember.CanAddPreviews = true
			ChatMember.RestrictedUntil = time.Now().Unix() + 60
			err := utils.Bot.Restrict(Border.Chat, &ChatMember)
			if err != nil {
				return err
			}
			Border.Users[i].Status = "accepted"
			Border.NeedUpdate = true
			err = context.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Добро пожаловать, %v!\nТеперь у тебя есть доступ к чату.", utils.UserFullName(context.Sender())), ShowAlert: true})
			if err != nil {
				return err
			}
		}
	}
	return context.Respond(&telebot.CallbackResponse{})
}
