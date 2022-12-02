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
		message = utils.Message{}
		result := utils.DB.Where(utils.Message{ChatID: context.Chat().ID}).Order("id desc").Offset(i).Last(&message)
		if result.Error != nil {
			continue
		}
		ricochetVictim = &tele.ChatMember{}
		ricochetVictim, err = utils.Bot.ChatMemberOf(context.Chat(), &tele.User{ID: message.UserID})
		if err != nil {
			continue
		}
		if ricochetVictim.Role == "member" {
			return context.Reply(fmt.Sprint(ricochetVictim.User))
		}
	}
	return context.Reply("can't find member user")
}
