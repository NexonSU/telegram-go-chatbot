package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send top 10 pidors of all time on /pidorall
func Pidorall(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var i = 0
	var username string
	var count int64
	var pidorall = "Топ-10 пидоров за всё время:\n\n"
	result, _ := utils.DB.Select("username, COUNT(*) as count").Table("pidor_stats, pidor_lists").Where("pidor_stats.user_id=pidor_lists.id").Group("user_id").Order("count DESC").Limit(10).Rows()
	for result.Next() {
		err := result.Scan(&username, &count)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		i++
		pidorall += fmt.Sprintf("%v. %v - %v раз(а)\n", i, username, count)
	}
	utils.DB.Model(utils.PidorList{}).Count(&count)
	pidorall += fmt.Sprintf("\nВсего участников — %v", count)
	_, err := utils.Bot.Reply(m, pidorall)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
