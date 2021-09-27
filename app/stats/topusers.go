package stats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func TopUsers(context telebot.Context) error {
	days := 30
	if len(context.Args()) == 1 {
		var err error
		days, err = strconv.Atoi(context.Data())
		if err != nil {
			return context.Reply("Ошибка определения дней.")
		}
		if days == 2077 {
			return context.Reply(&telebot.Video{File: telebot.File{FileID: "BAACAgIAAx0CRXO-MQADWWB4LQABzrOqWPkq-JXIi4TIixY4dwACPw4AArBgwUt5sRu-_fDR5x4E"}})
		}
	}
	days = days * -1
	from := time.Now().AddDate(0, 0, days)
	to := time.Now().Add(time.Hour * 24)
	popwords := fmt.Sprintf("Самые активные юзеры c %v:\n", from.Format("02.01.2006"))
	result, _ := utils.DB.Model(utils.Message{ChatID: context.Message().Chat.ID}).Select("user_id, COUNT(*) as count").Where("date BETWEEN ? AND ?", from, to).Group("user_id").Order("count DESC").Limit(10).Rows()
	var UserID int
	var count int
	for result.Next() {
		err := result.Scan(&UserID, &count)
		if err != nil {
			return err
		}
		user, err := utils.GetUserFromDB(strconv.Itoa(UserID))
		if err != nil {
			return err
		}
		popwords += fmt.Sprintf("%v	-	%v\n", count, user.FullName())
	}
	return context.Reply(popwords)
}
