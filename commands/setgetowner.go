package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Send Get to user on /get
func SetGetOwner(context tele.Context) error {
	var get utils.Get
	if len(context.Args()) != 1 || context.Message().ReplyTo == nil {
		return utils.ReplyAndRemove("Пример использования: <code>/setgetowner {гет}</code> в ответ пользователю, которого нужно задать владельцем.", context)
	}
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(context.Args()[0])}).First(&get)
	if result.RowsAffected != 0 {
		get.Creator = context.Message().ReplyTo.Sender.ID
		utils.DB.First(&get)
		if result.Error != nil {
			return result.Error
		}
		return utils.ReplyAndRemove(fmt.Sprintf("Владелец гета <code>%v</code> изменён на %v.", get.Name, utils.MentionUser(context.Message().ReplyTo.Sender)), context)
	} else {
		return utils.ReplyAndRemove(fmt.Sprintf("Гет <code>%v</code> не найден.", context.Data()), context)
	}
}
