package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send formatted text on /me
func Me(context telebot.Context) error {
	var text = strings.Split(context.Text(), " ")
	if len(text) == 1 {
		return context.Reply("Пример использования:\n<code>/me {делает что-то}</code>")
	}
	utils.Bot.Delete(context.Message())
	return context.Send(fmt.Sprintf("<code>%v %v</code>", utils.UserFullName(context.Sender()), strings.Join(text[1:], " ")))
}
