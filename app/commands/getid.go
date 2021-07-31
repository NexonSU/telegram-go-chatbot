package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send userid on /getid
func Getid(context telebot.Context) error {
	var err error
	if !utils.IsAdminOrModer(context.Sender().Username) {
		if context.Chat().Username != utils.Config.Telegram.Chat {
			return err
		}
		err := context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			return err
		}
		return err
	}
	var text = strings.Split(context.Text(), " ")
	if context.Message().ReplyTo != nil && context.Message().ReplyTo.OriginalSender != nil {
		_, err := utils.Bot.Send(context.Sender(), fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", context.Message().ReplyTo.OriginalSender.FirstName, context.Message().ReplyTo.OriginalSender.LastName, context.Message().ReplyTo.OriginalSender.Username, context.Message().ReplyTo.OriginalSender.ID))
		if err != nil {
			return err
		}
	} else if context.Message().ReplyTo != nil {
		_, err := utils.Bot.Send(context.Sender(), fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", context.Message().ReplyTo.Sender.FirstName, context.Message().ReplyTo.Sender.LastName, context.Message().ReplyTo.Sender.Username, context.Message().ReplyTo.Sender.ID))
		if err != nil {
			return err
		}
	} else if len(text) == 2 {
		target, _, err := utils.FindUserInMessage(context)
		if err != nil {
			err := context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				return err
			}
			return err
		}
		_, err = utils.Bot.Send(context.Sender(), fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", target.FirstName, target.LastName, target.Username, target.ID))
		if err != nil {
			return err
		}
	} else {
		_, err := utils.Bot.Send(context.Sender(), fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", context.Sender().FirstName, context.Sender().LastName, context.Sender().Username, context.Sender().ID))
		if err != nil {
			return err
		}
	}
	return err
}
