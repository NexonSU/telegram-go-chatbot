package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Delete Get in DB on /del
func Del(context telebot.Context) error {
	var text = strings.Split(context.Text(), " ")
	if len(text) != 2 {
		return context.Reply("Пример использования: <code>/del {гет}</code>")
	}
	result := utils.DB.Delete(&utils.Get{Name: strings.ToLower(text[1])})
	if result.RowsAffected != 0 {
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> удалён.", text[1]))
	} else {
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
	}
}
