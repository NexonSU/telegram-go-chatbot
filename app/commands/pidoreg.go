package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

//Send DB result on /pidoreg
func Pidoreg(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	var pidor utils.PidorList
	result := utils.DB.First(&pidor, context.Sender().ID)
	if result.RowsAffected != 0 {
		err := context.Reply("Эй, ты уже в игре!")
		if err != nil {
			return err
		}
	} else {
		pidor = utils.PidorList(*context.Sender())
		result = utils.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(pidor)
		if result.Error != nil {
			err := context.Reply(fmt.Sprintf("Не удалось зарегистрироваться:\n<code>%v</code>.", result.Error))
			if err != nil {
				return err
			}
			return err
		}
		err := context.Reply("OK! Ты теперь участвуешь в игре <b>Пидор Дня</b>!")
		if err != nil {
			return err
		}
	}
	return err
}
