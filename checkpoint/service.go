package checkpoint

import (
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

var welcomeMessageText = ""
var welcomeGet = utils.Get{Data: "Ответь на это сообщение и расскажи о себе."}

func welcomeMessageUpdate() error {
	welcomeMessageUsers := ""
	restricted := utils.DB.Find(&utils.RestrictedUsers)
	if restricted.Error != nil {
		return restricted.Error
	}
	for _, user := range utils.RestrictedUsers {
		if time.Now().Unix()-user.Since > 120 {
			delete := utils.DB.Delete(&user)
			if delete.Error != nil {
				return delete.Error
			}
			err := utils.Bot.Ban(&tele.Chat{ID: utils.Config.Chat}, &tele.ChatMember{User: &tele.User{ID: user.UserID}, RestrictedUntil: time.Now().Unix() + 21600})
			if err != nil {
				return err
			}
			continue
		}
		if utf8.RuneCountInString(welcomeMessageUsers) < 2000 {
			welcomeMessageUsers = fmt.Sprintf("%v, %v", welcomeMessageUsers, utils.MentionUser(&tele.User{ID: user.UserID, FirstName: user.UserFirstName, LastName: user.UserLastName}))
		}
	}
	//usertext & welcomeMessage check
	if welcomeMessageUsers == "" {
		if utils.WelcomeMessageID != 0 {
			err := utils.Bot.Delete(&tele.Message{ID: utils.WelcomeMessageID, Chat: &tele.Chat{ID: utils.Config.Chat}})
			if err != nil {
				return err
			}
			utils.WelcomeMessageID = 0
		}
		return nil
	}
	if utf8.RuneCountInString(welcomeMessageUsers) > 2000 {
		welcomeMessageUsers = fmt.Sprintf("%v и другие уважаемые цыгане!\nБотов в очереди: %v", welcomeMessageUsers, len(utils.RestrictedUsers))
	} else {
		welcomeMessageUsers = fmt.Sprintf("%v!", welcomeMessageUsers)
	}
	//welcome message text
	utils.DB.Where(&utils.Get{Name: "welcome"}).First(&welcomeGet)
	//welcome message create\update
	if utils.WelcomeMessageID == 0 {
		welcomeMessageText = fmt.Sprintf("Привет%v\n%v", welcomeMessageUsers, welcomeGet.Data)
		m, err := utils.Bot.Send(&tele.Chat{ID: utils.Config.Chat}, welcomeMessageText, &tele.SendOptions{DisableWebPagePreview: true})
		if err != nil {
			return err
		}
		utils.WelcomeMessageID = m.ID
	} else if welcomeMessageText != fmt.Sprintf("Привет%v\n%v", welcomeMessageUsers, welcomeGet.Data) {
		welcomeMessageText = fmt.Sprintf("Привет%v\n%v", welcomeMessageUsers, welcomeGet.Data)
		_, err := utils.Bot.Edit(&tele.Message{ID: utils.WelcomeMessageID, Chat: &tele.Chat{ID: utils.Config.Chat}}, welcomeMessageText, &tele.SendOptions{DisableWebPagePreview: true})
		if err != nil {
			return err
		}
	}
	return nil
}

func welcomeMessageUpdateService() {
	for {
		err := welcomeMessageUpdate()
		if err != nil {
			log.Print("welcomeMessageUpdate: ")
			log.Println(err.Error())
		}
		time.Sleep(time.Second * time.Duration(2))
	}
}

func init() {
	go welcomeMessageUpdateService()
}
