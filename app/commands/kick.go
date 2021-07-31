package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Kick user on /kick
func Kick(context telebot.Context) error {
	var err error
	if !utils.IsAdminOrModer(context.Sender().Username) {
		if context.Chat().Username != utils.Config.Telegram.Chat {
			return err
		}
		err := context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			return err
		}
		return err
	}
	var text = strings.Split(context.Text(), " ")
	if (context.Message().ReplyTo == nil && len(text) == 1) || (context.Message().ReplyTo != nil && len(text) != 2) {
		err := context.Reply("Пример использования: <code>/kick {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/kick</code>")
		if err != nil {
			return err
		}
		return err
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		err := context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			return err
		}
		return err
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		err := context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
		if err != nil {
			return err
		}
		return err
	}
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Ban(context.Chat(), TargetChatMember)
	if err != nil {
		err := context.Reply(fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			return err
		}
		return err
	}
	err = utils.Bot.Unban(context.Chat(), &target)
	if err != nil {
		err := context.Reply(fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			return err
		}
		return err
	}
	err = context.Reply(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> исключен.", target.ID, utils.UserFullName(&target)))
	if err != nil {
		return err
	}
	return err
}
