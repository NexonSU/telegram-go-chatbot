package checkpoint

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func OnLeft(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat {
		return err
	}
	err = utils.Bot.Delete(context.Message())
	for i, user := range Border.Users {
		if user.User.ID == context.Sender().ID && user.Status == "pending" {
			err := utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: user.User, RestrictedUntil: time.Now().Unix() + 7200})
			if err != nil {
				continue
			}
			Border.Users[i].Status = "banned"
			Border.Users[i].Reason = "сбежал"
			Border.NeedUpdate = true
		}
	}
	if err != nil {
		return err
	}
	return err
}
