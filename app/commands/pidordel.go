package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Remove user in DB on /pidordel
func Pidordel(context telebot.Context) error {
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
	var user telebot.User
	var pidor utils.PidorList
	user, _, err := utils.FindUserInMessage(context)
	if err != nil {
		err := context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			return err
		}
		return err
	}
	pidor = utils.PidorList(user)
	result := utils.DB.Delete(&pidor)
	if result.RowsAffected != 0 {
		err := context.Reply(fmt.Sprintf("Пользователь %v удалён из игры <b>Пидор Дня</b>!", utils.MentionUser(&user)))
		if err != nil {
			return err
		}
	} else {
		err := context.Reply(fmt.Sprintf("Не удалось удалить пользователя:\n<code>%v</code>", result.Error.Error()))
		if err != nil {
			return err
		}
	}
	return err
}
