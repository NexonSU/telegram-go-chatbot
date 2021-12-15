package checkpoint

import (
	"log"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

func restrictUpdate() error {
	restricted := utils.DB.Find(&utils.RestrictedUsers)
	if restricted.Error != nil {
		log.Println(restricted.Error)
	}
	for _, user := range utils.RestrictedUsers {
		if user.Since > time.Now().Unix()-120 {
			continue
		}
		delete := utils.DB.Delete(&user)
		if delete.Error != nil {
			log.Println(delete.Error)
		}
		err := utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, &telebot.User{ID: user.UserID})
		if err != nil {
			log.Println(err)
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

func restrictService(init bool) error {
	if init {
		go restrictService(false)
		return nil
	}
	for {
		err := restrictUpdate()
		if err != nil {
			log.Println(err.Error())
		}
		time.Sleep(time.Second * time.Duration(60))
	}
}

var _ = restrictService(true)
