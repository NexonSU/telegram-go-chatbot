package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send list of Gets to user on /getall
func Getall(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var getall []string
	var get utils.Get
	result, _ := utils.DB.Model(&utils.Get{}).Rows()
	for result.Next() {
		err := utils.DB.ScanRows(result, &get)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		getall = append(getall, get.Name)
	}
	_, err := utils.Bot.Reply(m, fmt.Sprintf("Доступные геты: %v", strings.Join(getall[:], ", ")))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
