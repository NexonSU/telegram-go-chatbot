package duel

import (
	"fmt"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func Deny(context telebot.Context) error {
	victim := context.Message().Entities[0].User
	if victim.ID != context.Sender().ID {
		return context.Respond(&telebot.CallbackResponse{Text: utils.GetNope()})
	}
	err := utils.Bot.Respond(context.Callback(), &telebot.CallbackResponse{})
	if err != nil {
		return err
	}
	busy["russianroulette"] = false
	busy["russianroulettePending"] = false
	busy["russianrouletteInProgress"] = false
	return context.Edit(fmt.Sprintf("%v отказался от дуэли.", context.Sender().FullName()))
}
