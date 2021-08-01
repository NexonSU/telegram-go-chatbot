package middleware

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func SysLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Telegram.Admins {
			if b == context.Sender().Username {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		for _, b := range utils.Config.Telegram.Moders {
			if b == context.Sender().Username {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		if utils.Config.Telegram.Chat == context.Chat().Username {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func AdminLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Telegram.Admins {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		for _, b := range utils.Config.Telegram.Moders {
			if b == context.Sender().Username {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		if utils.Config.Telegram.Chat == context.Chat().Username {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func ModerLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Telegram.Admins {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		for _, b := range utils.Config.Telegram.Moders {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		if utils.Config.Telegram.Chat == context.Chat().Username {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func ChatLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.Telegram.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Telegram.Admins {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		for _, b := range utils.Config.Telegram.Moders {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		if utils.Config.Telegram.Chat == context.Chat().Username {
			return next(context)
		}
		return nil
	}
}
