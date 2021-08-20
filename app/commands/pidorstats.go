package commands

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send top 10 pidors of year on /pidorstats
func Pidorstats(context telebot.Context) error {
	var i = 0
	var year = time.Now().Year()
	var username string
	var count int64
	if len(context.Args()) == 1 {
		argYear, err := strconv.Atoi(context.Data())
		if err != nil {
			return context.Reply("Ошибка определения года.\nУкажите год с 2019.")
		}
		if argYear == 2077 {
			return context.Reply(&telebot.Video{File: telebot.File{FileID: "BAACAgIAAx0CRXO-MQADWWB4LQABzrOqWPkq-JXIi4TIixY4dwACPw4AArBgwUt5sRu-_fDR5x4E"}})
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
			return err
		}
		i++
		pidorall += fmt.Sprintf("%v. %v - %v раз(а)\n", i, username, count)
	}
	utils.DB.Model(utils.PidorList{}).Count(&count)
	pidorall += fmt.Sprintf("\nВсего участников — %v", count)
	return context.Reply(pidorall)
}
