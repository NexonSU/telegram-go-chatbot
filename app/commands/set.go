package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm/clause"
	"strings"
)

//Save Get to DB on /set
func Set(m *tb.Message) {
	var get utils.Get
	var text = strings.Split(m.Text, " ")
	if (m.ReplyTo == nil && len(text) < 3) || (m.ReplyTo != nil && len(text) != 2) {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/set {гет} {значение}</code>\nИли отправь в ответ на какое-либо сообщение <code>/set {гет}</code>")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	get.Name = strings.ToLower(text[1])
	if m.ReplyTo == nil && len(text) > 2 {
		get.Type = "Text"
		get.Data = strings.Join(text[2:], " ")
	} else if m.ReplyTo != nil && len(text) == 2 {
		get.Caption = m.ReplyTo.Caption
		switch {
		case m.ReplyTo.Animation != nil:
			get.Type = "Animation"
			get.Data = m.ReplyTo.Animation.FileID
		case m.ReplyTo.Audio != nil:
			get.Type = "Audio"
			get.Data = m.ReplyTo.Audio.FileID
		case m.ReplyTo.Photo != nil:
			get.Type = "Photo"
			get.Data = m.ReplyTo.Photo.FileID
		case m.ReplyTo.Video != nil:
			get.Type = "Video"
			get.Data = m.ReplyTo.Video.FileID
		case m.ReplyTo.Voice != nil:
			get.Type = "Voice"
			get.Data = m.ReplyTo.Voice.FileID
		case m.ReplyTo.Document != nil:
			get.Type = "Document"
			get.Data = m.ReplyTo.Document.FileID
		case m.ReplyTo.Text != "":
			get.Type = "Text"
			get.Data = m.ReplyTo.Text
		default:
			_, err := utils.Bot.Reply(m, "Не удалось распознать файл в сообщении, возможно, он не поддерживается.")
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
	}
	result := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(get)
	if result.Error != nil {
		utils.ErrorReporting(result.Error, m)
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось сохранить гет <code>%v</code>.", get.Name))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err := utils.Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> сохранён как <code>%v</code>.", get.Name, get.Type))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
