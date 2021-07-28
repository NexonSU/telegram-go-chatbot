package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Unban user on /unban
func Unban(m *tb.Message) {
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
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/unban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unban</code>")
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
	if utils.Bot.Me.ID == target.ID {
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgQAAx0CQvXPNQABH62yYQHUkpaPOe79NW4ZnwYZWCNJXW8AAgoBAAK-qkVQnRXXGK03dEMgBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	err = utils.Bot.Unban(m.Chat, &target)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка разбана пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err = utils.Bot.Reply(m, fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> разбанен.", target.ID, utils.UserFullName(&target)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
