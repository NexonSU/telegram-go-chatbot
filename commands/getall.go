package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Send list of Gets to user on /getall
func Getall(context tele.Context) error {
	var getall string
	var get utils.Get
	result, _ := utils.DB.Model(&utils.Get{}).Rows()
	for result.Next() {
		if getall == "" {
			getall = "Доступные геты: "
		} else {
			getall += ", "
		}
		err := utils.DB.ScanRows(result, &get)
		if err != nil {
			return err
		}
		getall += get.Name
		if len([]rune(getall)) > 4000 {
			utils.Bot.Send(context.Sender(), getall)
			getall = ""
		}
	}
	utils.Bot.Send(context.Sender(), getall)
	return utils.SendAndRemove("Список отправлен в личку.\nЕсли список не пришел, то убедитесь, что бот запущен и не заблокирован в личке.", context)
}
