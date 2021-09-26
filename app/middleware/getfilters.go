package middleware

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func GetFilterCreator(next telebot.HandlerFunc) telebot.HandlerFunc {
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
			var get utils.Get
			if context.Message().ReplyTo == nil && len(context.Args()) == 0 {
				return next(context)
			}
			result := utils.DB.Where(&utils.Get{Name: strings.ToLower(context.Args()[0])}).First(&get)
			if result.RowsAffected != 0 {
				if get.Creator == context.Sender().ID {
					return next(context)
				}
				creator, err := utils.GetUserFromDB(fmt.Sprintf("%v", get.Creator))
				if err != nil {
					return err
				}
				return context.Reply(fmt.Sprintf("Данный гет могут изменять либо администраторы, либо %v.", utils.UserFullName(&creator)))
			}
			return next(context)
		}
		return nil
	}
}
