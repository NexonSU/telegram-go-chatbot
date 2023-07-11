package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Unban user on /unban
func Unban(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return utils.ReplyAndRemove("Пример использования: <code>/unban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unban</code>", context)
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	if utils.Bot.Me.ID == target.ID {
		return context.Reply(&tele.Animation{File: tele.File{FileID: "CgACAgQAAx0CQvXPNQABH62yYQHUkpaPOe79NW4ZnwYZWCNJXW8AAgoBAAK-qkVQnRXXGK03dEMgBA"}})
	}
	err = utils.Bot.Unban(context.Chat(), &target)
	if err != nil {
		return err
	}
	return utils.ReplyAndRemove(fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> разбанен.", target.ID, utils.UserFullName(&target)), context)
}
