package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

//Send warning to user on /warn
func Warn(context telebot.Context) error {
	var warn utils.Warn
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return context.Reply("Пример использования: <code>/warn {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/warn</code>")
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
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
	}).Create(&warn)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Не удалось выдать предупреждение:\n<code>%v</code>.", result.Error))
	}
	if warn.Amount == 1 {
		return context.Send(fmt.Sprintf("%v, у тебя 1 предупреждение.\nЕсль получишь 3 предупреждения за 2 недели, то будешь исключен из чата.", utils.UserFullName(&target)))
	}
	if warn.Amount == 2 {
		return context.Send(fmt.Sprintf("%v, у тебя 2 предупреждения.\nЕсли в течении недели получишь ещё одно, то будешь исключен из чата.", utils.UserFullName(&target)))
	}
	if warn.Amount == 3 {
		untildate := time.Now().AddDate(0, 0, 7).Unix()
		TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
		if err != nil {
			return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
		}
		TargetChatMember.RestrictedUntil = untildate
		err = utils.Bot.Ban(context.Chat(), TargetChatMember)
		if err != nil {
			return context.Reply(fmt.Sprintf("Ошибка бана пользователя:\n<code>%v</code>", err.Error()))
		}
		return context.Reply(fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> забанен%v, т.к. набрал 3 предупреждения.", target.ID, utils.UserFullName(&target), utils.RestrictionTimeMessage(untildate)))
	}
	return err
}
