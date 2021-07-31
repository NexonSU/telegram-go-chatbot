package roulette

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func Deny(c *telebot.Callback) {
	err := utils.Bot.Respond(c, &telebot.CallbackResponse{})
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return err
	}
	victim := c.Message.Entities[0].User
	if victim.ID != c.Sender.ID {
		err := utils.Bot.Respond(c, &telebot.CallbackResponse{})
		if err != nil {
			utils.ErrorReporting(err, c.Message)
			return err
		}
		return err
	}
	busy["russianroulette"] = false
	busy["russianroulettePending"] = false
	busy["russianrouletteInProgress"] = false
	_, err = utils.Bot.Edit(c.Message, fmt.Sprintf("%v отказался от дуэли.", utils.UserFullName(c.Sender)))
	if err != nil {
		utils.ErrorReporting(err, c.Message)
		return err
	}
	return err
}
