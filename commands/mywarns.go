package commands

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Send warning amount on /mywarns
func Mywarns(context tele.Context) error {
	// prt will replace fmt package to format text according plurals defined in utils package
	// If no plural rule matched it will be ignored and processed as usual formatting
	prt := message.NewPrinter(language.Russian)

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
	return utils.ReplyAndRemove(prt.Sprintf("У тебя %d предупреждений.", warn.Amount), context)
}
