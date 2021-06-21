package utils

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strconv"
	"strings"
	"time"
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
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		}
	} else {
		settings.Poller = &tb.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.Telegram.AllowedUpdates,
		}
	}
	var Bot, err = tb.NewBot(settings)
	if err != nil {
		log.Println(Config.Telegram.BotApiUrl)
		log.Fatal(err)
	}
	return *Bot
}

var Bot = BotInit()

func FindUserInMessage(m tb.Message) (tb.User, int64, error) {
	var user tb.User
	var err error = nil
	var untildate = time.Now().Unix()
	var text = strings.Split(m.Text, " ")
	if m.ReplyTo != nil {
		user = *m.ReplyTo.Sender
		if len(text) == 2 {
			addtime, err := strconv.ParseInt(text[1], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	} else {
		if len(text) == 1 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		user, err = GetUserFromDB(text[1])
		if err != nil {
			return user, untildate, err
		}
		if len(text) == 3 {
			addtime, err := strconv.ParseInt(text[2], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	}
	return user, untildate, err
}
