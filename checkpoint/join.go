package checkpoint

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/valyala/fastjson"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

func UserJoin(context tele.Context) error {
	//joined user
	User := context.ChatMember().NewChatMember.User
	//kick user
	return utils.Bot.Unban(&tele.Chat{ID: utils.Config.Chat}, User)
	//CAS ban check
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", User.ID))
	if err != nil {
		_ = utils.Bot.Unban(&tele.Chat{ID: utils.Config.Chat}, User)
		return err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(httpResponse.Body)
	jsonBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		_ = utils.Bot.Unban(&tele.Chat{ID: utils.Config.Chat}, User)
		return err
	}
	if fastjson.GetBool(jsonBytes, "ok") {
		err := utils.Bot.Ban(&tele.Chat{ID: utils.Config.Chat}, &tele.ChatMember{User: User})
		if err != nil {
			_ = utils.Bot.Unban(&tele.Chat{ID: utils.Config.Chat}, User)
			return err
		}
	}
	//user chat restrict
	restrictUser := utils.CheckPointRestrict{
		UserID:        User.ID,
		UserFirstName: User.FirstName,
		UserLastName:  User.LastName,
		Since:         time.Now().Unix(),
	}
	restrict := utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&restrictUser)
	if restrict.Error != nil {
		_ = utils.Bot.Unban(&tele.Chat{ID: utils.Config.Chat}, User)
		return restrict.Error
	}
	restricted := utils.DB.Find(&utils.RestrictedUsers)
	if restricted.Error != nil {
		return restricted.Error
	}
	return nil
}
