package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

//Mute user on /mute
func Mute(m *tb.Message) {
	if !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Admins) && !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Moders) {
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var text = strings.Split(m.Text, " ")
	if (m.ReplyTo == nil && len(text) < 2) || (m.ReplyTo != nil && len(text) > 2) {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/mute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/mute</code>\nЕсли нужно замьютить на время, то добавь время в секундах через пробел.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	target, untildate, err := utils.FindUserInMessage(*m)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя или время ограничения:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	TargetChatMember, err := utils.Bot.ChatMemberOf(m.Chat, &target)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	TargetChatMember.CanSendMessages = false
	TargetChatMember.RestrictedUntil = untildate
	err = utils.Bot.Restrict(m.Chat, TargetChatMember)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка ограничения пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err = utils.Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> больше не может отправлять сообщения%v.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
