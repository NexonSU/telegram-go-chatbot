package bets

import (
	"fmt"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Add bet
func Bet(context tele.Context) error {
	var bet utils.Bets
	if len(context.Args()) < 2 {
		return context.Reply("Пример использования: <code>/bet 30.06.2023 ставлю жопу, что TESVI будет говном</code>")
	}
	date, err := time.Parse("02.01.2006", context.Args()[0])
	if err != nil {
		return err
	}
	if date.Unix() < time.Now().Local().Unix() {
		return fmt.Errorf("минимальная дата: %v", time.Now().Local().Format("02.01.2006"))
	}
	bet.UserID = context.Sender().ID
	bet.Timestamp = date.Unix()
	bet.Text = strings.Join(context.Args()[1:], " ")
	result := utils.DB.Create(&bet)
	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("такая ставка уже добавлена")
		}
		return result.Error
	}
	return context.Reply(fmt.Sprintf("Ставка добавлена.\nДата: <code>%v</code>.\nТекст: <code>%v</code>.", time.Unix(bet.Timestamp, 0).Format("02.01.2006"), bet.Text))
}
