package commands

import (
	"strconv"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//List add pidors from DB on /pidorlist
func Pidorlist(m *tb.Message) {
	if !utils.IsAdminOrModer(m.Sender.Username) {
		if m.Chat.Username != utils.Config.Telegram.Chat {
			return
		}
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var pidorlist string
	var pidor utils.PidorList
	var i = 0
	result, _ := utils.DB.Model(&utils.PidorList{}).Rows()
	for result.Next() {
		err := utils.DB.ScanRows(result, &pidor)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		i++
		pidorlist += strconv.Itoa(i) + ". @" + pidor.Username + " (" + strconv.Itoa(pidor.ID) + ")\n"
		if len(pidorlist) > 3900 {
			_, err = utils.Bot.Send(m.Sender, pidorlist)
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			pidorlist = ""
		}
	}
	_, err := utils.Bot.Send(m.Sender, pidorlist)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	_, err = utils.Bot.Reply(m, "Список отправлен в личку.\nЕсли список не пришел, то убедитесь, что бот запущен и не заблокирован в личке.")
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
