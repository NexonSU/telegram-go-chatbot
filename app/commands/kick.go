package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
	"time"
)

//Kick user on /kick
func Kick(m *tb.Message) {
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
	var text = strings.Split(m.Text, " ")
	if (m.ReplyTo == nil && len(text) == 1) || (m.ReplyTo != nil && len(text) != 2) {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/kick {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/kick</code>")
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
	TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
	err = utils.Bot.Ban(m.Chat, TargetChatMember)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	err = utils.Bot.Unban(m.Chat, &target)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err = utils.Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> исключен.", target.ID, utils.UserFullName(&target)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
