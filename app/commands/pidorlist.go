package commands

import (
	"strconv"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//List add pidors from DB on /pidorlist
func Pidorlist(context telebot.Context) error {
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
	var pidorlist string
	var pidor utils.PidorList
	var i = 0
	result, _ := utils.DB.Model(&utils.PidorList{}).Rows()
	for result.Next() {
		err := utils.DB.ScanRows(result, &pidor)
		if err != nil {
			return err
		}
		i++
		pidorlist += strconv.Itoa(i) + ". @" + pidor.Username + " (" + strconv.FormatInt(pidor.ID, 10) + ")\n"
		if len(pidorlist) > 3900 {
			_, err = utils.Bot.Send(context.Sender(), pidorlist)
			if err != nil {
				return err
			}
			pidorlist = ""
		}
	}
	_, err := utils.Bot.Send(context.Sender(), pidorlist)
	if err != nil {
		return err
	}
	err = context.Reply("Список отправлен в личку.\nЕсли список не пришел, то убедитесь, что бот запущен и не заблокирован в личке.")
	if err != nil {
		return err
	}
	return err
}
