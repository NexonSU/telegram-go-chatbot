package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/gorm/clause"
	"time"
)

func Suicide(m *tb.Message) {
	err := utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	ChatMember, err := utils.Bot.ChatMemberOf(m.Chat, m.Sender)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("<code>👻 %v возродился у костра.</code>", utils.UserFullName(m.Sender)))
		if err != nil {
			utils.ErrorReporting(err, m)
		}
		return
	}
	var duelist utils.Duelist
	result := utils.DB.Model(utils.Duelist{}).Where(m.Sender.ID).First(&duelist)
	if result.RowsAffected == 0 {
		duelist.UserID = m.Sender.ID
		duelist.Kills = 0
		duelist.Deaths = 0
	}
	duelist.Deaths++
	result = utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(duelist)
	if result.Error != nil {
		utils.ErrorReporting(result.Error, m)
		return
	}
	ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(600*duelist.Deaths)).Unix()
	err = utils.Bot.Restrict(m.Chat, ChatMember)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	_, err = utils.Bot.Send(m.Chat, fmt.Sprintf("<code>💥 %v выбрал лёгкий путь.\nРеспавн через %v0 минут.</code>", utils.UserFullName(m.Sender), duelist.Deaths))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
