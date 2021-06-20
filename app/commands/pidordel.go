package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Remove user in DB on /pidordel
func Pidordel(m *tb.Message) {
	if !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Admins) && !utils.StringInSlice(m.Sender.Username, utils.Config.Telegram.Moders) {
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var user tb.User
	var pidor utils.PidorList
	user, _, err := utils.FindUserInMessage(*m)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	pidor = utils.PidorList(user)
	result := utils.DB.Delete(&pidor)
	if result.RowsAffected != 0 {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Пользователь %v удалён из игры <b>Пидор Дня</b>!", utils.MentionUser(&user)))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось удалить пользователя:\n<code>%v</code>", result.Error.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
