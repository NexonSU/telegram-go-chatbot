package bets

import (
	"fmt"
	"strconv"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// List all bets
func AllBets(context tele.Context) error {
	var betlist string
	var bet utils.Bets
	var user tele.User
	var i = 0
	var from int64
	if len(context.Args()) > 0 {
		if context.Args()[0] == "all" {
			from = 0
		}
	}
	from = time.Now().Local().Unix() - 86400
	result, _ := utils.DB.Model(&utils.Bets{}).Where("timestamp > ?", from).Rows()
	for result.Next() {
		err := utils.DB.ScanRows(result, &bet)
		if err != nil {
			return err
		}
		i++
		user, err = utils.GetUserFromDB(strconv.FormatInt(bet.UserID, 10))
		if err != nil {
			return err
		}
		betlist += fmt.Sprintf("%v, %v:<pre>%v</pre>\n", time.Unix(bet.Timestamp, 0).Format("02.01.2006"), utils.UserFullName(&user), bet.Text)
		if len(betlist) > 3900 {
			err = context.Reply(betlist)
			if err != nil {
				return err
			}
			betlist = ""
		}
	}
	return context.Reply(betlist)
}
