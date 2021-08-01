package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Adds nope text to DB
func AddNope(context telebot.Context) error {
	var nope utils.Nope
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply("Пример использования: <code>/addnope {текст}</code>\nИли отправь в ответ на какое-либо с текстом <code>/addnope</code>")
	}
	if context.Message().ReplyTo == nil {
		nope.Text = strings.ToLower(context.Data())
	} else {
		if context.Message().ReplyTo.Text != "" {
			nope.Text = strings.ToLower(context.Message().ReplyTo.Text)
		} else {
			return context.Reply("Я не смог найти текст в указанном сообщении.")
		}
	}
	if len([]rune(nope.Text)) > 50 {
		return context.Reply("Nope не может быть длиннее 50 символов.")
	}
	result := utils.DB.Create(nope)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Не удалось добавить nope, ошибка:\n<code>%v</code>", result.Error.Error()))
	}
	return context.Reply(fmt.Sprintf("Nope добавлен как <code>%v</code>.", nope.Text))
}
