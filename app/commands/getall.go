package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

//Send list of Gets to user on /getall
func Getall(m *tb.Message) {
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
	return
}
