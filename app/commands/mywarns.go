package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

//Send warning amount on /mywarns
func Mywarns(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var warn utils.Warn
	result := utils.DB.First(&warn, m.Sender.ID)
	if result.RowsAffected != 0 {
		warn.Amount = warn.Amount - int(time.Now().Sub(warn.LastWarn).Hours()/24/7)
		if warn.Amount < 0 {
			warn.Amount = 0
		}
	} else {
		warn.UserID = m.Sender.ID
		warn.LastWarn = time.Unix(0, 0)
		warn.Amount = 0
	}
	warnStrings := []string{"предупреждений", "предупреждение", "предупреждения", "предупреждения"}
	_, err := utils.Bot.Reply(m, fmt.Sprintf("У тебя %v %v.", warn.Amount, warnStrings[warn.Amount]))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
