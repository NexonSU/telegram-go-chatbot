package commands

import (
	"fmt"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

//Send formatted text on /me
func Me(context telebot.Context) error {
	if len(context.Args()) == 0 {
		return context.Reply("Пример использования:\n<code>/me {делает что-то}</code>")
	}
	utils.Bot.Delete(context.Message())
	return context.Send(fmt.Sprintf("<code>%v %v</code>", utils.UserFullName(context.Sender()), context.Data()))
}
