package commands

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/gocolly/colly"
	"github.com/kkdai/youtube/v2"
	tele "gopkg.in/telebot.v3"
)

// Convert given file
func Download(context tele.Context) error {
	if (context.Message().ReplyTo == nil && len(context.Args()) == 0) || (context.Message().ReplyTo != nil && len(context.Args()) > 1) {
		return context.Reply("Пример использования: <code>/download {ссылка на ютуб/твиттер}</code>\nИли отправь в ответ на какое-либо сообщение с ссылкой <code>/download</code>")
	}

	videoID := ""
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
			if strings.Contains(text, "youtu.be/") {
				videoID = strings.Split(strings.Split(text, "youtu.be/")[1], "?")[0]
				service = "youtube"
			}
			if strings.Contains(text, "watch?v=") {
				videoID = strings.Split(strings.Split(text, "watch?v=")[1], "&")[0]
				service = "youtube"
			}
			if strings.Contains(text, "twitter.com/") {
				videoID = text
				service = "twitter"
			}
		}
	}

	context.Notify(tele.RecordingVideo)

	if service == "youtube" {
		client := youtube.Client{}

		video, err := client.GetVideo(videoID)
		if err != nil {
			return err
		}

		format := &youtube.Format{}
		formats := video.Formats.WithAudioChannels()

		for i, _ := range formats {
			if (formats[i].ContentLength < 50000000 && formats[i].QualityLabel != "" && formats[i].ContentLength != 0) ||
				(formats[0].QualityLabel == "720p" && video.Duration < 300000000000) {
				format = &formats[i]
				break
			}
		}

		if format.ItagNo == 0 {
			return context.Reply("Видео слишком большое для скачивания.")
		}

		stream, _, err := client.GetStream(video, format)
		if err != nil {
			return err
		}

		return context.Reply(&tele.Video{
			File:      tele.FromReader(stream),
			MIME:      "video/mp4",
			Height:    format.Height,
			Width:     format.Width,
			Streaming: true,
			Thumbnail: &tele.Photo{
				Width:  int(video.Thumbnails[len(video.Thumbnails)-1].Width),
				Height: int(video.Thumbnails[len(video.Thumbnails)-1].Height),
				File:   tele.FromURL(video.Thumbnails[len(video.Thumbnails)-1].URL),
			},
			FileName: videoID + ".mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
	if service == "twitter" {
		downloader := NewTwitterVideoDownloader(videoID)
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
