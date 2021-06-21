package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

//Send top 10 pidors of year on /pidorstats
func Pidorstats(m *tb.Message) {
	var text = strings.Split(m.Text, " ")
	var i = 0
	var year = time.Now().Year()
	var username string
	var count int64
	if len(text) == 2 {
		argYear, err := strconv.Atoi(text[1])
		if err != nil {
			_, err := utils.Bot.Reply(m, "Ошибка определения года.\nУкажите год с 2019 по предыдущий.")
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
		if argYear < year && argYear > 2018 {
			year = argYear
		}
	}
	var pidorall = "Топ-10 пидоров за " + strconv.Itoa(year) + " год:\n\n"
	result, _ := utils.DB.Select("username, COUNT(*) as count").Table("pidor_stats, pidor_lists").Where("pidor_stats.user_id=pidor_lists.id").Where("date BETWEEN ? AND ?", time.Date(year, 1, 1, 0, 0, 0, 0, time.Local), time.Date(year+1, 1, 1, 0, 0, 0, 0, time.Local)).Group("user_id").Order("count DESC").Limit(10).Rows()
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
