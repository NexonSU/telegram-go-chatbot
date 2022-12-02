package commands

import (
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// test
func Test(context tele.Context) error {
	var ricochetVictim *tele.ChatMember
	var message utils.Message
	var err error
	for i := 1; i < 100; i++ {
		result := utils.DB.Where(utils.Message{ChatID: context.Chat().ID}).Order("id desc").Group("user_id").Offset(i).Last(&message)
		if result.Error != nil {
			continue
		}
		ricochetVictim, err = utils.Bot.ChatMemberOf(context.Chat(), &tele.User{ID: message.UserID})
		if err != nil {
			continue
		}
	}
	return context.Reply(fmt.Sprint(ricochetVictim.User))
}
