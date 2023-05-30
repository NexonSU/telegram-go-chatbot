package commands

import (
	cntx "context"
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/wader/goutubedl"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Download(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) > 1) {
		return context.Reply("Пример использования: <code>/download {ссылка на ютуб/твиттер}</code>\nИли отправь в ответ на какое-либо сообщение с ссылкой <code>/download</code>")
	}

	link := ""
	message := &tele.Message{}
	service := ""

	if context.Message().ReplyTo == nil {
		message = context.Message()
	} else {
		message = context.Message().ReplyTo
	}

	for _, entity := range message.Entities {
		if entity.Type == tele.EntityURL {
			text := message.EntityText(entity)
			link = text
			if strings.Contains(text, "youtu") {
				service = "youtube"
			}
			if strings.Contains(text, "twitter.com/") {
				service = "twitter"
			}
		}
	}

	context.Notify(tele.RecordingVideo)

	if service == "youtube" {
		goutubedl.Path = "yt-dlp"

		result, err := goutubedl.New(cntx.Background(), link, goutubedl.Options{})
		if err != nil {
			return err
		}

		if result.Info.Filesize > 50000000 || result.Info.FilesizeApprox > 50000000 {
			if !strings.Contains(link, "/clip/") {
				return context.Reply("Видео больше 50МБ")
			}
		}

		downloadResult, err := result.Download(cntx.Background(), "best")
		if err != nil {
			return err
		}
		defer downloadResult.Close()

		filename := strings.Split(link, "/")[len(strings.Split(link, "/"))-1]
		filename = strings.ReplaceAll(filename, "watch?v=", "")

		return context.Reply(&tele.Video{
			File:      tele.FromReader(downloadResult),
			MIME:      "video/mp4",
			Height:    int(result.Info.Height),
			Width:     int(result.Info.Width),
			Streaming: true,
			FileName:  filename + ".mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	if service == "twitter" {
		downloader := NewTwitterVideoDownloader(link)
		fileName := downloader.Download()
		context.Reply(&tele.Video{
			File:      tele.FromDisk(fileName),
			MIME:      "video/mp4",
			Streaming: true,
			FileName:  fileName,
		}, &tele.SendOptions{AllowWithoutReply: true})
		return os.Remove(fileName)
	}
	return context.Reply("Ссылка не найдена или сервис не поддерживается.")
}

type TwitterVideoDownloader struct {
	video_url    string
	bearer_token string
	xguest_token string
}

func NewTwitterVideoDownloader(url string) *TwitterVideoDownloader {
	self := new(TwitterVideoDownloader)
	self.video_url = url
	return self
}

func (self *TwitterVideoDownloader) GetBearerToken() string {
	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		pattern, _ := regexp.Compile(`"Bearer.*?"`)
		self.bearer_token = strings.Trim(pattern.FindString(string(r.Body)), `"`)
	})

	c.Visit("https://abs.twimg.com/web-video-player/TwitterVideoPlayerIframe.cefd459559024bfb.js")

	return self.bearer_token
}

func (self *TwitterVideoDownloader) GetXGuestToken() string {
	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Authorization", self.bearer_token)
	})

	c.OnResponse(func(r *colly.Response) {
		pattern, _ := regexp.Compile(`[0-9]+`)
		self.xguest_token = pattern.FindString(string(r.Body))
	})

	c.Post("https://api.twitter.com/1.1/guest/activate.json", nil)

	return self.xguest_token
}

func (self *TwitterVideoDownloader) GetM3U8Urls() string {
	var m3u8_urls string

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Authorization", self.bearer_token)
		r.Headers.Set("x-guest-token", self.xguest_token)
	})

	c.OnResponse(func(r *colly.Response) {
		pattern, _ := regexp.Compile(`https.*m3u8`)
		m3u8_urls = strings.ReplaceAll(pattern.FindString(string(r.Body)), "\\", "")
	})

	url := "https://api.twitter.com/1.1/videos/tweet/config/" +
		strings.Split(self.video_url, "/status/")[1] +
		".json"

	c.Visit(url)

	return m3u8_urls
}

func (self *TwitterVideoDownloader) GetM3U8Url(m3u8_urls string) string {
	var m3u8_url string

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		pattern, _ := regexp.Compile(`.*m3u8`)
		m3u8_urls := pattern.FindAllString(string(r.Body), -1)
		m3u8_url = "https://video.twimg.com" + m3u8_urls[len(m3u8_urls)-1]
	})

	c.Visit(m3u8_urls)

	return m3u8_url
}

func (self *TwitterVideoDownloader) Download() string {
	self.GetBearerToken()
	self.GetXGuestToken()
	m3u8_urls := self.GetM3U8Urls()
	m3u8_url := self.GetM3U8Url(m3u8_urls)

	sum := md5.Sum([]byte(m3u8_url))
	filename := hex.EncodeToString(sum[:]) + ".mp4"

	cmd := exec.Command("ffmpeg", "-y", "-i", m3u8_url, "-c", "copy", filename)
	cmd.Run()

	return filename
}
