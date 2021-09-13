package commands

import (
	"io"
	"os"
	"path/filepath"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gopkg.in/tucnak/telebot.v3"
)

//Convert given file
func Convert(context telebot.Context) error {
	if context.Message().ReplyTo.Video == nil && context.Message().ReplyTo.Audio == nil && context.Message().ReplyTo.Voice == nil {
		return context.Reply("Пример использования: /convert в ответ на какое-либо сообщение с медиа.")
	}
	var InputFile *os.File
	var OutputFile *os.File
	var ReadBuffer io.ReadCloser
	var err error
	var KwArgs map[string]interface{}
	var Caption string
	var Duration int
	var FileName string
	var Title string
	var Performer string
	var Extension string
	switch {
	case context.Message().ReplyTo.Audio != nil:
		context.Notify(telebot.RecordingAudio)
		Extension = "mp3"
		InputFile, err = os.CreateTemp("", context.Message().ReplyTo.Audio.FileName)
		if err != nil {
			return err
		}
		ReadBuffer, err = utils.Bot.File(&telebot.File{FileID: context.Message().ReplyTo.Audio.FileID})
		if err != nil {
			return err
		}
		InputFile.ReadFrom(ReadBuffer)
		OutputFile, err = os.CreateTemp("", "Audio_*.mp3")
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame"}
		if len(context.Args()) == 1 && context.Args()[0] == "ogg" {
			Extension = "ogg"
			KwArgs = ffmpeg.KwArgs{"c:a": "libopus"}
			OutputFile, err = os.CreateTemp("", "Voice_*.ogg")
			if err != nil {
				return err
			}
		}
		Caption = context.Message().ReplyTo.Audio.Caption
		Duration = context.Message().ReplyTo.Audio.Duration
		FileName = context.Message().ReplyTo.Audio.FileName
		Title = context.Message().ReplyTo.Audio.Title
		Performer = context.Message().ReplyTo.Audio.Performer
	case context.Message().ReplyTo.Video != nil:
		context.Notify(telebot.RecordingVideo)
		Extension = "mp4"
		InputFile, err = os.CreateTemp("", context.Message().ReplyTo.Video.FileName)
		if err != nil {
			return err
		}
		ReadBuffer, err = utils.Bot.File(&telebot.File{FileID: context.Message().ReplyTo.Video.FileID})
		if err != nil {
			return err
		}
		InputFile.ReadFrom(ReadBuffer)
		OutputFile, err = os.CreateTemp("", "Video_*.mp4")
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:v": "libx265", "preset": "fast", "crf": 26}
		if len(context.Args()) == 1 && context.Args()[0] == "gif" {
			KwArgs = ffmpeg.KwArgs{"c:v": "libx265", "an": "", "preset": "fast", "crf": 26}
		}
		if len(context.Args()) == 1 && context.Args()[0] == "mp3" {
			context.Notify(telebot.RecordingAudio)
			Extension = "mp3"
			KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame", "vn": ""}
			OutputFile, err = os.CreateTemp("", "Audio_*.mp3")
			if err != nil {
				return err
			}
		}
		Caption = context.Message().ReplyTo.Video.Caption
		Duration = context.Message().ReplyTo.Video.Duration
		FileName = context.Message().ReplyTo.Video.FileName
	case context.Message().ReplyTo.Voice != nil:
		context.Notify(telebot.RecordingAudio)
		Extension = "ogg"
		InputFile, err = os.CreateTemp("", context.Message().ReplyTo.Voice.FileName)
		if err != nil {
			return err
		}
		ReadBuffer, err = utils.Bot.File(&telebot.File{FileID: context.Message().ReplyTo.Voice.FileID})
		if err != nil {
			return err
		}
		InputFile.ReadFrom(ReadBuffer)
		OutputFile, err = os.CreateTemp("", "Voice_*.ogg")
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.KwArgs{"c:a": "libopus"}
		if len(context.Args()) == 1 && context.Args()[0] == "mp3" {
			Extension = "mp3"
			KwArgs = ffmpeg.KwArgs{"c:a": "libmp3lame"}
			OutputFile, err = os.CreateTemp("", "Audio_*.mp3")
			if err != nil {
				return err
			}
		}
		Caption = context.Message().ReplyTo.Voice.Caption
		Duration = context.Message().ReplyTo.Voice.Duration
		FileName = context.Message().ReplyTo.Voice.FileName
	default:
		return context.Reply("Пример использования: /convert в ответ на какое-либо сообщение с аудио или видео.")
	}
	err = ffmpeg.Input(InputFile.Name()).Output(OutputFile.Name(), KwArgs).OverWriteOutput().WithOutput(nil, os.Stdout).Run()
	if err != nil {
		return err
	}
	FileName = FileName[:len(FileName)-len(filepath.Ext(FileName))]
	if Extension == "mp3" {
		return context.Reply(&telebot.Audio{
			File:      telebot.FromDisk(OutputFile.Name()),
			Duration:  Duration,
			Caption:   Caption,
			Title:     Title,
			Performer: Performer,
			MIME:      "audio/mp3",
			FileName:  FileName + ".mp3",
		})
	}
	if Extension == "ogg" {
		return context.Reply(&telebot.Voice{
			File:     telebot.FromDisk(OutputFile.Name()),
			Duration: Duration,
			Caption:  Caption,
			MIME:     "audio/ogg",
		})
	}
	if Extension == "mp4" {
		return context.Reply(&telebot.Video{
			File:              telebot.FromDisk(OutputFile.Name()),
			Width:             context.Message().ReplyTo.Video.Width,
			Height:            context.Message().ReplyTo.Video.Height,
			Duration:          Duration,
			Caption:           Caption,
			Thumbnail:         context.Message().ReplyTo.Video.Thumbnail,
			SupportsStreaming: context.Message().ReplyTo.Video.SupportsStreaming,
			MIME:              "video/mp4",
			FileName:          FileName + ".mp4",
		})
	}
	return nil
}
