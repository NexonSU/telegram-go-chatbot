package stats

import (
	"fmt"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func RemoveWord(context telebot.Context) error {
	if len(context.Args()) != 1 {
		return context.Reply("Укажите слово.")
	}
	result := utils.DB.Where("text = ?", context.Data()).Delete(&utils.Word{})
	if result.RowsAffected != 0 {
		return context.Reply("Слово удалено.")
	} else {
		return context.Reply(fmt.Sprintf("Не удалось удалить слово:\n<code>%v</code>", result.Error.Error()))
	}
}
