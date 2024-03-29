package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Send formatted text on /me
func Me(context tele.Context) error {
	if len(context.Args()) == 0 {
		return utils.ReplyAndRemove("Пример использования:\n<code>/me {делает что-то}</code>", context)
	}
	utils.Bot.Delete(context.Message())
	return context.Send(fmt.Sprintf("<code>%v %v</code>", strings.Replace(utils.UserFullName(context.Sender()), "💥", "", -1), context.Data()))
}
