package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

//Delete Get in DB on /del
func Del(m *tb.Message) {
	var text = strings.Split(m.Text, " ")
	if len(text) != 2 {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/del {гет}</code>")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	result := utils.DB.Delete(&utils.Get{Name: strings.ToLower(text[1])})
	if result.RowsAffected != 0 {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> удалён.", text[1]))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
