package middleware

import (
	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func SysLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Admins {
			if b == context.Sender().Username {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		for _, b := range utils.Config.Moders {
			if b == context.Sender().Username {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		if utils.Config.Chat == context.Chat().Username {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func AdminLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Admins {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		for _, b := range utils.Config.Moders {
			if b == context.Sender().Username {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		if utils.Config.Chat == context.Chat().Username {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func ModerLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Admins {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		for _, b := range utils.Config.Moders {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		if utils.Config.Chat == context.Chat().Username {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func ChatLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if utils.Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range utils.Config.Admins {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		for _, b := range utils.Config.Moders {
			if b == context.Sender().Username {
				return next(context)
			}
		}
		if utils.Config.Chat == context.Chat().Username {
			return next(context)
		}
		return nil
	}
}
