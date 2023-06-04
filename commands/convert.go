package commands

import (
	"bytes"
	"os"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Convert(context tele.Context) error {
	if context.Message().ReplyTo == nil {
		return context.Reply("Пример использования: <code>/convert</code> в ответ на какое-либо сообщение с медиа-файлом.\nДопольнительные параметры: gif,mp3,ogg,jpg.")
	}
	if context.Message().ReplyTo.Media() == nil {
		return context.Reply("Какого-либо медиа файла нет в указанном сообщении.")
	}

	media := context.Message().ReplyTo.Media()
	defaultKwArgs := ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": ""}
	KwArgs := ffmpeg.KwArgs{"format": "mp4", "c:v": "libx264", "preset": "fast", "crf": 26, "movflags": "frag_keyframe+empty_moov+faststart", "c:a": "aac"}
	mime := "video/mp4"
	fileName := media.MediaFile().FileID + ".mp4"
	arg := media.MediaType()
	action := "upload_video"
	if len(context.Args()) == 1 {
		arg = strings.ToLower(context.Args()[0])
		switch arg {
		case "mp3", "audio", "ogg", "voice":
			if !utils.StringInSlice(media.MediaType(), []string{"video", "voice", "audio", "document"}) {
				return context.Reply("Неподдерживаемая операция")
			}
		case "jpg", "photo":
			if !utils.StringInSlice(media.MediaType(), []string{"photo", "animation", "video", "document"}) {
				return context.Reply("Неподдерживаемая операция")
			}
		case "gif", "animation":
			if !utils.StringInSlice(media.MediaType(), []string{"video", "animation", "document"}) {
				return context.Reply("Неподдерживаемая операция")
			}
		case "video", "video_note", "document":
			break
		default:
			return context.Reply("Неподдерживаемая операция")
		}
	}
	if arg == "sticker" && (context.Message().ReplyTo.Sticker.Animated || context.Message().ReplyTo.Sticker.Video) {
		arg = "gif"
	}
	switch arg {
	case "mp3", "audio":
		KwArgs = ffmpeg.KwArgs{"map": "a:0", "format": "mp3", "c:a": "libmp3lame"}
		mime = "audio/mp3"
		fileName = media.MediaFile().FileID + ".mp3"
		action = "upload_audio"
	case "ogg", "voice":
		KwArgs = ffmpeg.KwArgs{"map": "a:0", "format": "ogg", "c:a": "libopus"}
		mime = "audio/ogg"
		fileName = media.MediaFile().FileID + ".ogg"
		action = "record_voice"
	case "jpg", "photo", "sticker":
		KwArgs = ffmpeg.KwArgs{"vf": "select=eq(n\\,0)", "format": "image2"}
		mime = "image/jpeg"
		fileName = media.MediaFile().FileID + ".jpg"
		action = "upload_photo"
	case "gif", "animation":
		KwArgs = ffmpeg.KwArgs{"map": "v:0", "format": "mp4", "c:v": "libx264", "an": "", "preset": "fast", "crf": 26, "movflags": "frag_keyframe+empty_moov+faststart"}
	}

	var done = make(chan bool, 1)
	go func() {
		for {
			select {
			case <-done:
				return
			default:
				context.Notify(tele.ChatAction(action))
			}
			time.Sleep(time.Second * 5)
		}
	}()
	defer func() {
		done <- true
	}()

	buf := bytes.NewBuffer(nil)
	err := utils.Bot.Download(media.MediaFile(), "/tmp/"+media.MediaFile().FileID+".mp4")
	if err != nil {
		return err
	}
	err = ffmpeg.Input("/tmp/"+media.MediaFile().FileID+".mp4").Output("pipe:", ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{defaultKwArgs, KwArgs})).WithOutput(buf, os.Stdout).Run()
	if err != nil {
		return err
	}

	os.Remove("/tmp/" + media.MediaFile().FileID + ".mp4")

	return context.Reply(&tele.Document{
		File:     tele.FromReader(buf),
		MIME:     mime,
		FileName: fileName,
	}, &tele.SendOptions{AllowWithoutReply: true})
}
