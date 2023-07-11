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
		return utils.ReplyAndRemove("Пример использования: <code>/addbless {текст}</code>\nИли отправь в ответ на сообщение с текстом <code>/addbless</code>", context)
	}
	if context.Message().ReplyTo == nil {
		bless.Text = context.Data()
	} else {
		if context.Message().ReplyTo.Text != "" {
			bless.Text = context.Message().ReplyTo.Text
		} else {
			return utils.ReplyAndRemove("Я не смог найти текст в указанном сообщении.", context)
		}
	}
	if len([]rune(bless.Text)) > 200 {
		return utils.ReplyAndRemove("Bless не может быть длиннее 200 символов.", context)
	}
	result := utils.DB.Create(&bless)
	if result.Error != nil {
		return result.Error
	}
	return utils.ReplyAndRemove(fmt.Sprintf("Bless добавлен как <code>%v</code>.", bless.Text), context)
}
