package commands

import (
	cntx "context"
	"time"

	"github.com/wader/goutubedl"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Download(context tele.Context) error {
	if context.Message().ReplyTo == nil && len(context.Args()) < 1 {
		return context.Reply("Пример использования: <code>/download {ссылка на ютуб/твиттер}</code>\nИли отправь в ответ на какое-либо сообщение с ссылкой <code>/download</code>")
	}

	link := ""
	message := &tele.Message{}

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
		if entity.Type == tele.EntityURL {
			text := message.EntityText(entity)
			link = text
		}
	}

	if link == "" {
		return context.Reply("Ссылка ненайдена.")
	}

	goutubedl.Path = "yt-dlp"

	result, err := goutubedl.New(cntx.Background(), link, goutubedl.Options{})
	if err != nil {
		return err
	}

	if result.Info.Duration > 3600 {
		return context.Reply("Максимальная длина видео 60 минут.")
	}

	ytdlpResult, err := result.Download(cntx.Background(), "bv*[ext=mp4]+ba[ext=m4a]/b[ext=mp4] / bv*+ba/b")
	if err != nil {
		return err
	}
	defer ytdlpResult.Close()

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

	return context.Reply(&tele.Document{File: tele.FromReader(ytdlpResult), FileName: result.Info.Title + ".mp4"})
}
