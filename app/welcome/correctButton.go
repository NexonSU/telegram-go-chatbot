package welcome

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func OnClickCorrectButton(c *tb.Callback) {
	for i, e := range Border.Users {
		if e.User.ID == c.Sender.ID && e.Status == "pending" {
			var ChatMember tb.ChatMember
			ChatMember.User = c.Sender
			ChatMember.CanSendMessages = true
			ChatMember.CanSendMedia = true
			ChatMember.CanSendPolls = true
			ChatMember.CanSendOther = true
			ChatMember.CanAddPreviews = true
			ChatMember.RestrictedUntil = time.Now().Unix() + 60
			err := utils.Bot.Restrict(Border.Chat, &ChatMember)
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
			Border.Users[i].Status = "accepted"
			Border.NeedUpdate = true
			err = utils.Bot.Respond(c, &tb.CallbackResponse{Text: fmt.Sprintf("Добро пожаловать, %v!\nТеперь у тебя есть доступ к чату.", utils.UserFullName(c.Sender)), ShowAlert: true})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
		}
	}
	err := utils.Bot.Respond(c, &tb.CallbackResponse{})
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return
	}
}
