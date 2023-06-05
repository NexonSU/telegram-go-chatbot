package commands

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
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
	var extension string
	var targetArg string

	switch media.MediaType() {
	case "audio":
		extension = "mp3"
	case "voice":
		extension = "ogg"
	case "photo":
		extension = "jpg"
	case "sticker":
		extension = "webp"
	case "animation", "video", "video_note", "document":
		extension = "mp4"
	}

	if media.MediaType() == "sticker" {
		if context.Message().ReplyTo.Sticker.Animated || context.Message().ReplyTo.Sticker.Video {
			extension = "webm"
		}
	}

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

	filePath := fmt.Sprintf("%v/%v.%v", os.TempDir(), media.MediaFile().FileID, extension)

	file, err := utils.Bot.FileByID(media.MediaFile().FileID)
	if err != nil {
		return err
	}

	MarshalledMessage, _ := json.MarshalIndent(file, "", "    ")
	JsonMessage := html.EscapeString(string(MarshalledMessage))
	return context.Reply(fmt.Sprintf("\n\nMessage:\n<pre>%v</pre>", JsonMessage))

	return nil
	return utils.FFmpegConvert(context, filePath, targetArg)
}
