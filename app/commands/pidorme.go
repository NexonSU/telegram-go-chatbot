package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

//Send DB stats on /pidorme
func Pidorme(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var pidor utils.PidorStats
	var countYear int64
	var countAlltime int64
	pidor.UserID = m.Sender.ID
	utils.DB.Model(&utils.PidorStats{}).Where(pidor).Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local), time.Now()).Count(&countYear)
	utils.DB.Model(&utils.PidorStats{}).Where(pidor).Count(&countAlltime)
	_, err := utils.Bot.Reply(m, fmt.Sprintf("В этом году ты был пидором дня — %v раз!\nЗа всё время ты был пидором дня — %v раз!", countYear, countAlltime))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
