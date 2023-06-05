package commands

import (
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Convert(context tele.Context) error {
	if context.Message().ReplyTo == nil {
		for _, entity := range context.Message().Entities {
			if entity.Type == tele.EntityURL {
				return Download(context)
			}
		}
		return context.Reply("Пример использования: <code>/convert</code> в ответ на какое-либо сообщение с медиа-файлом.\nДопольнительные параметры: gif,mp3,ogg,jpg.")
	}
	if context.Message().ReplyTo.Media() == nil {
		for _, entity := range context.Message().ReplyTo.Entities {
			if entity.Type == tele.EntityURL {
				return Download(context)
			}
		}
		return context.Reply("Какого-либо медиа файла нет в указанном сообщении.")
	}

	media := context.Message().ReplyTo.Media()
	var targetArg string

	targetArg = media.MediaType()
	if len(context.Args()) == 1 {
		targetArg = strings.ToLower(context.Args()[0])
	}

	var done = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				context.Notify(tele.ChatAction(tele.UploadingDocument))
			}
			time.Sleep(time.Second * 5)
		}
	}()
	defer func() {
		done <- true
	}()

	file, err := utils.Bot.FileByID(media.MediaFile().FileID)
	if err != nil {
		return err
	}

	return utils.FFmpegConvert(context, file.FilePath, targetArg)
}
