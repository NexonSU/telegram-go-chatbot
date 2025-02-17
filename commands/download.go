package commands

import (
	cntx "context"
	"fmt"
	"os"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/lrstanley/go-ytdlp"
	tele "gopkg.in/telebot.v3"
)

// Convert given  file
func Download(context tele.Context) error {
	var filePath string

	context.Delete()

	if context.Message().ReplyTo == nil && len(context.Args()) < 1 {
		return utils.ReplyAndRemove("Пример использования: <code>/download {ссылка на ютуб/твиттер}</code>\nИли отправь в ответ на какое-либо сообщение с ссылкой <code>/download</code>", context)
	}

	link := ""
	message := &tele.Message{}
	arg := "video"

	if context.Message().ReplyTo == nil && len(context.Args()) == 2 {
		arg = context.Args()[1]
	}

	if context.Message().ReplyTo != nil && len(context.Args()) == 1 {
		arg = context.Args()[0]
	}

	if context.Message().ReplyTo == nil {
		message = context.Message()
	} else {
		message = context.Message().ReplyTo
	}

	var done = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				context.Notify(tele.RecordingVideo)
			}
			time.Sleep(time.Second * 5)
		}
	}()
	defer func() {
		done <- true
	}()

	for _, entity := range message.Entities {
		if entity.Type == tele.EntityURL || entity.Type == tele.EntityTextLink {
			link = entity.URL
			if link == "" {
				link = message.EntityText(entity)
			}
		}
	}

	ytdlp.MustInstall(cntx.TODO(), nil)

	filePath = fmt.Sprintf("%v/%v.mp4", os.TempDir(), context.Message().ID)

	dl := ytdlp.New().FormatSort("res,ext:mp4:m4a").RecodeVideo("mp4").Output(filePath)

	_, err := dl.Run(cntx.TODO(), link)
	if err != nil {
		return err
	}

	done <- true

	var done2 = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done2:
				return
			default:
				context.Notify(tele.UploadingVideo)
			}
			time.Sleep(time.Second * 5)
		}
	}()
	defer func() {
		done2 <- true
	}()

	return utils.FFmpegConvert(context, filePath, arg)
}
