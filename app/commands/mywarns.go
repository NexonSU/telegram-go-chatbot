package commands

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send warning amount on /mywarns
func Mywarns(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	var warn utils.Warn
	result := utils.DB.First(&warn, context.Sender().ID)
	if result.RowsAffected != 0 {
		warn.Amount = warn.Amount - int(time.Since(warn.LastWarn).Hours()/24/7)
		if warn.Amount < 0 {
			warn.Amount = 0
		}
	} else {
		warn.UserID = context.Sender().ID
		warn.LastWarn = time.Unix(0, 0)
		warn.Amount = 0
	}
	warnStrings := []string{"предупреждений", "предупреждение", "предупреждения", "предупреждения"}
	return context.Reply(fmt.Sprintf("У тебя %v %v.", warn.Amount, warnStrings[warn.Amount]))
}
