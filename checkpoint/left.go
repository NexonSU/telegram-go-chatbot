package checkpoint

import (
	"log"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

func UserLeft(context tele.Context) error {
	for _, user := range utils.RestrictedUsers {
		if user.UserID != context.ChatMember().NewChatMember.User.ID {
			continue
		}
		delete := utils.DB.Delete(&user)
		if delete.Error != nil {
			log.Println(delete.Error)
		}
		restricted := utils.DB.Find(&utils.RestrictedUsers)
		if restricted.Error != nil {
			log.Println(restricted.Error)
		}
		if utils.DB.First(&utils.CheckPointRestrict{WelcomeMessageID: user.WelcomeMessageID}).RowsAffected == 0 {
			if utils.WelcomeMessageID == user.WelcomeMessageID {
				utils.WelcomeMessageID = 0
			}
			utils.Bot.Delete(&tele.Message{ID: user.WelcomeMessageID, Chat: &tele.Chat{ID: utils.Config.Chat}})
		}
	}
	return nil
}
