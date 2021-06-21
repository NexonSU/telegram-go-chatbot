package roulette

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

func Deny(c *tb.Callback) {
	err := utils.Bot.Respond(c, &tb.CallbackResponse{})
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return
	}
	victim := c.Message.Entities[0].User
	if victim.ID != c.Sender.ID {
		err := utils.Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			utils.ErrorReporting(err, c.Message)
			return
		}
		return
	}
	busy["russianroulette"] = false
	busy["russianroulettePending"] = false
	busy["russianrouletteInProgress"] = false
	_, err = utils.Bot.Edit(c.Message, fmt.Sprintf("%v отказался от дуэли.", utils.UserFullName(c.Sender)))
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return
	}
}
