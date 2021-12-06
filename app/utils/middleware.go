package utils

import (
	"fmt"
	"strings"

	"gopkg.in/tucnak/telebot.v3"
)

func ChatOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.Chat == context.Chat().ID {
			return next(context)
		}
		return nil
	}
}

func ChannelOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.Channel == context.Chat().ID {
			return next(context)
		}
		return nil
	}
}

func CommentChatOnly(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.CommentChat == context.Chat().ID {
			return next(context)
		}
		return nil
	}
}

func GetFilterCreator(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range Config.Admins {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		for _, b := range Config.Moders {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		if Config.Chat == context.Chat().ID {
			var get Get
			if context.Message().ReplyTo == nil && len(context.Args()) == 0 {
				return next(context)
			}
			result := DB.Where(&Get{Name: strings.ToLower(context.Args()[0])}).First(&get)
			if result.RowsAffected != 0 {
				if get.Creator == context.Sender().ID {
					return next(context)
				}
				creator, err := GetUserFromDB(fmt.Sprintf("%v", get.Creator))
				if err != nil {
					return err
				}
				return context.Reply(fmt.Sprintf("Данный гет могут изменять либо администраторы, либо %v.", UserFullName(&creator)))
			}
			return next(context)
		}
		return nil
	}
}

func SysLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range Config.Admins {
			if b == context.Sender().ID {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		for _, b := range Config.Moders {
			if b == context.Sender().ID {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		if Config.Chat == context.Chat().ID {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func AdminLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range Config.Admins {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		for _, b := range Config.Moders {
			if b == context.Sender().ID {
				return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			}
		}
		if Config.Chat == context.Chat().ID {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func ModerLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range Config.Admins {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		for _, b := range Config.Moders {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		if Config.Chat == context.Chat().ID {
			return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		}
		return nil
	}
}

func ChatLevel(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(context telebot.Context) error {
		if Config.SysAdmin == context.Sender().ID {
			return next(context)
		}
		for _, b := range Config.Admins {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		for _, b := range Config.Moders {
			if b == context.Sender().ID {
				return next(context)
			}
		}
		if Config.Chat == context.Chat().ID {
			return next(context)
		}
		return nil
	}
}

type MiddlewarePoller struct {
	Capacity int // Default: 1
	Poller   telebot.Poller
	Filter   func(*telebot.Update) bool
}

func NewMiddlewarePoller(original telebot.Poller, filter func(*telebot.Update) bool) *MiddlewarePoller {
	return &MiddlewarePoller{
		Poller: original,
		Filter: filter,
	}
}

func (p *MiddlewarePoller) Poll(b *telebot.Bot, dest chan telebot.Update, stop chan struct{}) {
	if p.Capacity < 1 {
		p.Capacity = 1
	}

	middle := make(chan telebot.Update, p.Capacity)
	stopPoller := make(chan struct{})

	go p.Poller.Poll(b, middle, stopPoller)

	for {
		select {
		case <-stop:
			close(stopPoller)
			return
		case upd := <-middle:
			if p.Filter(&upd) {
				dest <- upd
			}
		}
	}
}
