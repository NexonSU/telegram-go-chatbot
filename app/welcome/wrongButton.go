package welcome

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func OnClickWrongButton(context telebot.Context) error {
	for i, e := range Border.Users {
		if e.User.ID == context.Sender().ID && e.Status == "pending" {
			err := utils.Bot.Respond(context.Callback(), &telebot.CallbackResponse{Text: "Это неверный ответ, пока.", ShowAlert: true})
			if err != nil {
				return err
			}
			err = utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: context.Sender(), RestrictedUntil: time.Now().Unix() + 7200})
			if err != nil {
				return err
			}
			Border.Users[i].Status = "banned"
			Border.Users[i].Reason = "неверный ответ"
			Border.NeedUpdate = true
		}
	}
	return context.Respond(&telebot.CallbackResponse{})
}
