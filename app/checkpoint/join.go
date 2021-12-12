package checkpoint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/valyala/fastjson"
	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

type welcomeMessage struct {
	ID    int
	time  int64
	users int
	text  string
}

var WelcomeMessage welcomeMessage

func UserJoin(context telebot.Context) error {
	//joined user
	User := context.ChatMember().NewChatMember.User
	//CAS ban check
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", User.ID))
	if err != nil {
		_ = utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, User)
		return err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(httpResponse.Body)
	jsonBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		_ = utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, User)
		return err
	}
	if fastjson.GetBool(jsonBytes, "ok") {
		err := utils.Bot.Ban(&telebot.Chat{ID: utils.Config.Chat}, &telebot.ChatMember{User: User})
		if err != nil {
			_ = utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, User)
			return err
		}
	}
	//user chat restrict
	restrictUser := utils.CheckPointRestrict{UserID: User.ID, Since: time.Now().Unix()}
	restrict := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&restrictUser)
	if restrict.Error != nil {
		_ = utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, User)
		return restrict.Error
	}
	restricted := utils.DB.Find(&utils.RestrictedUsers)
	if restricted.Error != nil {
		return restricted.Error
	}
	//welcome message text
	var welcomeGet utils.Get
	utils.DB.Where(&utils.Get{Name: "welcome"}).First(&welcomeGet)
	if welcomeGet.Data == "" {
		welcomeGet.Data = "Ответь на это сообщение и расскажи о себе."
	}
	WelcomeMessage.users++
	//welcome message create\update
	if time.Now().Unix()-WelcomeMessage.time > 60 && utils.LastChatMessageID-WelcomeMessage.ID > 0 {
		WelcomeMessage.time = time.Now().Unix()
		WelcomeMessage.users = 1
		WelcomeMessage.text = fmt.Sprintf("Привет %v!", utils.MentionUser(User))
		m, err := utils.Bot.Send(&telebot.Chat{ID: utils.Config.Chat}, WelcomeMessage.text+"\n"+welcomeGet.Data)
		if err != nil {
			_ = utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, User)
			return err
		}
		WelcomeMessage.ID = m.ID
		utils.WelcomeMessageID = m.ID
	} else if len(WelcomeMessage.text) < 3500 &&
		!strings.ContainsAny(WelcomeMessage.text, fmt.Sprint(User.ID)) {
		WelcomeMessage.text = strings.Replace(WelcomeMessage.text, "Привет ", fmt.Sprintf("Привет %v, ", utils.MentionUser(User)), 1)
		if WelcomeMessage.users > 5 && time.Now().Unix()-WelcomeMessage.time < 5 {
			return nil
		}
		WelcomeMessage.time = time.Now().Unix()
		_, err := utils.Bot.Edit(&telebot.Message{ID: WelcomeMessage.ID, Chat: &telebot.Chat{ID: utils.Config.Chat}}, WelcomeMessage.text+"\n"+welcomeGet.Data)
		if err != nil {
			_ = utils.Bot.Unban(&telebot.Chat{ID: utils.Config.Chat}, User)
			return err
		}
	}
	return nil
}
