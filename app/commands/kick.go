package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Kick user on /kick
func Kick(context telebot.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply("Пример использования: <code>/kick {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/kick</code>")
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
	}
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Ban(context.Chat(), TargetChatMember)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
	}
	err = utils.Bot.Unban(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
	}
	return context.Reply(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> исключен.", target.ID, utils.UserFullName(&target)))
}
