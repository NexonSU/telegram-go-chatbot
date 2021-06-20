package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"net/url"
	"strings"
)

//Reply google URL on "google"
func Google(m *tb.Message) {
	var target = *m
	var text = strings.Split(m.Text, " ")
	if len(text) == 1 {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Пример использования:\n<code>/google {запрос}</code>"))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	if m.ReplyTo != nil {
		target = *m.ReplyTo
	}
	_, err := utils.Bot.Reply(&target, fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(text[1:], " "))))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
