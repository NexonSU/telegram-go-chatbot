package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Delete Get in DB on /del
func Del(context telebot.Context) error {
	if len(context.Args()) != 1 {
		return context.Reply("Пример использования: <code>/del {гет}</code>")
	}
	result := utils.DB.Delete(&utils.Get{Name: strings.ToLower(context.Data())})
	if result.RowsAffected != 0 {
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> удалён.", context.Data()))
	} else {
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> не найден.", context.Data()))
	}
}
