package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Convert(context tele.Context) error {
	var err error
	var KwArgs map[string]interface{}
	var FileName string
	var Title string
	var Performer string
	var Width int
	var Height int
	var Caption string
	var Duration int
	var Thumbnail *tele.Photo
	var Streaming bool
	var MediaType string
	utils.Bot.URL = "https://api.telegram.org"
	KwArgs = ffmpeg.KwArgs{"loglevel": "debug", "map": "0", "format": "mp4", "c:v": "libx264", "preset": "fast", "crf": 26, "movflags": "frag_keyframe+empty_moov+faststart", "c:a": "aac"}
	RequestedMediaType := ""
	if len(context.Args()) == 1 {
		arg := strings.ToLower(context.Args()[0])
		if arg == "mp3" || arg == "audio" {
			RequestedMediaType = "Audio"
		}
		if arg == "gif" || arg == "animation" {
			RequestedMediaType = "Animation"
		}
		if arg == "ogg" || arg == "voice" {
			RequestedMediaType = "Voice"
		}
		if arg == "mp4" || arg == "video" {
			RequestedMediaType = "Video"
		}
	}
	switch {
	case context.Message().ReplyTo.Audio != nil:
		context.Notify(tele.RecordingAudio)
		Caption = context.Message().ReplyTo.Audio.Caption
		Duration = context.Message().ReplyTo.Audio.Duration
		FileName = context.Message().ReplyTo.Audio.FileName
		Title = context.Message().ReplyTo.Audio.Title
		Performer = context.Message().ReplyTo.Audio.Performer
		MediaType = "Audio"
	case context.Message().ReplyTo.Document != nil && context.Message().ReplyTo.Document.MIME[0:5] == "video":
		context.Notify(tele.RecordingVideo)
		Caption = context.Message().ReplyTo.Document.Caption
		FileName = context.Message().ReplyTo.Document.FileName
		MediaType = "Video"
	case context.Message().ReplyTo.Video != nil:
		context.Notify(tele.RecordingVideo)
		Width = context.Message().ReplyTo.Video.Width
		Height = context.Message().ReplyTo.Video.Height
		Caption = context.Message().ReplyTo.Video.Caption
		Duration = context.Message().ReplyTo.Video.Duration
		Thumbnail = context.Message().ReplyTo.Video.Thumbnail
		Streaming = context.Message().ReplyTo.Video.Streaming
		FileName = context.Message().ReplyTo.Video.FileName
		MediaType = "Video"
	case context.Message().ReplyTo.Voice != nil:
		context.Notify(tele.RecordingAudio)
		Caption = context.Message().ReplyTo.Voice.Caption
		Duration = context.Message().ReplyTo.Voice.Duration
		MediaType = "Voice"
	default:
		return context.Reply("Пример использования: /convert в ответ на какое-либо сообщение с аудио или видео.", &tele.SendOptions{AllowWithoutReply: true})
	}
	if RequestedMediaType == "Animation" {
		MediaType = "Animation"
		KwArgs = ffmpeg.KwArgs{"map": "v:0", "format": "mp4", "c:v": "libx264", "an": "", "preset": "fast", "crf": 26, "movflags": "frag_keyframe+empty_moov+faststart"}
	}
	if RequestedMediaType == "Audio" {
		MediaType = "Audio"
		KwArgs = ffmpeg.KwArgs{"map": "a:0", "format": "mp3", "c:a": "libmp3lame"}
	}
	if RequestedMediaType == "Voice" {
		MediaType = "Voice"
		KwArgs = ffmpeg.KwArgs{"map": "a:0", "format": "ogg", "c:a": "libopus"}
	}
	buf := bytes.NewBuffer(nil)
	fileReader, err := utils.Bot.File(context.Message().ReplyTo.Media().MediaFile())
	if err != nil {
		return err
	}
	err = ffmpeg.Input("pipe:").Output("pipe:", KwArgs).WithInput(fileReader).WithOutput(buf, os.Stdout).Run()
	if err != nil {
		return err
	}
	if MediaType == "Audio" {
		context.Reply(&tele.Audio{
			File:      tele.FromReader(buf),
			Duration:  Duration,
			Caption:   Caption,
			Title:     Title,
			Performer: Performer,
			MIME:      "audio/mp3",
			FileName:  FileName[:len(FileName)-len(filepath.Ext(FileName))] + ".mp3",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	if MediaType == "Voice" {
		context.Reply(&tele.Voice{
			File:     tele.FromReader(buf),
			Duration: Duration,
			Caption:  Caption,
			MIME:     "audio/ogg",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	if MediaType == "Video" {
		context.Reply(&tele.Video{
			File:      tele.FromReader(buf),
			Width:     Width,
			Height:    Height,
			Duration:  Duration,
			Caption:   Caption,
			Thumbnail: Thumbnail,
			Streaming: Streaming,
			MIME:      "video/mp4",
			FileName:  FileName[:len(FileName)-len(filepath.Ext(FileName))] + ".mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	if MediaType == "Animation" {
		context.Reply(&tele.Animation{
			File:      tele.FromReader(buf),
			Width:     Width,
			Height:    Height,
			Duration:  Duration,
			Caption:   Caption,
			Thumbnail: Thumbnail,
			MIME:      "video/mp4",
			FileName:  FileName[:len(FileName)-len(filepath.Ext(FileName))] + ".mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	return nil
}
