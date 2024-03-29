package bets

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Remove bet
func DelBet(context tele.Context) error {
	var bet utils.Bets
	if len(context.Args()) < 2 {
		return utils.ReplyAndRemove("Пример использования: <code>/delbet 30.06.2023 ставлю жопу, что TESVI будет говном</code>", context)
	}
	date, err := time.Parse("02.01.2006", context.Args()[0])
	if err != nil {
		return err
	}
	bet.UserID = context.Sender().ID
	bet.Timestamp = date.Unix()
	bet.Text = strings.Join(context.Args()[1:], " ")
	if err != nil {
		return err
	}
	result := utils.DB.Delete(&bet)
	if result.RowsAffected != 0 {
		return utils.ReplyAndRemove(fmt.Sprintf("Ставка удалена:\n%v, %v:<pre>%v</pre>\n", time.Unix(bet.Timestamp, 0).Format("02.01.2006"), utils.UserFullName(context.Sender()), html.EscapeString(bet.Text)), context)
	} else {
		return utils.ReplyAndRemove("Твоя ставка не найдена по указанным параметрам.", context)
	}
}
