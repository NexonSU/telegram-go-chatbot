package middleware

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func ChatOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.Chat == context.Chat().Username {
			return next(context)
		}
		return nil
	}
}

func ChannelOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.Channel == context.Chat().Username {
			return next(context)
		}
		return nil
	}
}

func CommentChatOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.CommentChat == context.Chat().ID {
			return next(context)
		}
		return nil
	}
}
