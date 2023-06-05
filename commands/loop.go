package commands

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Invert given file
func Loop(context tele.Context) error {
	if context.Message().ReplyTo == nil {
		return context.Reply("Пример использования: <code>/loop</code> в ответ на какое-либо сообщение с видео.")
	}
	if context.Message().ReplyTo.Media() == nil {
		return context.Reply("Какого-либо видео нет в указанном сообщении.")
	}

	media := context.Message().ReplyTo.Media()

	targetArg := media.MediaType()
	if len(context.Args()) == 1 {
		targetArg = strings.ToLower(context.Args()[0])
	}

	var extension string
	switch targetArg {
	case "animation":
		extension = "mp4"
		targetArg = "animation"
	default:
		return context.Reply("Неподдерживаемая операция")
	}

	targetArg = targetArg + "_loop"

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

	err := utils.Bot.Download(media.MediaFile(), filePath)
	if err != nil {
		return err
	}

	return utils.FFmpegConvert(context, filePath, targetArg)
}
