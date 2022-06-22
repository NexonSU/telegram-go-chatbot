package checkpoint

import (
	"time"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

func UserLeft(context tele.Context) error {
	for _, user := range utils.RestrictedUsers {
		if user.UserID != context.ChatMember().NewChatMember.User.ID {
			continue
		}
		delete := utils.DB.Delete(&user)
		if delete.Error != nil {
			return delete.Error
		}
		err := utils.Bot.Ban(&tele.Chat{ID: utils.Config.Chat}, &tele.ChatMember{User: &tele.User{ID: user.UserID}, RestrictedUntil: time.Now().Unix() + 3600})
		if err != nil {
			return err
		}
	}
	return nil
}
