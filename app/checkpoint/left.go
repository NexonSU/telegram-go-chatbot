package checkpoint

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func UserLeft(context telebot.Context) error {
	for i, user := range Border.Users {
		if user.User.ID == context.ChatMember().NewChatMember.User.ID && user.Status == "pending" {
			err := utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: user.User, RestrictedUntil: time.Now().Unix() + 7200})
			if err != nil {
				return err
			}
			Border.Users[i].Status = "banned"
			Border.Users[i].Reason = "сбежал"
			Border.NeedUpdate = true
		}
	}
	return nil
}
