package pidor

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

//Send DB stats on /pidorme
func Pidorme(context tele.Context) error {
	var pidor utils.PidorStats
	var countYear int64
	var countAlltime int64
	pidor.UserID = context.Sender().ID
	utils.DB.Model(&utils.PidorStats{}).Where(pidor).Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local), time.Now()).Count(&countYear)
	utils.DB.Model(&utils.PidorStats{}).Where(pidor).Count(&countAlltime)
	return context.Reply(fmt.Sprintf("В этом году ты был пидором дня — %v раз!\nЗа всё время ты был пидором дня — %v раз!", countYear, countAlltime))
}
