package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Unban user on /unban
func Unban(context telebot.Context) error {
	var text = strings.Split(context.Text(), " ")
	if (context.Message().ReplyTo == nil && len(text) != 2) || (context.Message().ReplyTo != nil && len(text) != 1) {
		return context.Reply("Пример использования: <code>/unban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unban</code>")
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	if utils.Bot.Me.ID == target.ID {
		return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgQAAx0CQvXPNQABH62yYQHUkpaPOe79NW4ZnwYZWCNJXW8AAgoBAAK-qkVQnRXXGK03dEMgBA"}})
	}
	err = utils.Bot.Unban(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка разбана пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Reply(fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> разбанен.", target.ID, utils.UserFullName(&target)))
}
