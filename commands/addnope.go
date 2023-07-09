package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Adds nope text to DB
func AddNope(context tele.Context) error {
	var nope utils.Nope
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return utils.SendAndRemove("Пример использования: <code>/addnope {текст}</code>\nИли отправь в ответ на сообщение с текстом <code>/addnope</code>", context)
	}
	if context.Message().ReplyTo == nil {
		nope.Text = strings.TrimPrefix(context.Text(), strings.Split(context.Text(), " ")[0]+" ")
	} else {
		if context.Message().ReplyTo.Text != "" {
			nope.Text = strings.ToLower(context.Message().ReplyTo.Text)
		} else {
			return utils.SendAndRemove("Я не смог найти текст в указанном сообщении.", context)
		}
	}
	result := utils.DB.Create(&nope)
	if result.Error != nil {
		return utils.SendAndRemove(fmt.Sprintf("Не удалось добавить nope, ошибка:\n<code>%v</code>", result.Error.Error()), context)
	}
	return utils.SendAndRemove(fmt.Sprintf("Nope добавлен как <code>%v</code>.", nope.Text), context)
}
