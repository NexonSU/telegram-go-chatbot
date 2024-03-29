package pidor

import (
	"strconv"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// List add pidors from DB on /pidorlist
func Pidorlist(context tele.Context) error {
	var pidorlist string
	var pidor utils.PidorList
	var i = 0
	var err error
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
	_, err = utils.Bot.Send(context.Sender(), pidorlist)
	if err != nil {
		return err
	}
	return utils.ReplyAndRemove("Список отправлен в личку.\nЕсли список не пришел, то убедитесь, что бот запущен и не заблокирован в личке.", context)
}
