package commands

import (
	"os"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

//Restart bot on /restart
func Restart(context telebot.Context) error {
	utils.Bot.Delete(context.Message())
	os.Exit(0)
	return nil
}
