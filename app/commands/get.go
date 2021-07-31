package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send Get to user on /get
func Get(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	var get utils.Get
	var text = strings.Split(context.Text(), " ")
	if len(text) != 2 {
		return context.Reply("Пример использования: <code>/get {гет}</code>")
	}
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(text[1])}).First(&get)
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
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
	}
}
