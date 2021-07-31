package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

//Save Get to DB on /set
func Set(context telebot.Context) error {
	var get utils.Get
	if (context.Message().ReplyTo == nil && len(context.Args()) < 2) || (context.Message().ReplyTo != nil && len(context.Args()) != 1) {
		return context.Reply("Пример использования: <code>/set {гет} {значение}</code>\nИли отправь в ответ на какое-либо сообщение <code>/set {гет}</code>")
	}
	get.Name = strings.ToLower(context.Args()[0])
	if context.Message().ReplyTo == nil && len(context.Args()) > 1 {
		get.Type = "Text"
		get.Data = strings.Join(context.Args()[1:], " ")
	} else if context.Message().ReplyTo != nil && len(context.Args()) == 1 {
		get.Caption = context.Message().ReplyTo.Caption
		switch {
		case context.Message().ReplyTo.Animation != nil:
			get.Type = "Animation"
			get.Data = context.Message().ReplyTo.Animation.FileID
		case context.Message().ReplyTo.Audio != nil:
			get.Type = "Audio"
			get.Data = context.Message().ReplyTo.Audio.FileID
		case context.Message().ReplyTo.Photo != nil:
			get.Type = "Photo"
			get.Data = context.Message().ReplyTo.Photo.FileID
		case context.Message().ReplyTo.Video != nil:
			get.Type = "Video"
			get.Data = context.Message().ReplyTo.Video.FileID
		case context.Message().ReplyTo.Voice != nil:
			get.Type = "Voice"
			get.Data = context.Message().ReplyTo.Voice.FileID
		case context.Message().ReplyTo.Document != nil:
			get.Type = "Document"
			get.Data = context.Message().ReplyTo.Document.FileID
		case context.Message().ReplyTo.Text != "":
			get.Type = "Text"
			get.Data = context.Message().ReplyTo.Text
		default:
			return context.Reply("Не удалось распознать файл в сообщении, возможно, он не поддерживается.")
		}
	}
	get.Creator = context.Sender().ID
	result := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(get)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Не удалось сохранить гет <code>%v</code>.", get.Name))
	}
	return context.Reply(fmt.Sprintf("Гет <code>%v</code> сохранён как <code>%v</code>.", get.Name, get.Type))
}
