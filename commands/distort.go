package commands

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	_ "golang.org/x/image/bmp"

	"github.com/Jeffail/tunny"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
)

var DistortBusy bool

// Distort given file
func Distort(context tele.Context) error {
	if context.Message().ReplyTo == nil {
		return context.Reply("Пример использования: <code>/distort</code> в ответ на какое-либо сообщение с видео.")
	}
	if context.Message().ReplyTo.Media() == nil {
		return context.Reply("Какого-либо видео нет в указанном сообщении.")
	}

	media := context.Message().ReplyTo.Media()
	extension := ""
	outputKwArgs := ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": ""}

	switch media.MediaType() {
	case "video":
		extension = "mp4"
		if context.Message().ReplyTo.Video.Duration > 60 {
			return context.Reply("Слишком длинное видео. Лимит 60 секунд.")
		}
	case "animation":
		extension = "mp4"
		outputKwArgs = ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": "", "an": ""}
		if context.Message().ReplyTo.Animation.Duration > 60 {
			return context.Reply("Слишком длинная гифка. Лимит 60 секунд.")
		}
	default:
		return context.Reply("Неподдерживаемая операция")
	}

	if DistortBusy {
		return context.Reply("Команда занята")
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
		DistortBusy = false
		done <- true
	}()
	DistortBusy = true

	jobStarted := time.Now().Unix()

	file, err := utils.Bot.FileByID(media.MediaFile().FileID)
	if err != nil {
		return err
	}

	workdir := fmt.Sprintf("%v/telegram-go-chatbot-distort/%v", os.TempDir(), media.MediaFile().FileID)
	inputFile := file.FilePath
	outputFile := fmt.Sprintf("%v/output.%v", workdir, extension)

	if err := os.Mkdir(workdir, os.ModePerm); err != nil {
		return err
	}
	defer func(workdir string) {
		os.RemoveAll(workdir)
	}(workdir)

	err = ffmpeg.Input(inputFile).Output(workdir + "/%09d.bmp").OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	files, err := filepath.Glob(workdir + "/*.bmp")
	if err != nil {
		return err
	}

	f, err := os.Open(files[0])
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
	scale := 0

	pool := tunny.NewFunc(runtime.NumCPU(), func(payload interface{}) interface{} {
		payloadCommand := strings.Fields(payload.(string))
		return exec.Command(payloadCommand[0], payloadCommand[1:]...).Run()
	})
	defer pool.Close()

	for i, file := range files {
		scale = 512 - (i * 340 / len(files))
		command := fmt.Sprintf("convert %v -liquid-rescale %vx%v -resize %vx%v %v", file, scale, scale, width, height, file)
		go func(command string) {
			if pool.Process(command) != nil {
				err = pool.Process(command).(error)
			}
		}(command)
	}

	for {
		time.Sleep(1 * time.Second)
		if time.Now().Unix()-jobStarted > 300 {
			return fmt.Errorf("слишком долгое выполнение операции")
		}
		if pool.QueueLength() == 0 {
			break
		}
	}
	if err != nil {
		return err
	}

	err = ffmpeg.Input(workdir+"/%09d.bmp").Output(outputFile, outputKwArgs).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	DistortBusy = false
	switch media.MediaType() {
	case "video":
		return context.Reply(&tele.Video{
			File:      tele.FromDisk(outputFile),
			FileName:  media.MediaFile().FileID + "." + extension,
			Streaming: true,
			Width:     width,
			Height:    height,
			MIME:      "video/mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	case "animation":
		return context.Reply(&tele.Animation{
			File:     tele.FromDisk(outputFile),
			FileName: media.MediaFile().FileID + "." + extension,
			Width:    width,
			Height:   height,
			MIME:     "video/mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	default:
		return context.Reply(&tele.Document{
			File:     tele.FromDisk(outputFile),
			FileName: media.MediaFile().FileID + "." + extension,
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
}

func init() {
	dir := fmt.Sprintf("%v/telegram-go-chatbot-distort", os.TempDir())
	if err := os.Mkdir(dir, os.ModePerm); err != nil {
		os.RemoveAll(dir)
		os.Mkdir(dir, os.ModePerm)
	}
	DistortBusy = false
}
