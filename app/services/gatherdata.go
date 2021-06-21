package services

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Gather user data on incoming text message
func OnText(m *tb.Message) {
	err := utils.GatherData(m.Sender)
	if err != nil {
		utils.ErrorReporting(err, m)
	}
}
