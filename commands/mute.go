package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Mute user on /mute
func Mute(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) > 1) {
		return utils.ReplyAndRemove("Пример использования: <code>/mute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/mute</code>\nЕсли нужно замьютить на время, то добавь время в секундах через пробел.", context)
	}
	target, untildate, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return err
	}
	TargetChatMember.CanSendMessages = false
	TargetChatMember.RestrictedUntil = untildate
	if utils.Bot.Restrict(context.Chat(), TargetChatMember) != nil {
		return utils.ReplyAndRemove(fmt.Sprintf("Ошибка ограничения пользователя:\n<code>%v</code>", err.Error()), context)
	}
	return utils.ReplyAndRemove(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> больше не может отправлять сообщения%v.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)), context)
}
