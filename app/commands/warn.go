package commands

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm/clause"
)

//Send warning to user on /warn
func Warn(m *tb.Message) {
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
	var warn utils.Warn
	var text = strings.Split(m.Text, " ")
	if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/warn {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/warn</code>")
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
	result := utils.DB.First(&warn, target.ID)
	if result.RowsAffected != 0 {
		warn.Amount = warn.Amount - int(time.Since(warn.LastWarn).Hours()/24/7)
		if warn.Amount < 0 {
			warn.Amount = 0
		}
		warn.Amount = warn.Amount + 1
	} else {
		warn.Amount = 1
	}
	warn.UserID = target.ID
	warn.LastWarn = time.Now()
	result = utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(warn)
	if result.Error != nil {
		utils.ErrorReporting(result.Error, m)
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось выдать предупреждение:\n<code>%v</code>.", result.Error))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	if warn.Amount == 1 {
		_, err := utils.Bot.Send(m.Chat, fmt.Sprintf("%v, у тебя 1 предупреждение.\nЕсль получишь 3 предупреждения за 2 недели, то будешь исключен из чата.", utils.MentionUser(&target)))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
	if warn.Amount == 2 {
		_, err := utils.Bot.Send(m.Chat, fmt.Sprintf("%v, у тебя 2 предупреждения.\nЕсли в течении недели получишь ещё одно, то будешь исключен из чата.", utils.MentionUser(&target)))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
	if warn.Amount == 3 {
		untildate := time.Now().AddDate(0, 0, 7).Unix()
		TargetChatMember, err := utils.Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember.RestrictedUntil = untildate
		err = utils.Bot.Ban(m.Chat, TargetChatMember)
		if err != nil {
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка бана пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = utils.Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> забанен%v, т.к. набрал 3 предупреждения.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
