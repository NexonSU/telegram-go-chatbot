package commands

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Unban user on /unban
func Unban(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply("Пример использования: <code>/unban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unban</code>")
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	if utils.Bot.Me.ID == target.ID {
		return context.Reply(&tele.Animation{File: tele.File{FileID: "CgACAgQAAx0CQvXPNQABH62yYQHUkpaPOe79NW4ZnwYZWCNJXW8AAgoBAAK-qkVQnRXXGK03dEMgBA"}})
	}
	err = utils.Bot.Unban(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка разбана пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Reply(fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> разбанен.", target.ID, utils.UserFullName(&target)))
}
