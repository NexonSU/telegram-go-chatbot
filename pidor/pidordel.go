package pidor

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Remove user in DB on /pidordel
func Pidordel(context tele.Context) error {
	var pidor utils.PidorList
	user, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	pidor = utils.PidorList(user)
	result := utils.DB.Delete(&pidor)
	if result.RowsAffected != 0 {
		return utils.SendAndRemove(fmt.Sprintf("Пользователь %v удалён из игры <b>Пидор Дня</b>!", utils.MentionUser(&user)), context)
	} else {
		return utils.SendAndRemove(fmt.Sprintf("Не удалось удалить пользователя:\n<code>%v</code>", result.Error.Error()), context)
	}
}
