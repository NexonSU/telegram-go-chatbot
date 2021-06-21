package userActions

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

var nopes = []string{"неа", "не", "нет", "не то", "не попал"}

func OnClickWrongButton(c *tb.Callback) {
	for _, element := range c.Message.Entities {
		if element.User.ID == c.Sender.ID {
			err := utils.Bot.Respond(c, &tb.CallbackResponse{Text: nopes[utils.RandInt(0, len(nopes)-1)]})
			if err != nil {
				utils.ErrorReporting(err, c.Message)
				return
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
