package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/tucnak/telebot.v3"
)

//Convert given file
func Convert(context telebot.Context) error {
	if context.Message().ReplyTo.Video == nil && context.Message().ReplyTo.Audio == nil && context.Message().ReplyTo.Voice == nil {
		return context.Reply("Пример использования: /convert в ответ на какое-либо сообщение с медиа.", &telebot.SendOptions{AllowWithoutReply: true})
	}
	var InputFilePath string
	var OutputFilePath string
	var err error
	var KwArgs map[string]interface{}
	var Caption string
	var Duration int
	var FileName string
	var Title string
	var Performer string
	var Extension string
	TempName := time.Now().UnixNano()
	utils.Bot.URL = "https://api.telegram.org"
	switch {
	case context.Message().ReplyTo.Audio != nil:
		context.Notify(telebot.RecordingAudio)
		Caption = context.Message().ReplyTo.Audio.Caption
		Duration = context.Message().ReplyTo.Audio.Duration
		FileName = context.Message().ReplyTo.Audio.FileName
		Title = context.Message().ReplyTo.Audio.Title
		Performer = context.Message().ReplyTo.Audio.Performer
		Extension = "mp3"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v%v", os.TempDir(), TempName, filepath.Ext(context.Message().ReplyTo.Audio.FileName))
		err = utils.Bot.Download(&context.Message().ReplyTo.Audio.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame", "timelimit": 60}
		if len(context.Args()) == 1 && context.Args()[0] == "ogg" {
			Extension = "ogg"
			KwArgs = ffmpeg.KwArgs{"c:a": "libopus", "timelimit": 60}
		}
	case context.Message().ReplyTo.Video != nil:
		context.Notify(telebot.RecordingVideo)
		Caption = context.Message().ReplyTo.Video.Caption
		Duration = context.Message().ReplyTo.Video.Duration
		FileName = context.Message().ReplyTo.Video.FileName
		Extension = "mp4"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v%v", os.TempDir(), TempName, filepath.Ext(context.Message().ReplyTo.Video.FileName))
		err = utils.Bot.Download(&context.Message().ReplyTo.Video.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:v": "libx264", "preset": "fast", "crf": 26, "timelimit": 900, "movflags": "+faststart", "c:a": "aac"}
		if len(context.Args()) == 1 && context.Args()[0] == "gif" {
			KwArgs = ffmpeg.KwArgs{"c:v": "libx264", "an": "", "preset": "fast", "crf": 26, "timelimit": 900, "movflags": "+faststart", "c:a": "aac"}
		}
		if len(context.Args()) == 1 && context.Args()[0] == "mp3" {
			context.Notify(telebot.RecordingAudio)
			Extension = "mp3"
			KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame", "vn": ""}
		}
	case context.Message().ReplyTo.Voice != nil:
		context.Notify(telebot.RecordingAudio)
		Caption = context.Message().ReplyTo.Voice.Caption
		Duration = context.Message().ReplyTo.Voice.Duration
		FileName = context.Message().ReplyTo.Voice.FileName
		Extension = "ogg"
		InputFilePath = fmt.Sprintf("%v/convert_input_%v%v", os.TempDir(), TempName, filepath.Ext(context.Message().ReplyTo.Voice.FileName))
		err = utils.Bot.Download(&context.Message().ReplyTo.Voice.File, InputFilePath)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:a": "libopus", "timelimit": 60}
		if len(context.Args()) == 1 && context.Args()[0] == "mp3" {
			Extension = "mp3"
			KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame", "timelimit": 60}
		}
	default:
		return context.Reply("Пример использования: /convert в ответ на какое-либо сообщение с аудио или видео.", &telebot.SendOptions{AllowWithoutReply: true})
	}
	FileNameWOExt := FileName[:len(FileName)-len(filepath.Ext(FileName))]
	OutputFilePath = fmt.Sprintf("%v/convert_output_%v.%v", os.TempDir(), TempName, Extension)
	err = ffmpeg.Input(InputFilePath).Output(OutputFilePath, KwArgs).OverWriteOutput().WithOutput(nil, os.Stdout).Run()
	if err != nil {
		return err
	}
	os.Remove(InputFilePath)
	if Extension == "mp3" {
		context.Reply(&telebot.Audio{
			File:      telebot.FromDisk(OutputFilePath),
			Duration:  Duration,
			Caption:   Caption,
			Title:     Title,
			Performer: Performer,
			MIME:      "audio/mp3",
			FileName:  FileNameWOExt + ".mp3",
		}, &telebot.SendOptions{AllowWithoutReply: true})
	}
	if Extension == "ogg" {
		context.Reply(&telebot.Voice{
			File:     telebot.FromDisk(OutputFilePath),
			Duration: Duration,
			Caption:  Caption,
			MIME:     "audio/ogg",
		}, &telebot.SendOptions{AllowWithoutReply: true})
	}
	if Extension == "mp4" {
		context.Reply(&telebot.Video{
			File:              telebot.FromDisk(OutputFilePath),
			Width:             context.Message().ReplyTo.Video.Width,
			Height:            context.Message().ReplyTo.Video.Height,
			Duration:          Duration,
			Caption:           Caption,
			Thumbnail:         context.Message().ReplyTo.Video.Thumbnail,
			SupportsStreaming: context.Message().ReplyTo.Video.SupportsStreaming,
			MIME:              "video/mp4",
			FileName:          FileNameWOExt + ".mp4",
		}, &telebot.SendOptions{AllowWithoutReply: true})
	}
	os.Remove(OutputFilePath)
	return nil
}
