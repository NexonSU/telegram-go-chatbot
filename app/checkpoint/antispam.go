package checkpoint

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	m "github.com/keighl/metabolize"
	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
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
	return context.Reply(fmt.Sprintf("%v спамер на %v процентов.", utils.UserFullName(&user), spamchance))
}

func AddToWhiteList(context telebot.Context) error {
	if context.Data() == "" {
		return context.Reply("Нужно указать URL или его часть.")
	}
	var link utils.AntiSpamLink
	link.URL = context.Data()
	link.Type = "whitelist"
	result := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&link)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Ошибка запроса: <code>%v</code>", result.Error.Error()))
	}
	return context.Reply(fmt.Sprintf("URL <code>%v</code> добавлен в белый список.", link.URL))
}

func AddToBlackList(context telebot.Context) error {
	if context.Data() == "" {
		return context.Reply("Нужно указать URL или его часть.")
	}
	var link utils.AntiSpamLink
	link.URL = context.Data()
	link.Type = "blacklist"
	result := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&link)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Ошибка запроса: <code>%v</code>", result.Error.Error()))
	}
	return context.Reply(fmt.Sprintf("URL <code>%v</code> добавлен в черный список.", link.URL))
}

func ListAntispamLinks(context telebot.Context) error {
	var list = "Список URL фильтров:\n\n"
	result, err := utils.DB.Model(utils.AntiSpamLink{}).Rows()
	if err != nil {
		return err
	}
	for result.Next() {
		var link utils.AntiSpamLink
		err := utils.DB.ScanRows(result, &link)
		if err != nil {
			return err
		}
		list += fmt.Sprintf("%v - %v\n", link.URL, link.Type)
	}
	return context.Reply(list, telebot.SendOptions{DisableWebPagePreview: true})
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

func UrlFilter(context telebot.Context) error {
	for _, entity := range context.Message().Entities {
		if entity.Type == "url" {
			var link utils.AntiSpamLink
			runes := []rune(context.Message().Text)
			url := string(runes[entity.Offset : entity.Offset+entity.Length])
			result := utils.DB.Where("url LIKE ?", url).First(&link)
			if result.Error != nil {
				return nil
			}
			if link.Type == "blacklist" {
				return context.Delete()
			}
			if GetSpamChance(*context.Sender()) > 50 && result.RowsAffected == 0 {
				return context.Delete()
			}
		}
	}
	return nil
}
