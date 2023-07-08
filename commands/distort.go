package commands

import (
	"bufio"
	"fmt"
	"image"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	_ "image/png"

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
	additionalInputArgs := ""

	switch media.MediaType() {
	case "video", "animation", "photo", "audio", "voice", "sticker":
		break
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

	framerate := "30/1"

	if media.MediaType() != "audio" && media.MediaType() != "voice" {
		frames := data.FirstVideoStream().NbFrames
		framerate = data.FirstVideoStream().AvgFrameRate

		if frames == "" {
			frames = "1"
		}

		framesInt, err := strconv.Atoi(frames)
		if err != nil {
			return err
		}

		if framesInt > 1000 {
			return context.Reply("Видео слишком длинное. Максимум 1000 фреймов.")
		}
	}

	if err := os.Mkdir(workdir, os.ModePerm); err != nil {
		return context.Reply("Обработка файла уже выполняется")
	}
	defer func(workdir string) {
		os.RemoveAll(workdir)
	}(workdir)

	if media.MediaType() == "video" {
		ffmpeg.Input(inputFile).Output(workdir + "/input_audio.mp3").OverWriteOutput().ErrorToStdOut().Run()
		ffmpeg.Input(workdir+"/input_audio.mp3").Output(workdir+"/audio.mp3", ffmpeg.KwArgs{"filter_complex": "vibrato=f=10:d=0.7"}).OverWriteOutput().ErrorToStdOut().Run()
		additionalInputArgs = "-i " + workdir + "/audio.mp3 -c:a aac"
	}

	if media.MediaType() == "audio" || media.MediaType() == "voice" {
		ffmpeg.Input(inputFile).Output(workdir + "/input_audio.mp3").OverWriteOutput().ErrorToStdOut().Run()
		err = ffmpeg.Input(workdir+"/input_audio.mp3").Output(workdir+"/audio.mp3", ffmpeg.KwArgs{"filter_complex": "vibrato=f=10:d=0.7"}).OverWriteOutput().ErrorToStdOut().Run()
		if err != nil {
			return err
		}
		return context.Reply(&tele.Audio{
			File:     tele.FromDisk(workdir + "/audio.mp3"),
			FileName: media.MediaFile().FileID + ".mp3",
			MIME:     "video/mp3",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}

	err = ffmpeg.Input(inputFile).Output(workdir + "/%09d.png").OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	if media.MediaType() == "photo" || (media.MediaType() == "sticker" && !context.Message().ReplyTo.Sticker.Animated && !context.Message().ReplyTo.Sticker.Video) {
		framerate = "15/1"
		src := workdir + "/000000001.png"
		for i := 2; i < 31; i++ {
			dst := fmt.Sprintf("%v/%09d.png", workdir, i)

			sourceFileStat, err := os.Stat(src)
			if err != nil {
				return err
			}

			if !sourceFileStat.Mode().IsRegular() {
				return fmt.Errorf("%s is not a regular file", src)
			}

			source, err := os.Open(src)
			if err != nil {
				return err
			}
			defer source.Close()

			destination, err := os.Create(dst)
			if err != nil {
				return err
			}
			defer destination.Close()
			_, err = io.Copy(destination, source)
			if err != nil {
				return err
			}
		}
	}

	files, err := filepath.Glob(workdir + "/*.png")
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

	if width%2 != 0 {
		width++
	}

	if height%2 != 0 {
		height++
	}

	pool := tunny.NewFunc(runtime.NumCPU()-1, func(payload interface{}) interface{} {
		payloadCommand := strings.Fields(payload.(string))
		return exec.Command(payloadCommand[0], payloadCommand[1:]...).Run()
	})
	defer pool.Close()

	for i, file := range files {
		scale = 90 - (i * 65 / len(files))
		command := fmt.Sprintf("convert %v -liquid-rescale %v%% -resize %vx%v! %v", file, scale, width, height, file)
		go func(command string) {
			if pool.Process(command) != nil {
				err = pool.Process(command).(error)
			}
		}(command)
	}

	for {
		time.Sleep(1 * time.Second)
		if time.Now().Unix()-jobStarted > 300 {
			return context.Reply("Слишком долгое выполнение операции")
		}
		if pool.QueueLength() == 0 {
			break
		}
	}
	if err != nil {
		return err
	}

	ffmpegCommand := fmt.Sprintf("ffmpeg -y -framerate %v -i %v/%%09d.png %v -c:v: libx264 -preset fast -crf 26 -pix_fmt yuv420p -movflags +faststart -hide_banner -loglevel fatal %v", framerate, workdir, additionalInputArgs, outputFile)
	ffmpegCommandExec := strings.Fields(ffmpegCommand)
	err = exec.Command(ffmpegCommandExec[0], ffmpegCommandExec[1:]...).Run()
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
