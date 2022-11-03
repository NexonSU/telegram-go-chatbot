package commands

import (
	"os"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

//Restart bot on /restart
func Restart(context tele.Context) error {
	utils.Bot.Delete(context.Message())
	os.Exit(0)
	return nil
}
