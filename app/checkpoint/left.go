package checkpoint

import (
	"log"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func UserLeft(context telebot.Context) error {
	for _, user := range utils.RestrictedUsers {
		if user.UserID != context.Sender().ID {
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
			utils.Bot.Delete(&telebot.Message{ID: user.WelcomeMessageID, Chat: &telebot.Chat{ID: utils.Config.Chat}})
		}
	}
	return nil
}
