package commands

import (
	"fmt"
	"net/url"
	"strings"

	"gopkg.in/tucnak/telebot.v3"
)

//Reply google URL on "google"
func Google(context telebot.Context) error {
	var text = strings.Split(context.Text(), " ")
	if len(text) == 1 {
		return context.Reply("Пример использования:\n<code>/google {запрос}</code>")
	}
	if context.Message().ReplyTo != nil {
		context.Message().Sender = context.Message().ReplyTo.Sender
	}
	return context.Reply(fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(text[1:], " "))))
}
