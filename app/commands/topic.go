package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

//Change chat name on /topic
func Topic(m *tb.Message) {
	if !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Admins) && !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Moders) {
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var text = strings.Split(m.Text, " ")
	if len(text) < 2 {
		_, err := utils.Bot.Reply(m, "Пример использования:\n<code>/topic {новая тема чата}</code>")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	err := utils.Bot.SetGroupTitle(m.Chat, fmt.Sprintf("Zavtrachat | %v", strings.Join(text[1:], " ")))
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка изменения названия чата:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}

		return
	}
}
