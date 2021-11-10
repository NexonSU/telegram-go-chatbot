package checkpoint

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func UserJoin(context telebot.Context) error {
	ChatMember := context.ChatMember().NewChatMember
	Border.Users = append(Border.Users, BorderUser{
		User:     ChatMember.User,
		Status:   "pending",
		Role:     string(ChatMember.Role),
		JoinedAt: time.Now().Unix(),
	})
	Border.NeedCreate = true
	if ChatMember.Role == "member" {
		ChatMember.CanSendMessages = false
		ChatMember.RestrictedUntil = time.Now().Unix() + 300
		return utils.Bot.Restrict(Border.Chat, ChatMember)
	}
	return nil
}
