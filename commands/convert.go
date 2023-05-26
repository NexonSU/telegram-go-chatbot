package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Convert(context tele.Context) error {
	var InputFilePath string
	var OutputFilePath string
	var err error
	var KwArgs map[string]interface{}
	var FileName string
	var Title string
	var Performer string
	var Extension string
	var Width int
	var Height int
	var Caption string
	var Duration int
	var Thumbnail *tele.Photo
	var Streaming bool
	var MediaType string
	TempName := time.Now().UnixNano()
	utils.Bot.URL = "https://api.telegram.org"
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
		Extension = "mp3"
		MediaType = "Audio"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v%v", os.TempDir(), TempName, filepath.Ext(context.Message().ReplyTo.Audio.FileName))
		err = utils.Bot.Download(&context.Message().ReplyTo.Audio.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame", "timelimit": 60}
	case context.Message().ReplyTo.Document != nil && context.Message().ReplyTo.Document.MIME[0:5] == "video":
		context.Notify(tele.RecordingVideo)
		Caption = context.Message().ReplyTo.Document.Caption
		FileName = context.Message().ReplyTo.Document.FileName
		Extension = "mp4"
		MediaType = "Video"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v%v", os.TempDir(), TempName, filepath.Ext(context.Message().ReplyTo.Document.FileName))
		err = utils.Bot.Download(&context.Message().ReplyTo.Document.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:v": "libx264", "preset": "fast", "crf": 26, "timelimit": 900, "movflags": "+faststart", "c:a": "aac"}
	case context.Message().ReplyTo.Video != nil:
		context.Notify(tele.RecordingVideo)
		Width = context.Message().ReplyTo.Video.Width
		Height = context.Message().ReplyTo.Video.Height
		Caption = context.Message().ReplyTo.Video.Caption
		Duration = context.Message().ReplyTo.Video.Duration
		Thumbnail = context.Message().ReplyTo.Video.Thumbnail
		Streaming = context.Message().ReplyTo.Video.Streaming
		FileName = context.Message().ReplyTo.Video.FileName
		Extension = "mp4"
		MediaType = "Video"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v%v", os.TempDir(), TempName, filepath.Ext(context.Message().ReplyTo.Video.FileName))
		err = utils.Bot.Download(&context.Message().ReplyTo.Video.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:v": "libx264", "preset": "fast", "crf": 26, "timelimit": 900, "movflags": "+faststart", "c:a": "aac"}
	case context.Message().ReplyTo.Voice != nil:
		context.Notify(tele.RecordingAudio)
		Caption = context.Message().ReplyTo.Voice.Caption
		Duration = context.Message().ReplyTo.Voice.Duration
		Extension = "ogg"
		MediaType = "Voice"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v", os.TempDir(), TempName)
		err = utils.Bot.Download(&context.Message().ReplyTo.Voice.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:a": "libopus", "timelimit": 60}
	default:
		return context.Reply("Пример использования: /convert в ответ на какое-либо сообщение с аудио или видео.", &tele.SendOptions{AllowWithoutReply: true})
	}
	OutputFilePath = fmt.Sprintf("%v/convert_output_%v.%v", os.TempDir(), TempName, Extension)
	if RequestedMediaType == "Animation" {
		MediaType = "Animation"
		KwArgs = ffmpeg.KwArgs{"c:v": "libx264", "an": "", "preset": "fast", "crf": 26, "timelimit": 900, "movflags": "+faststart"}
	}
	if RequestedMediaType == "Audio" {
		Extension = "mp3"
		MediaType = "Audio"
		KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame", "timelimit": 60}
	}
	if RequestedMediaType == "Voice" {
		Extension = "ogg"
		MediaType = "Voice"
		KwArgs = ffmpeg.KwArgs{"c:a": "libopus", "timelimit": 60}
	}
	err = ffmpeg.Input(InputFilePath).Output(OutputFilePath, KwArgs).OverWriteOutput().WithOutput(nil, os.Stdout).Run()
	if err != nil {
		return err
	}
	os.Remove(InputFilePath)
	if MediaType == "Audio" {
		context.Reply(&tele.Audio{
			File:      tele.FromDisk(OutputFilePath),
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
			File:     tele.FromDisk(OutputFilePath),
			Duration: Duration,
			Caption:  Caption,
			MIME:     "audio/ogg",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	if MediaType == "Video" {
		context.Reply(&tele.Video{
			File:      tele.FromDisk(OutputFilePath),
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
			File:      tele.FromDisk(OutputFilePath),
			Width:     Width,
			Height:    Height,
			Duration:  Duration,
			Caption:   Caption,
			Thumbnail: Thumbnail,
			MIME:      "video/mp4",
			FileName:  FileName[:len(FileName)-len(filepath.Ext(FileName))] + ".mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	os.Remove(OutputFilePath)
	return nil
}
