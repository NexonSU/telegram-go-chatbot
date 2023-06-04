package commands

import (
	"bytes"
	"os"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
)

// Invert given file
func Invert(context tele.Context) error {
	if context.Message().ReplyTo == nil {
		return context.Reply("Пример использования: <code>/convert</code> в ответ на какое-либо сообщение с видео.")
	}
	if context.Message().ReplyTo.Media() == nil {
		return context.Reply("Какого-либо видео нет в указанном сообщении.")
	}
	media := context.Message().ReplyTo.Media()
	outputKwArgs := ffmpeg.KwArgs{"map": "0"}
	inputKwArgs := ffmpeg.KwArgs{"map": "0"}
	switch media.MediaType() {
	case "video":
		inputKwArgs = ffmpeg.KwArgs{"loglevel": "debug", "hide_banner": "", "format": "mp4"}
		outputKwArgs = ffmpeg.KwArgs{"format": "mp4", "map": "0", "c:v": "libx264", "preset": "fast", "crf": 26, "movflags": "frag_keyframe+empty_moov+faststart", "c:a": "aac", "vf": "reverse", "af": "areverse"}
	case "animation":
		inputKwArgs = ffmpeg.KwArgs{"loglevel": "debug", "hide_banner": "", "format": "mp4"}
		outputKwArgs = ffmpeg.KwArgs{"format": "mp4", "map": "v:0", "c:v": "libx264", "preset": "fast", "crf": 26, "movflags": "frag_keyframe+empty_moov+faststart", "vf": "reverse"}
	default:
		return context.Reply("Неподдерживаемая операция")
	}

	var done = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				context.Notify(tele.ChatAction("upload_video"))
			}
			time.Sleep(time.Second * 5)
		}
	}()
	defer func() {
		done <- true
	}()

	buf := bytes.NewBuffer(nil)
	fileReader, err := utils.Bot.File(media.MediaFile())
	if err != nil {
		return err
	}
	err = ffmpeg.Input("pipe:", inputKwArgs).Output("pipe:", outputKwArgs).WithInput(fileReader).WithOutput(buf, os.Stdout).Run()
	if err != nil {
		return err
	}

	return context.Reply(&tele.Document{
		File:     tele.FromReader(buf),
		MIME:     "video/mp4",
		FileName: media.MediaFile().FileID + ".mp4",
	}, &tele.SendOptions{AllowWithoutReply: true})
}
