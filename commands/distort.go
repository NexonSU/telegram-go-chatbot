package commands

import (
	"bufio"
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "golang.org/x/image/bmp"

	cntx "context"

	"github.com/Jeffail/tunny"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/vansante/go-ffprobe.v2"
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
	outputKwArgs := ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": ""}
	inputKwArgs := ffmpeg.KwArgs{}

	switch media.MediaType() {
	case "video":
		break
	case "animation":
		outputKwArgs = ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": "", "an": ""}
	case "sticker":
		if !context.Message().ReplyTo.Sticker.Animated && !context.Message().ReplyTo.Sticker.Video {
			return context.Reply("Неподдерживаемая операция")
		}
		outputKwArgs = ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": "", "an": ""}
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
	outputFile := fmt.Sprintf("%v/output.mp4", workdir)

	ctx, cancelFn := cntx.WithTimeout(cntx.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, inputFile)
	if err != nil {
		return err
	}

	frames, err := strconv.Atoi(data.FirstVideoStream().NbFrames)
	if err != nil {
		return err
	}

	if frames > 1000 {
		return context.Reply("Видео слишком длинное. Максимум 1000 фреймов.")
	}

	if err := os.Mkdir(workdir, os.ModePerm); err != nil {
		return context.Reply("Обработка файла уже выполняется")
	}
	defer func(workdir string) {
		os.RemoveAll(workdir)
	}(workdir)

	if media.MediaType() == "video" {
		ffmpeg.Input(inputFile).Output(workdir + "/audio.mp3").OverWriteOutput().ErrorToStdOut().Run()
		outputKwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{outputKwArgs, {"i": workdir + "/audio.mp3", "filter_complex": "vibrato=f=8,aphaser=type=t:speed=2:decay=0.6"}})
	}

	inputKwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{inputKwArgs, {"framerate": data.FirstVideoStream().AvgFrameRate}})

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

	pool := tunny.NewFunc(runtime.NumCPU()-1, func(payload interface{}) interface{} {
		payloadCommand := strings.Fields(payload.(string))
		return exec.Command(payloadCommand[0], payloadCommand[1:]...).Run()
	})
	defer pool.Close()

	for i, file := range files {
		scale = 100 - (i * 75 / len(files))
		command := fmt.Sprintf("convert %v -liquid-rescale %v%% -resize %vx%v %v", file, scale, width, height, file)
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

	err = ffmpeg.Input(workdir+"/%09d.bmp", inputKwArgs).Output(outputFile, outputKwArgs).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	DistortBusy = false
	switch media.MediaType() {
	case "video":
		return context.Reply(&tele.Video{
			File:      tele.FromDisk(outputFile),
			FileName:  media.MediaFile().FileID + ".mp4",
			Streaming: true,
			Width:     width,
			Height:    height,
			MIME:      "video/mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	case "animation", "sticker":
		return context.Reply(&tele.Animation{
			File:     tele.FromDisk(outputFile),
			FileName: media.MediaFile().FileID + ".mp4",
			Width:    width,
			Height:   height,
			MIME:     "video/mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	default:
		return context.Reply(&tele.Document{
			File:     tele.FromDisk(outputFile),
			FileName: media.MediaFile().FileID + ".mp4",
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
