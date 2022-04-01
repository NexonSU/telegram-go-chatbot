package commands

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

//Kill user on /blessing, /suicide
func Blessing(context tele.Context) error {
	err := context.Delete()
	if err != nil {
		return err
	}
	ChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), context.Sender())
	if err != nil {
		return err
	}
	if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
		return context.Send(fmt.Sprintf("<code>👻 %v возродился у костра.</code>", utils.UserFullName(context.Sender())))
	}
	var duelist utils.Duelist
	result := utils.DB.Model(utils.Duelist{}).Where(context.Sender().ID).First(&duelist)
	if result.RowsAffected == 0 {
		duelist.UserID = context.Sender().ID
		duelist.Kills = 0
		duelist.Deaths = 0
	}
	duelist.Deaths++
	result = utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&duelist)
	if result.Error != nil {
		return result.Error
	}
	duration := utils.RandInt(1, duelist.Deaths+1)
	prependText := ""
	if utils.RandInt(0, 100) >= 98 {
		duration = duration * 10
		prependText = "критически "
	}
	if duration > 600 {
		duration = 600
	}
	ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(60*duration)).Unix()
	err = utils.Bot.Restrict(context.Chat(), ChatMember)
	if err != nil {
		return err
	}
	reason := []string{
		"выбрал лёгкий путь",
		"сыграл в ящик",
		"слил своё HP до нуля",
		"приказал долго жить",
		"покинул этот скорбный мир",
		"пагиб",
		"разбежавшись прыгнул со скалы",
		"разогнал RTX 4090 Ti",
		"принял ислам",
		"пьёт чай и кушоет конфеты, никакова суецыда",
	}
	return context.Send(fmt.Sprintf("<code>💥 %v %v%v.\nРеспавн через %v мин.</code>", utils.UserFullName(context.Sender()), prependText, reason[rand.Intn(len(reason))], duration))
}
