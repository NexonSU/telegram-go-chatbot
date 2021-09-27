package stats

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func PopWords(context telebot.Context) error {
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
	popwords := fmt.Sprintf("Самые популярные слова c %v:\n", from.Format("02.01.2006"))
	result, _ := utils.DB.Model(utils.Word{ChatID: context.Message().Chat.ID}).Select("text, COUNT(*) as count").Where("date BETWEEN ? AND ?", from, to).Group("text").Order("count DESC").Limit(10).Rows()
	var word string
	var count int
	for result.Next() {
		err := result.Scan(&word, &count)
		if err != nil {
			return err
		}
		popwords += fmt.Sprintf("%v	-	%v\n", count, word)
	}
	return context.Reply(popwords)
}
