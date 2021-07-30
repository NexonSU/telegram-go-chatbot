package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send formatted text on /me
func Me(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var text = strings.Split(m.Text, " ")
	if len(text) == 1 {
		_, err := utils.Bot.Reply(m, "Пример использования:\n<code>/me {делает что-то}</code>")
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
