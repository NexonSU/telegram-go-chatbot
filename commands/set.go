package commands

import (
	"fmt"
	"strings"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	"gorm.io/gorm/clause"
)

//Save Get to DB on /set
func Set(context tele.Context) error {
	var get utils.Get
	//args check
	if (context.Message().ReplyTo == nil && len(context.Args()) < 2) || (context.Message().ReplyTo != nil && len(context.Args()) == 0) {
		return context.Reply("Пример использования: <code>/set {гет} {значение}</code>\nИли отправь в ответ на какое-либо сообщение <code>/set {гет}</code>")
	}
	//ownership check
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(context.Args()[0])}).First(&get)
	if result.RowsAffected != 0 {
		creator, err := utils.GetUserFromDB(fmt.Sprint(get.Creator))
		if err != nil {
			return err
		}
		if get.Creator != context.Sender().ID && !utils.IsAdminOrModer(context.Sender().ID) {
			return context.Reply(fmt.Sprintf("Данный гет могут изменять либо администраторы, либо %v.", utils.UserFullName(&creator)))
		}
	}
	//filling Get from message
	if context.Message().ReplyTo == nil {
		get.Name = strings.ToLower(context.Args()[0])
		get.Title = context.Args()[0]
		get.Type = "Text"
		get.Data = utils.GetHtmlText(*context.Message())
	} else {
		get.Name = strings.ToLower(context.Data())
		get.Title = context.Data()
		get.Caption = utils.GetHtmlText(*context.Message().ReplyTo)
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
			get.Data = utils.GetHtmlText(*context.Message().ReplyTo)
		default:
			return context.Reply("Не удалось распознать файл в сообщении, возможно, он не поддерживается.")
		}
	}
	get.Creator = context.Sender().ID
	//writing get to DB
	result = utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&get)
	if result.Error != nil {
		return result.Error
	}
	return context.Reply(fmt.Sprintf("Гет <code>%v</code> сохранён как <code>%v</code>.", get.Name, get.Type))
}
