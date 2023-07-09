package pidor

import (
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

// Send DB result on /pidoreg
func Pidoreg(context tele.Context) error {
	var pidor utils.PidorList
	if utils.DB.First(&pidor, context.Sender().ID).RowsAffected != 0 {
		return utils.SendAndRemove("Эй, ты уже в игре!", context)
	} else {
		pidor = utils.PidorList(*context.Sender())
		result := utils.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(&pidor)
		if result.Error != nil {
			return result.Error
		}
		return utils.SendAndRemove("OK! Ты теперь участвуешь в игре <b>Пидор Дня</b>!", context)
	}
}
