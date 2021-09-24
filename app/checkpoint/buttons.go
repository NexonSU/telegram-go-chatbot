package checkpoint

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func ButtonCallback(context telebot.Context) error {
	if CorrectButton.Data == context.Data() {
		for i, e := range Border.Users {
			if e.User.ID == context.Sender().ID && e.Status == "pending" {
				if e.Role == "member" {
					var ChatMember telebot.ChatMember
					ChatMember.User = context.Sender()
					ChatMember.CanSendMessages = true
					ChatMember.RestrictedUntil = time.Now().Unix() + 40
					err := utils.Bot.Restrict(Border.Chat, &ChatMember)
					if err != nil {
						return err
					}
				}
				Border.Users[i].Status = "accepted"
				Border.NeedUpdate = true
				return context.Respond(&telebot.CallbackResponse{Text: fmt.Sprintf("Добро пожаловать, %v!\nТеперь у тебя есть доступ к чату.", utils.UserFullName(context.Sender())), ShowAlert: true})
			}
		}
		return context.Respond(&telebot.CallbackResponse{Text: utils.GetNope()})
	}
	if FirstWrongButton.Data == context.Data() || SecondWrongButton.Data == context.Data() || ThirdWrongButton.Data == context.Data() {
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
		return context.Respond(&telebot.CallbackResponse{Text: utils.GetNope()})
	}
	time.Sleep(2000 * time.Microsecond)
	return context.Respond(&telebot.CallbackResponse{Text: utils.GetNope()})
}
