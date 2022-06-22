package commands

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Mute user on /mute
func Mute(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) > 1) {
		return context.Reply("Пример использования: <code>/mute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/mute</code>\nЕсли нужно замьютить на время, то добавь время в секундах через пробел.")
	}
	target, untildate, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя или время ограничения:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember.CanSendMessages = false
	TargetChatMember.RestrictedUntil = untildate
	if utils.Bot.Restrict(context.Chat(), TargetChatMember) != nil {
		return context.Reply(fmt.Sprintf("Ошибка ограничения пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Reply(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> больше не может отправлять сообщения%v.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)))
}
