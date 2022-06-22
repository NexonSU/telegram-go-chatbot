package checkpoint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unicode/utf16"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	m "github.com/keighl/metabolize"
	"github.com/valyala/fastjson"
	"gorm.io/gorm/clause"
	"mvdan.cc/xurls/v2"
)

type MetaData struct {
	Title       string `meta:"og:title"`
	Description string `meta:"og:description,description"`
}

var MaximumIdFromDB = GetMaximumIdFromDB()

func GetMaximumIdFromDB() int64 {
	var user tele.User
	utils.DB.Last(&user).Limit(1)
	return user.ID
}

func CommandGetSpamChance(context tele.Context) error {
	var user tele.User
	var err error
	if len(context.Args()) == 0 && context.Message().ReplyTo == nil {
		user = *context.Sender()
	} else {
		user, _, err = utils.FindUserInMessage(context)
		if err != nil {
			return context.Reply(err.Error())
		}
	}
	spamchance := GetSpamChance(user)
	if spamchance < 0 {
		spamchance = 0
	}
	return context.Reply(fmt.Sprintf("%v спамер на %v%%.", utils.UserFullName(&user), spamchance))
}

func GetSpamChance(user tele.User) int {
	spamchance := 0
	//photos
	photos, _ := utils.Bot.ProfilePhotosOf(&user)
	photoCount := len(photos)
	if photoCount > 5 {
		photoCount = 5
	}
	spamchance -= photoCount*10 - 20
	//ID
	if user.ID > 5000000000 {
		spamchance += int(float64(user.ID)/float64(MaximumIdFromDB)*100) - 50
	} else {
		spamchance += int(float64(user.ID)/float64(2147483647)*100) - 75
	}
	//Bio
	if user.Username != "" {
		res, _ := http.Get(fmt.Sprintf("https://t.me/%v", user.Username))
		data := new(MetaData)
		if m.Metabolize(res.Body, data) == nil {
			if len(data.Description) > 15 && data.Description[:15] == "You can contact" {
				spamchance += 10
			} else {
				spamchance -= 10
				if strings.Contains(data.Description, "http") {
					spamchance += 40
				}
			}
		}
	} else {
		spamchance += 10
	}
	return spamchance
}

func AddAntispam(context tele.Context) error {
	var AntiSpam utils.AntiSpam
	if AntiSpam.Type == "" && xurls.Relaxed().FindString(context.Data()) != "" {
		AntiSpam.Text = xurls.Relaxed().FindString(context.Data())
		AntiSpam.Type = "URL"
	}
	if AntiSpam.Type == "" && context.Message().ReplyTo != nil && context.Message().ReplyTo.Sticker != nil {
		AntiSpam.Text = context.Message().ReplyTo.Sticker.SetName
		AntiSpam.Type = "StickerPack"
	}
	if AntiSpam.Type == "" && context.Data() != "" {
		AntiSpam.Text = context.Data()
		AntiSpam.Type = "Text"
	}
	if AntiSpam.Type == "" {
		return context.Reply("Нужно указать URL, текст или какое-либо сообщение.")
	}
	result := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&AntiSpam)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Ошибка запроса: <code>%v</code>", result.Error.Error()))
	}
	return context.Reply(fmt.Sprintf("%v <code>%v</code> добавлен в антиспам.", AntiSpam.Type, AntiSpam.Text))
}

func DelAntispam(context tele.Context) error {
	var AntiSpam utils.AntiSpam
	if AntiSpam.Type == "" && xurls.Relaxed().FindString(context.Data()) != "" {
		AntiSpam.Text = xurls.Relaxed().FindString(context.Data())
		AntiSpam.Type = "URL"
	}
	if AntiSpam.Type == "" && context.Message().ReplyTo != nil && context.Message().ReplyTo.Sticker != nil {
		AntiSpam.Text = context.Message().ReplyTo.Sticker.SetName
		AntiSpam.Type = "StickerPack"
	}
	if AntiSpam.Type == "" && context.Data() != "" {
		AntiSpam.Text = context.Data()
		AntiSpam.Type = "Text"
	}
	if AntiSpam.Type == "" {
		return context.Reply("Нужно указать URL, текст или какое-либо сообщение.")
	}
	result := utils.DB.Delete(&AntiSpam)
	if result.Error != nil {
		return context.Reply(fmt.Sprintf("Ошибка запроса: <code>%v</code>", result.Error.Error()))
	}
	if result.RowsAffected == 0 {
		return context.Reply("Ошибка: значение не найдено.")
	}
	return context.Reply(fmt.Sprintf("%v <code>%v</code> удалён из антиспама.", AntiSpam.Type, AntiSpam.Text))
}

