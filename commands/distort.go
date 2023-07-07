package commands

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
)

// Invert given file
func Distort(context tele.Context) error {
	if context.Message().ReplyTo == nil {
		return context.Reply("Пример использования: <code>/distort</code> в ответ на какое-либо сообщение с видео.")
	}
	if context.Message().ReplyTo.Media() == nil {
		return context.Reply("Какого-либо видео нет в указанном сообщении.")
	}

	media := context.Message().ReplyTo.Media()
	extension := ""

	switch media.MediaType() {
	case "video":
		extension = "mp4"
	case "animation":
		extension = "mp4"
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
				context.Notify(tele.ChatAction(tele.UploadingDocument))
			}
			time.Sleep(time.Second * 5)
		}
	}()
	defer func() {
		done <- true
	}()

	workdir := fmt.Sprintf("/tmp/%v", media.MediaFile().FileID)
	inputFile := fmt.Sprintf("%v/input.%v", workdir, extension)
	outputFile := fmt.Sprintf("%v/output.%v", workdir, extension)

	if err := os.Mkdir(workdir, os.ModePerm); err != nil {
		return err
	}
	defer func(workdir string) {
		os.RemoveAll(workdir)
	}(workdir)

	err := utils.Bot.Download(media.MediaFile(), inputFile)
	if err != nil {
		return err
	}

	err = ffmpeg.Input(inputFile).Output(workdir + "/%04d.png").OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	files, err := filepath.Glob(workdir + "/*.png")
	if err != nil {
		return err
	}

	f, err := os.Open(workdir + "/0001.png")
	if err != nil {
		return err
	}
	frameConfig, _, err := image.DecodeConfig(bufio.NewReader(f))
	if err != nil {
		return err
	}
	f.Close()
	width := frameConfig.Width
	height := frameConfig.Height
	scaleWidth := 0
	scaleHeight := 0

	for i, v := range files {
		scaleWidth = width - (i * 340 / len(files))
		scaleHeight = height - (i * 340 / len(files))
		err = exec.Command("convert", v, "-liquid-rescale", fmt.Sprintf("%vx%v", scaleWidth, scaleHeight), "-resize", fmt.Sprintf("%vx%v", width, height), v).Run()
		if err != nil {
			return err
		}
	}

	err = ffmpeg.Input(inputFile, ffmpeg.KwArgs{"i": workdir + "/%04d.png"}).Output(outputFile, ffmpeg.KwArgs{"i": workdir + "/%04d.png"}).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	return context.Reply(&tele.Video{
		File:      tele.FromDisk(outputFile),
		FileName:  media.MediaFile().FileID + "." + extension,
		Streaming: true,
		MIME:      "video/mp4",
	}, &tele.SendOptions{AllowWithoutReply: true})
}
