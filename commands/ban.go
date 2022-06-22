package commands

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Ban user on /ban
func Ban(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) > 1) {
		return context.Reply("Пример использования: <code>/ban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/ban</code>\nЕсли нужно забанить на время, то добавь время в секундах через пробел.")
	}
	target, untildate, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя или время бана:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember.RestrictedUntil = untildate
	err = utils.Bot.Ban(context.Chat(), TargetChatMember)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка бана пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Reply(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> забанен%v.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)))
}
