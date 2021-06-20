package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

//Send formatted text on /me
func Me(m *tb.Message) {
	var text = strings.Split(m.Text, " ")
	if len(text) == 1 {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Пример использования:\n<code>/me {делает что-то}</code>"))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	err := utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	_, err = utils.Bot.Send(m.Chat, fmt.Sprintf("<code>%v %v</code>", utils.UserFullName(m.Sender), strings.Join(text[1:], " ")))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
