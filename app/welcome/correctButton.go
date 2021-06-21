package welcome

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func OnClickCorrectButton(c *tb.Callback) {
	for _, element := range c.Message.Entities {
		if element.User.ID == c.Sender.ID {
			err := utils.Bot.Respond(c, &tb.CallbackResponse{Text: fmt.Sprintf("Добро пожаловать, %v!\nТеперь у тебя есть доступ к чату.", utils.UserFullName(c.Sender)), ShowAlert: true})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
			ChatMember, err := utils.Bot.ChatMemberOf(c.Message.Chat, c.Sender)
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
			ChatMember.CanSendMessages = true
			ChatMember.RestrictedUntil = time.Now().Add(time.Hour).Unix()
			err = utils.Bot.Promote(c.Message.Chat, ChatMember)
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
			}
			if len(c.Message.Entities) == 1 {
				if Message.ID == c.Message.ID {
					Message.Unixtime = 0
				}
				err = utils.Bot.Delete(c.Message)
				if err != nil {
					utils.ErrorReporting(err, c.Message)
					return
				}
			} else {
				text := "Добро пожаловать"
				for _, element := range c.Message.Entities {
					if element.User.ID != c.Sender.ID {
						text += ", " + utils.MentionUser(element.User)
					}
				}
				text += "!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста?"
				_, err = utils.Bot.Edit(c.Message, text, &Selector)
				if err != nil {
					utils.ErrorReporting(err, c.Message)
					return
				}
			}
			return
		}
	}
	err := utils.Bot.Respond(c, &tb.CallbackResponse{})
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return
	}
}
