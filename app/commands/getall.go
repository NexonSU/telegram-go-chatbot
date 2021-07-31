package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send list of Gets to user on /getall
func Getall(context telebot.Context) error {
	var getall []string
	var get utils.Get
	result, _ := utils.DB.Model(&utils.Get{}).Rows()
	for result.Next() {
		err := utils.DB.ScanRows(result, &get)
		if err != nil {
			return err
		}
		getall = append(getall, get.Name)
	}
	return context.Reply(fmt.Sprintf("Доступные геты: %v", strings.Join(getall[:], ", ")))
}
