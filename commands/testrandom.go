package commands

import (
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Kill user on /blessing, /suicide
func TestRandom(context tele.Context) error {
	text := "1000xRandInt(0, 9):\n"
	numbers := [10]int{}
	for i := 0; i < 1000; i++ {
		numbers[utils.RandInt(0, 9)] += 1
	}
	for number, count := range numbers {
		text = fmt.Sprintf("%v%v - %v\n", text, number, count)
	}
	return context.Reply(text)
}
