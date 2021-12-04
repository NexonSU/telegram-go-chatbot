package utils

import (
	"fmt"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v3"
)

func BotInit() telebot.Bot {
	if Config.Token == "" {
		log.Fatal("Telegram Bot token not found in config.json")
	}
	if Config.Chat == 0 {
		log.Fatal("Chat username not found in config.json")
	}
	settings := telebot.Settings{
		URL:       Config.BotApiUrl,
		Token:     Config.Token,
		ParseMode: telebot.ModeHTML,
		Poller: &telebot.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.AllowedUpdates,
		},
	}
	if Config.EndpointPublicURL != "" || Config.Listen != "" {
		settings.Poller = &telebot.Webhook{
			Listen: Config.Listen,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: Config.EndpointPublicURL,
			},
			MaxConnections: Config.MaxConnections,
			AllowedUpdates: Config.AllowedUpdates,
		}
	} else {
		settings.Poller = &telebot.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.AllowedUpdates,
		}
	}
	var Bot, err = telebot.NewBot(settings)
	if err != nil {
		log.Println(Config.BotApiUrl)
		log.Fatal(err)
	}
	if Config.SysAdmin != 0 {
		Bot.Send(telebot.ChatID(Config.SysAdmin), fmt.Sprintf("%v has finished starting up.", MentionUser(Bot.Me)))
	}
	return *Bot
}

var Bot = BotInit()
