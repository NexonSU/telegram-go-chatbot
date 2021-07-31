package welcome

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func OnClickWrongButton(c *telebot.Callback) {
	for i, e := range Border.Users {
		if e.User.ID == c.Sender.ID && e.Status == "pending" {
			err := utils.Bot.Respond(c, &telebot.CallbackResponse{Text: "Это неверный ответ, пока.", ShowAlert: true})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return err
			}
			err = utils.Bot.Ban(Border.Chat, &telebot.ChatMember{User: c.Sender, RestrictedUntil: time.Now().Unix() + 7200})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return err
			}
			Border.Users[i].Status = "banned"
			Border.Users[i].Reason = "неверный ответ"
			Border.NeedUpdate = true
		}
	}
	err := utils.Bot.Respond(c, &telebot.CallbackResponse{})
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return err
	}
	return err
}
