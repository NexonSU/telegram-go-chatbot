package roulette

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func Deny(context telebot.Context) error {
	err := utils.Bot.Respond(context.Callback(), &telebot.CallbackResponse{})
	if err != nil {
		return err
	}
	victim := context.Message().Entities[0].User
	if victim.ID != context.Sender().ID {
		return context.Respond(&telebot.CallbackResponse{})
	}
	busy["russianroulette"] = false
	busy["russianroulettePending"] = false
	busy["russianrouletteInProgress"] = false
	return context.Edit(fmt.Sprintf("%v отказался от дуэли.", utils.UserFullName(context.Sender())))
}
