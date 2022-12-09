package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Adds bless text to DB
func AddBless(context tele.Context) error {
	var bless utils.Bless
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply("Пример использования: <code>/addbless {текст}</code>\nИли отправь в ответ на сообщение с текстом <code>/addbless</code>")
	}
	if context.Message().ReplyTo == nil {
		bless.Text = context.Data()
	} else {
		if context.Message().ReplyTo.Text != "" {
			bless.Text = context.Message().ReplyTo.Text
		} else {
			return context.Reply("Я не смог найти текст в указанном сообщении.")
		}
	}
	if len([]rune(bless.Text)) > 200 {
		return context.Reply("Bless не может быть длиннее 200 символов.")
	}
	result := utils.DB.Create(&bless)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Не удалось добавить bless, ошибка:\n<code>%v</code>", result.Error.Error()))
	}
	return context.Reply(fmt.Sprintf("Bless добавлен как <code>%v</code>.", bless.Text))
}
