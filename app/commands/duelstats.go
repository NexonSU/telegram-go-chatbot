package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send user utils.Duelist stats on /duelstats
func Duelstats(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(m.Sender.Username) {
		return
	}
	var duelist utils.Duelist
	result := utils.DB.Model(utils.Duelist{}).Where(m.Sender.ID).First(&duelist)
	if result.RowsAffected == 0 {
		_, err := utils.Bot.Reply(m, "У тебя нет статистики.")
		if err != nil {
			utils.ErrorReporting(err, m)
		}
		return
	}
	_, err := utils.Bot.Reply(m, fmt.Sprintf("Побед: %v\nСмертей: %v", duelist.Kills, duelist.Deaths))
	if err != nil {
		utils.ErrorReporting(err, m)
	}
}
