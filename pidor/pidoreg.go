package pidor

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	"gorm.io/gorm/clause"
)

//Send DB result on /pidoreg
func Pidoreg(context tele.Context) error {
	var pidor utils.PidorList
	if utils.DB.First(&pidor, context.Sender().ID).RowsAffected != 0 {
		return context.Reply("Эй, ты уже в игре!")
	} else {
		pidor = utils.PidorList(*context.Sender())
		result := utils.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&pidor)
		if result.Error != nil {
			return context.Reply(fmt.Sprintf("Не удалось зарегистрироваться:\n<code>%v</code>.", result.Error))
		}
		return context.Reply("OK! Ты теперь участвуешь в игре <b>Пидор Дня</b>!")
	}
}
