package utils

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/tucnak/telebot.v3"
)

func BotInit() telebot.Bot {
	if Config.Telegram.Token == "" {
		log.Fatal("Telegram Bot token not found in config.json")
	}
	if Config.Telegram.Chat == "" {
		log.Fatal("Chat username not found in config.json")
	}
	settings := telebot.Settings{
		URL:       Config.Telegram.BotApiUrl,
		Token:     Config.Telegram.Token,
		ParseMode: telebot.ModeHTML,
		Poller: &telebot.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		},
	}
	if Config.Webhook.EndpointPublicURL != "" || Config.Webhook.Listen != "" {
		settings.Poller = &telebot.Webhook{
			Listen: Config.Webhook.Listen,
			Endpoint: &telebot.WebhookEndpoint{
				PublicURL: Config.Webhook.EndpointPublicURL,
			},
			MaxConnections: Config.Webhook.MaxConnections,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		}
	} else {
		settings.Poller = &telebot.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		}
	}
	settings.Poller = telebot.NewMiddlewarePoller(settings.Poller, func(upd *telebot.Update) bool {
		if upd.Message != nil && upd.Message.Sender != nil {
			GatherData(upd.Message.Sender)
		}

		if upd.ChatMember != nil {
			MarshalledMessage, _ := json.MarshalIndent(upd.ChatMember, "", "    ")
			log.Println(string(MarshalledMessage))
		}

		return true
	})
	var Bot, err = telebot.NewBot(settings)
	if err != nil {
		log.Println(Config.Telegram.BotApiUrl)
		log.Fatal(err)
	}
	return *Bot
}

var Bot = BotInit()
