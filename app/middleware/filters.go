package middleware

import (
	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func ChatOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Chat == context.Chat().Username {
			return next(context)
		}
		return nil
	}
}

func ChannelOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Channel == context.Chat().Username {
			return next(context)
		}
		return nil
	}
}

func CommentChatOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.CommentChat == context.Chat().ID {
			return next(context)
		}
		return nil
	}
}
