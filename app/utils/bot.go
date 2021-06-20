package utils

import (
	"errors"
	tb "gopkg.in/tucnak/telebot.v2"
	"strconv"
	"strings"
	"time"
)

var Bot, _ = tb.NewBot(tb.Settings{
	URL:       Config.Telegram.BotApiUrl,
	Token:     Config.Telegram.Token,
	ParseMode: tb.ModeHTML,
	Poller: &tb.LongPoller{
		Timeout:        10 * time.Second,
		AllowedUpdates: Config.Webhook.AllowedUpdates,
	},
})

func FindUserInMessage(m tb.Message) (tb.User, int64, error) {
	var user tb.User
	var err error = nil
	var untildate = time.Now().Unix()
	var text = strings.Split(m.Text, " ")
	if m.ReplyTo != nil {
		user = *m.ReplyTo.Sender
		if len(text) == 2 {
			addtime, err := strconv.ParseInt(text[1], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	} else {
		if len(text) == 1 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		user, err = GetUserFromDB(text[1])
		if err != nil {
			return user, untildate, err
		}
		if len(text) == 3 {
			addtime, err := strconv.ParseInt(text[2], 10, 64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	}
	return user, untildate, err
}
