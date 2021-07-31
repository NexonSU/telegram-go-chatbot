package commands

import (
	"strings"

	"gopkg.in/tucnak/telebot.v3"
)

// Sed Replace text in target message
func Sed(context telebot.Context) error {
	var text = strings.Split(context.Text(), " ")
	var foo = strings.Split(text[1], "/")[1]
	var bar = strings.Split(text[1], "/")[2]
	if context.Message().ReplyTo == nil || foo == "" || bar == "" {
		return context.Reply("Пример использования:\n/sed {патерн вида s/foo/bar/} в ответ на сообщение.")
	}
	return context.Reply(strings.ReplaceAll(context.Message().ReplyTo.Text, foo, bar))
}
