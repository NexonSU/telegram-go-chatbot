package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

// Sed Replace text in target message
func Sed(m *tb.Message) {
	var text = strings.Split(m.Text, " ")
	var foo = strings.Split(text[1], "/")[1]
	var bar = strings.Split(text[1], "/")[2]
	if m.ReplyTo != nil && foo != "" && bar != "" {
		_, err := utils.Bot.Reply(m, strings.ReplaceAll(m.ReplyTo.Text, foo, bar))
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
