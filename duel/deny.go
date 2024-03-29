package duel

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

func Deny(context tele.Context) error {
	victim := context.Message().Entities[0].User
	if victim.ID != context.Sender().ID {
		return context.Respond(&tele.CallbackResponse{Text: utils.GetNope()})
	}
	err := utils.Bot.Respond(context.Callback(), &tele.CallbackResponse{})
	if err != nil {
		return err
	}
	busy["russianroulette"] = false
	busy["russianroulettePending"] = false
	busy["russianrouletteInProgress"] = false
	return context.Edit(fmt.Sprintf("%v отказался от дуэли.", utils.UserFullName(context.Sender())))
}
