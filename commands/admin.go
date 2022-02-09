package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"gopkg.in/telebot.v3"
)

//Send admin list to user on /admin
func Admin(context telebot.Context) error {
	var get utils.Get
	result := utils.DB.Where(&utils.Get{Name: "admin"}).First(&get)
	if result.RowsAffected != 0 {
		switch {
		case get.Type == "Animation":
			return context.Reply(&telebot.Animation{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Audio":
			return context.Reply(&telebot.Audio{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Photo":
			return context.Reply(&telebot.Photo{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Video":
			return context.Reply(&telebot.Video{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Voice":
			return context.Reply(&telebot.Voice{
				File: telebot.File{FileID: get.Data},
			})
		case get.Type == "Document":
			return context.Reply(&telebot.Document{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Text":
			return context.Reply(get.Data)
		default:
			return context.Reply(fmt.Sprintf("Ошибка при определении типа гета, я не знаю тип <code>%v</code>.", get.Type))
		}
	} else {
		return context.Reply("Гет <code>admin</code> не найден.")
	}
}
