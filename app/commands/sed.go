package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"os/exec"
	"strings"
)

//Replace text in target message and send result on /sed
func Sed(m *tb.Message) {
	var text = strings.Split(m.Text, " ")
	if m.ReplyTo != nil {
		cmd := fmt.Sprintf("echo \"%v\" | sed \"%v\"", strings.ReplaceAll(m.ReplyTo.Text, "\"", "\\\""), strings.ReplaceAll(text[1], "\"", "\\\""))
		out, err := exec.Command("bash", "-c", cmd).Output()
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		_, err = utils.Bot.Reply(m, string(out))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else {
		_, err := utils.Bot.Reply(m, "Пример использования:\n/sed {патерн вида s/foo/bar/} в ответ на сообщение.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
