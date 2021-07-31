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
	return context.Send(fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(text[1:], " "))), &telebot.SendOptions{DisableWebPagePreview: true, ReplyTo: context.Message().ReplyTo, AllowWithoutReply: true})
}
