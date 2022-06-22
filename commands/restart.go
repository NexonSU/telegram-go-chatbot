package commands

import (
	"os"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Restart bot on /restart
func Restart(context tele.Context) error {
	utils.Bot.Delete(context.Message())
	os.Exit(0)
	return nil
}
