package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm/clause"
)

//Send DB result on /pidoreg
func Pidoreg(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var pidor utils.PidorList
	result := utils.DB.First(&pidor, m.Sender.ID)
	if result.RowsAffected != 0 {
		_, err := utils.Bot.Reply(m, "Эй, ты уже в игре!")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else {
		pidor = utils.PidorList(*m.Sender)
		result = utils.DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(pidor)
		if result.Error != nil {
			utils.ErrorReporting(result.Error, m)
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось зарегистрироваться:\n<code>%v</code>.", result.Error))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
		_, err := utils.Bot.Reply(m, "OK! Ты теперь участвуешь в игре <b>Пидор Дня</b>!")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
