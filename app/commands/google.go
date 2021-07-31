package commands

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Reply google URL on "google"
func Google(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	var text = strings.Split(context.Text(), " ")
	if len(text) == 1 {
		return context.Reply("Пример использования:\n<code>/google {запрос}</code>")
	}
	if context.Message().ReplyTo != nil {
		context.Message().Sender = context.Message().ReplyTo.Sender
	}
	return context.Reply(fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(text[1:], " "))))
}
