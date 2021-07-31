package services

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/valyala/fastjson"
	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

func ZavtraStreamCheckService() {
	for {
		delay := 240
		if time.Now().Hour() < 2 && time.Now().Hour() >= 20 {
			delay = 30
		}
		err := zavtraStreamCheck("youtube")
		if err != nil {
			log.Println(err.Error())
			_, _ = utils.Bot.Send(telebot.ChatID(utils.Config.Telegram.SysAdmin), fmt.Sprintf("ZavtraStreamCheck error:\n<code>%v</code>", err.Error()))
		}
		time.Sleep(time.Duration(delay) * time.Second)
	}
}

func zavtraStreamCheck(service string) error {
	if service == "youtube" {
		if utils.Config.Youtube.ApiKey == "" || utils.Config.Youtube.ChannelName == "" || utils.Config.Youtube.ChannelID == "" || utils.Config.Youtube.StreamChannel == "" {
			return nil
		}
		var stream utils.ZavtraStream
		var httpClient = &http.Client{Timeout: 10 * time.Second}
		r, err := httpClient.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%v&type=video&eventType=live&key=%v", utils.Config.Youtube.ChannelID, utils.Config.Youtube.ApiKey))
		if err != nil {
			return err
		}
		defer func(Body io.ReadCloser) {
			Body.Close()
		}(r.Body)
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		stream.Service = service
		utils.DB.First(&stream)
		results := fastjson.GetInt(jsonBytes, "pageInfo", "totalResults")
		if results != 0 {
			title := fastjson.GetString(jsonBytes, "items", "0", "snippet", "title")
			videoId := fastjson.GetString(jsonBytes, "items", "0", "id", "videoId")
			if stream.VideoID != videoId {
				thumbnail := fmt.Sprintf("https://i.ytimg.com/vi/%v/maxresdefault_live.jpg", videoId)
				caption := fmt.Sprintf("Стрим \"%v\" начался.\nhttps://youtube.com/%v/live", title, utils.Config.Youtube.ChannelName)
				chat, err := utils.Bot.ChatByID("@" + utils.Config.Youtube.StreamChannel)
				if err != nil {
					return err
				}
				_, err = utils.Bot.Send(chat, &telebot.Photo{File: telebot.File{FileURL: thumbnail}, Caption: caption})
				if err != nil {
					return err
				}
				stream.VideoID = videoId
			}
		}
		stream.LastCheck = time.Now()
		result := utils.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(stream)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
	return nil
}
