package welcome

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func OnLeft(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat {
		return
	}
	err := utils.Bot.Delete(m)
	for i, user := range Border.Users {
		if user.User.ID == m.Sender.ID && user.Status == "pending" {
			err := utils.Bot.Ban(Border.Chat, &tb.ChatMember{User: user.User, RestrictedUntil: time.Now().Unix() + 7200})
			if err != nil {
				continue
			}
			Border.Users[i].Status = "banned"
			Border.Users[i].Reason = "сбежал"
			Border.NeedUpdate = true
		}
	}
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
