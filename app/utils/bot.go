package utils

import (
	"log"
	"time"

	"github.com/NexonSU/telebot"
)

func BotInit() telebot.Bot {
	if Config.Token == "" {
		log.Fatal("Telegram Bot token not found in config.json")
	}
	if Config.Chat == "" {
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
	settings.Poller = telebot.NewMiddlewarePoller(settings.Poller, func(upd *telebot.Update) bool {
		if upd.Message != nil && upd.Message.Sender != nil {
			GatherData(upd.Message.Sender)
		}

		return true
	})
	var Bot, err = telebot.NewBot(settings)
	if err != nil {
		log.Println(Config.BotApiUrl)
		log.Fatal(err)
	}
	return *Bot
}

var Bot = BotInit()
