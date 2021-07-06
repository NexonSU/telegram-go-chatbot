package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Unmute user on /unmute
func Revive(m *tb.Message) {
	if !utils.IsAdminOrModer(m.Sender.Username) {
		if m.Chat.Username != utils.Config.Telegram.Chat {
			return
		}
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var target tb.User
	var text = strings.Split(m.Text, " ")
	if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/unmute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unmute</code>")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	target, _, err := utils.FindUserInMessage(*m)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
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
	TargetChatMember.CanSendMessages = true
	TargetChatMember.CanSendMedia = true
	TargetChatMember.CanSendPolls = true
	TargetChatMember.CanSendOther = true
	TargetChatMember.CanAddPreviews = true
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Restrict(m.Chat, TargetChatMember)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка возрождения пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err = utils.Bot.Reply(m, fmt.Sprintf("%v возродился в чате.", utils.MentionUser(&target)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
