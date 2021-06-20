package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

//Send Get to user on /get
func Get(m *tb.Message) {
	var get utils.Get
	var text = strings.Split(m.Text, " ")
	if len(text) != 2 {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/get {гет}</code>")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(text[1])}).First(&get)
	if result.RowsAffected != 0 {
		switch {
		case get.Type == "Animation":
			_, err := utils.Bot.Reply(m, &tb.Animation{
				File:    tb.File{FileID: get.Data},
				Caption: get.Caption,
			})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		case get.Type == "Audio":
			_, err := utils.Bot.Reply(m, &tb.Audio{
				File:    tb.File{FileID: get.Data},
				Caption: get.Caption,
			})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		case get.Type == "Photo":
			_, err := utils.Bot.Reply(m, &tb.Photo{
				File:    tb.File{FileID: get.Data},
				Caption: get.Caption,
			})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		case get.Type == "Video":
			_, err := utils.Bot.Reply(m, &tb.Video{
				File:    tb.File{FileID: get.Data},
				Caption: get.Caption,
			})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		case get.Type == "Voice":
			_, err := utils.Bot.Reply(m, &tb.Voice{
				File: tb.File{FileID: get.Data},
			})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		case get.Type == "Document":
			_, err := utils.Bot.Reply(m, &tb.Document{
				File:    tb.File{FileID: get.Data},
				Caption: get.Caption,
			})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		case get.Type == "Text":
			_, err := utils.Bot.Reply(m, get.Data)
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		default:
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка при определении типа гета, я не знаю тип <code>%v</code>.", get.Type))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		}
	} else {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
