package pidor

import (
	"strconv"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

//List add pidors from DB on /pidorlist
func Pidorlist(context telebot.Context) error {
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
				return context.Reply("Ошибка отправки списка. Убедитесь, что бот запущен и не заблокирован в личке.")
			}
			pidorlist = ""
		}
	}
	return context.Reply("Список отправлен в личку.\nЕсли список не пришел, то убедитесь, что бот запущен и не заблокирован в личке.")
}
