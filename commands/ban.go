package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Ban user on /ban
func Ban(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) > 1) {
		return utils.ReplyAndRemove("Пример использования: <code>/ban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/ban</code>\nЕсли нужно забанить на время, то добавь время в секундах через пробел.", context)
	}
	target, untildate, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return err
	}
	TargetChatMember.RestrictedUntil = untildate
	err = utils.Bot.Ban(context.Chat(), TargetChatMember)
	if err != nil {
		return err
	}
	return utils.ReplyAndRemove(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> забанен%v.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)), context)
}
