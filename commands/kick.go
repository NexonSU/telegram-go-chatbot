package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Kick user on /kick
func Kick(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return utils.SendAndRemove("Пример использования: <code>/kick {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/kick</code>", context)
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return err
	}
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Ban(context.Chat(), TargetChatMember)
	if err != nil {
		return err
	}
	err = utils.Bot.Unban(context.Chat(), &target)
	if err != nil {
		return err
	}
	return utils.SendAndRemove(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> исключен.", target.ID, utils.UserFullName(&target)), context)
}
