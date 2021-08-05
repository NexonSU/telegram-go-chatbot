package checkpoint

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	m "github.com/keighl/metabolize"
	"gopkg.in/tucnak/telebot.v3"
)

type MetaData struct {
	Title       string `meta:"og:title"`
	Description string `meta:"og:description,description"`
}

var MaximumIdFromDB = GetMaximumIdFromDB()

func GetMaximumIdFromDB() int64 {
	var user telebot.User
	utils.DB.Last(&user).Limit(1)
	return user.ID
}

func CommandGetSpamChance(context telebot.Context) error {
	user, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(err.Error())
	}
	spamchance := GetSpamChance(user)
	return context.Reply(fmt.Sprintf("%v", spamchance))
}

func GetSpamChance(user telebot.User) int {
	spamchance := 0
	//photos
	photos, _ := utils.Bot.ProfilePhotosOf(&user)
	photoCount := len(photos)
	if photoCount > 5 {
		photoCount = 5
	}
	spamchance -= photoCount*10 - 20
	log.Printf("%v - %v photos - %v", user.FirstName, photoCount, spamchance)
	//ID
	spamchance += int(float64(user.ID)/float64(MaximumIdFromDB)*100) - 50
	log.Printf("%v - id %v - %v", user.FirstName, user.ID, spamchance)
	//Bio
	if user.Username != "" {
		res, _ := http.Get(fmt.Sprintf("https://t.me/%v", user.Username))
		data := new(MetaData)
		if m.Metabolize(res.Body, data) == nil {
			if len(data.Description) > 15 && data.Description[:15] == "You can contact" {
				spamchance += 10
				log.Printf("%v - no bio - %v", user.FirstName, spamchance)
			} else {
				spamchance -= 10
				log.Printf("%v - has bio - %v", user.FirstName, spamchance)
				if strings.Contains(data.Description, "http") {
					spamchance += 40
				}
			}
		}
	} else {
		spamchance += 10
		log.Printf("%v - no username - %v", user.FirstName, spamchance)
	}
	return spamchance
}
