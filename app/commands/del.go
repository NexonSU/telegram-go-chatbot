package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Delete Get in DB on /del
func Del(context telebot.Context) error {
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
	if len(text) != 2 {
		err := context.Reply("Пример использования: <code>/del {гет}</code>")
		if err != nil {
			return err
		}
		return err
	}
	result := utils.DB.Delete(&utils.Get{Name: strings.ToLower(text[1])})
	if result.RowsAffected != 0 {
		err := context.Reply(fmt.Sprintf("Гет <code>%v</code> удалён.", text[1]))
		if err != nil {
			return err
		}
	} else {
		err := context.Reply(fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
		if err != nil {
			return err
		}
	}
	return err
}
