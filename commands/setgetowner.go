package commands

import (
	"fmt"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"gopkg.in/telebot.v3"
)

//Send Get to user on /get
func SetGetOwner(context telebot.Context) error {
	var get utils.Get
	if len(context.Args()) != 1 || context.Message().ReplyTo == nil {
		return context.Reply("Пример использования: <code>/setgetowner {гет}</code> в ответ пользователю, которого нужно задать владельцем.")
	}
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(context.Args()[0])}).First(&get)
	if result.RowsAffected != 0 {
		get.Creator = context.Message().ReplyTo.Sender.ID
		utils.DB.First(&get)
		if result.Error != nil {
			return context.Reply(fmt.Sprintf("Не удалось сохранить гет <code>%v</code>.", get.Name))
		}
		return context.Reply(fmt.Sprintf("Владелец гета <code>%v</code> изменён на %v.", get.Name, utils.MentionUser(context.Message().ReplyTo.Sender)))
	} else {
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> не найден.", context.Data()))
	}
}
