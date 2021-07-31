package commands

import (
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send text in chat on /say
func Say(context telebot.Context) error {
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
	if len(text) > 1 {
		err := utils.Bot.Delete(context.Message())
		if err != nil {
			return err
		}
		_, err = utils.Bot.Send(context.Chat(), strings.Join(text[1:], " "))
		if err != nil {
			return err
		}
	} else {
		err := context.Reply("Укажите сообщение.")
		if err != nil {
			return err
		}
	}
	return err
}
