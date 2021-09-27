package stats

import (
	"fmt"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func PopWords(context telebot.Context) error {
	result, _ := utils.DB.Model(utils.Word{ChatID: context.Message().Chat.ID}).Select("text, COUNT(*) as count").Group("text").Order("count DESC").Limit(10).Rows()
	var popwords string
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
