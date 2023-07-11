package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Unmute user on /unmute
func Revive(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return utils.ReplyAndRemove("Пример использования: <code>/unmute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unmute</code>", context)
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return err
	}
	TargetChatMember.CanSendMessages = true
	TargetChatMember.CanSendMedia = true
	TargetChatMember.CanSendPolls = true
	TargetChatMember.CanSendOther = true
	TargetChatMember.CanAddPreviews = true
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Restrict(context.Chat(), TargetChatMember)
	if err != nil {
		return err
	}
	return utils.ReplyAndRemove(fmt.Sprintf("%v возродился в чате.", utils.MentionUser(&target)), context)
}
