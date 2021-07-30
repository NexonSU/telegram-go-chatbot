package utils

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func BotInit() tb.Bot {
	if Config.Telegram.Token == "" {
		log.Fatal("Telegram Bot token not found in config.json")
	}
	if Config.Telegram.Chat == "" {
		log.Fatal("Chat username not found in config.json")
	}
	settings := tb.Settings{
		URL:       Config.Telegram.BotApiUrl,
		Token:     Config.Telegram.Token,
		ParseMode: tb.ModeHTML,
		Poller: &tb.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		},
	}
	if Config.Webhook.EndpointPublicURL != "" || Config.Webhook.Listen != "" {
		settings.Poller = &tb.Webhook{
			Listen: Config.Webhook.Listen,
			Endpoint: &tb.WebhookEndpoint{
				PublicURL: Config.Webhook.EndpointPublicURL,
			},
			MaxConnections: Config.Webhook.MaxConnections,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		}
	} else {
		settings.Poller = &tb.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		}
	}
	settings.Poller = tb.NewMiddlewarePoller(settings.Poller, func(upd *tb.Update) bool {
		if upd.Message != nil && upd.Message.Sender != nil {
			GatherData(upd.Message.Sender)
		}

		return true
	})
	var Bot, err = tb.NewBot(settings)
	if err != nil {
		log.Println(Config.Telegram.BotApiUrl)
		log.Fatal(err)
	}
	return *Bot
}

var Bot = BotInit()
