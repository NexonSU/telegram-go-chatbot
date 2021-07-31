package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send admin list to user on /admin
func Admin(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	var get utils.Get
	result := utils.DB.Where(&utils.Get{Name: "admin"}).First(&get)
	if result.RowsAffected != 0 {
		switch {
		case get.Type == "Animation":
			err = context.Reply(&telebot.Animation{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
			return err
		case get.Type == "Audio":
			err = context.Reply(&telebot.Audio{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
			return err
		case get.Type == "Photo":
			err = context.Reply(&telebot.Photo{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
			return err
		case get.Type == "Video":
			err = context.Reply(&telebot.Video{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
			return err
		case get.Type == "Voice":
			err = context.Reply(&telebot.Voice{
				File: telebot.File{FileID: get.Data},
			})
			return err
		case get.Type == "Document":
			err = context.Reply(&telebot.Document{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
			return err
		case get.Type == "Text":
			err = context.Reply(get.Data)
			return err
		default:
			err = context.Reply(fmt.Sprintf("Ошибка при определении типа гета, я не знаю тип <code>%v</code>.", get.Type))
			return err
		}
	} else {
		err = context.Reply("Гет <code>admin</code> не найден.")
		return err
	}
}
