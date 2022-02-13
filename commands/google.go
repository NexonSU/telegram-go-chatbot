package commands

import (
	"fmt"
	"net/url"

	tele "gopkg.in/telebot.v3"
)

//Reply google URL on "google"
func Google(context tele.Context) error {
	if len(context.Args()) == 0 {
		return context.Reply("Пример использования:\n<code>/google {запрос}</code>")
	}
	return context.Send(fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(context.Data())), &tele.SendOptions{DisableWebPagePreview: true, ReplyTo: context.Message().ReplyTo, AllowWithoutReply: true})
}
