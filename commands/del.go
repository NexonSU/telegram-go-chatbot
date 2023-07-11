package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Delete Get in DB on /del
func Del(context tele.Context) error {
	var get utils.Get
	//args check
	if len(context.Args()) == 0 {
		return utils.ReplyAndRemove("Пример использования: <code>/del {гет}</code>", context)
	}
	//ownership check
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(context.Data())}).First(&get)
	if result.RowsAffected == 0 {
		return utils.ReplyAndRemove(fmt.Sprintf("Гет <code>%v</code> не найден.", context.Data()), context)
	}
	creator, err := utils.GetUserFromDB(fmt.Sprint(get.Creator))
	if err != nil {
		return err
	}
	if get.Creator != context.Sender().ID && !utils.IsAdminOrModer(context.Sender().ID) {
		return utils.ReplyAndRemove(fmt.Sprintf("Данный гет могут изменять либо администраторы, либо %v.", utils.UserFullName(&creator)), context)
	}
	//removing Get
	result = utils.DB.Delete(&get)
	if result.Error != nil {
		return result.Error
	}
	return utils.ReplyAndRemove(fmt.Sprintf("Гет <code>%v</code> удалён.", context.Data()), context)
}
