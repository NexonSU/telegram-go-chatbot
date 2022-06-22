package pidor

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Remove user in DB on /pidordel
func Pidordel(context tele.Context) error {
	var pidor utils.PidorList
	user, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	pidor = utils.PidorList(user)
	result := utils.DB.Delete(&pidor)
	if result.RowsAffected != 0 {
		return context.Reply(fmt.Sprintf("Пользователь %v удалён из игры <b>Пидор Дня</b>!", utils.MentionUser(&user)))
	} else {
		return context.Reply(fmt.Sprintf("Не удалось удалить пользователя:\n<code>%v</code>", result.Error.Error()))
	}
}