func ListAntispam(context tele.Context) error {
	var list = "Список фильтров:\n\n"
	result, err := utils.DB.Model(utils.AntiSpam{}).Rows()
	if err != nil {
		return err
	}
	for result.Next() {
		var AntiSpam utils.AntiSpam
		err := utils.DB.ScanRows(result, &AntiSpam)
		if err != nil {
			return err
		}
		list += fmt.Sprintf("%v - %v\n", AntiSpam.Text, AntiSpam.Type)
	}
	return context.Reply(list, &tele.SendOptions{DisableWebPagePreview: true})
}

func SpamFilter(context tele.Context) error {
	if context.Sender().ID == 777000 {
		return nil
	}
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", context.Sender().ID))
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(httpResponse.Body)
	jsonBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}
	if fastjson.GetBool(jsonBytes, "ok") {
		text := fmt.Sprintf("Сообщение пользователя %v было удалено, т.к. он забанен CAS:\n<pre>%v</pre>", utils.MentionUser(context.Sender()), context.Message().Text)
		utils.Bot.Send(tele.ChatID(utils.Config.SysAdmin), text)
		return context.Delete()
	}
	if GetSpamChance(*context.Sender()) < 10 {
		return nil
	}
	var AntiSpam []utils.AntiSpam
	utils.DB.Find(&AntiSpam)
	if context.Message() != nil && context.Message().Sticker != nil {
		for _, AntiSpamEntry := range AntiSpam {
			if AntiSpamEntry.Type == "StickerPack" && context.Message().Sticker.SetName == AntiSpamEntry.Text {
				text := fmt.Sprintf("Стикер пользователя %v был удален, т.к. стикерпак %v запрещен:", utils.MentionUser(context.Sender()), AntiSpamEntry.Text)
				utils.Bot.Send(tele.ChatID(utils.Config.SysAdmin), text)
				utils.Bot.Send(tele.ChatID(utils.Config.SysAdmin), &tele.Sticker{
					File: tele.File{FileID: context.Message().Sticker.FileID},
				})
				return context.Delete()
			}
		}
		return nil
	}
	for _, entity := range context.Message().Entities {
		if entity.Type == "url" {
			url := string(utf16.Decode(utf16.Encode([]rune(context.Message().Text))[entity.Offset : entity.Offset+entity.Length]))
			for _, AntiSpamEntry := range AntiSpam {
				if AntiSpamEntry.Type == "URL" && strings.Contains(strings.ToLower(url), strings.ToLower(AntiSpamEntry.Text)) {
					text := fmt.Sprintf("Сообщение пользователя %v было удалено, т.к. URL запрещен:\n<pre>%v</pre>", utils.MentionUser(context.Sender()), context.Text())
					utils.Bot.Send(tele.ChatID(utils.Config.SysAdmin), text)
					return context.Delete()
				}
			}
		}
	}
	if context.Text() != "" {
		for _, AntiSpamEntry := range AntiSpam {
			if AntiSpamEntry.Type == "Text" && strings.Contains(strings.ToLower(context.Text()), strings.ToLower(AntiSpamEntry.Text)) {
				text := fmt.Sprintf("Сообщение пользователя %v было удалено, т.к. содержит запрещенный текст:\n<pre>%v</pre>", utils.MentionUser(context.Sender()), context.Text())
				utils.Bot.Send(tele.ChatID(utils.Config.SysAdmin), text)
				return context.Delete()
			}
		}
	}
	return nil
}
