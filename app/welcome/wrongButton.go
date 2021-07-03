package welcome

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func OnClickWrongButton(c *tb.Callback) {
	for i, e := range Border.Users {
		if e.User.ID == c.Sender.ID && e.Status == "pending" {
			err := utils.Bot.Respond(c, &tb.CallbackResponse{Text: "Это неверный ответ, пока.", ShowAlert: true})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
			err = utils.Bot.Ban(Border.Chat, &tb.ChatMember{User: c.Sender, RestrictedUntil: time.Now().Unix() + 7200})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
			Border.Users[i].Status = "banned"
			Border.Users[i].Reason = "неверный ответ"
			Border.NeedUpdate = true
		}
	}
	err := utils.Bot.Respond(c, &tb.CallbackResponse{})
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return
	}
}
