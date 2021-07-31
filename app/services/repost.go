package services

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Repost channel post to chat
func OnPost(context telebot.Context) error {
	var err error
	if context.Chat().Username == utils.Config.Telegram.Channel {
		chat, err := utils.Bot.ChatByID("@" + utils.Config.Telegram.Chat)
		if err != nil {
			return err
		}
		_, err = utils.Bot.Forward(chat, m)
		if err != nil {
			return err
		}
	}
	return err
}
