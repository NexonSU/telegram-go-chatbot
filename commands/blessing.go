package commands

import (
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

var firstSuicide int64
var lastSuicide int64
var burst int
var lastVideoSent int64

// Kill user on /blessing, /suicide
func Blessing(context tele.Context) error {
	// prt will replace fmt package to format text according plurals defined in utils package
	// If no plural rule matched it will be ignored and processed as usual formatting
	prt := message.NewPrinter(language.Russian)

	victim := context.Sender()
	ricochetText := ""

	err := context.Delete()
	if err != nil {
		return err
	}
	ChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), context.Sender())
	if err != nil {
		return err
	}
	if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
		var ricochetVictim *tele.ChatMember
		var lastMessage utils.Message
		for i := 1; i < 100; i++ {
			lastMessage = utils.Message{}
			result := utils.DB.Where(utils.Message{ChatID: context.Chat().ID}).Order("id desc").Offset(i).Last(&lastMessage)
			if result.Error != nil {
				continue
			}
			ricochetVictim = &tele.ChatMember{}
			ricochetVictim, err = utils.Bot.ChatMemberOf(context.Chat(), &tele.User{ID: lastMessage.UserID})
			if err != nil {
				continue
			}
			if ricochetVictim.Role == "member" {
				victim = ricochetVictim.User
				ChatMember = ricochetVictim
				ricochetText = prt.Sprintf("Пуля отскакивает от головы %v и летит в голову %v.\n", utils.MentionUser(context.Sender()), utils.MentionUser(victim))
			}
		}
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
	duration += 10
	prependText := ""
	additionalChance := int(time.Now().Unix() - lastSuicide)
	if additionalChance > 3600 {
		additionalChance = 3600
	}
	additionalChance = (3600 - additionalChance) / 360
	if context.Sender().IsPremium {
		duration = duration * 2
		prependText += "премиально "
	}
	if utils.RandInt(0, 100) >= 90-additionalChance {
		duration = duration * 10
		prependText += "критически "
	}
	if duration >= 1400 && duration <= 1500 {
		duration = 1488
	}
	if ricochetText != "" {
		duration = 1
	}
	ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(60*duration)).Unix()
	err = utils.Bot.Restrict(context.Chat(), ChatMember)
	if err != nil {
		return err
	}
	burst++
	if time.Now().Unix() > firstSuicide+120 {
		firstSuicide = time.Now().Unix()
		burst = 1
	}
	lastSuicide = time.Now().Unix()
	if burst > 3 && time.Now().Unix() > lastVideoSent+3600 {
		lastVideoSent = time.Now().Unix()
		return context.Send(&tele.Video{
			File: tele.File{
				FileID: "BAACAgIAAx0CReJGYgABAlMuYnagTilFaB8ke8Rw-dYLbfJ6iF8AAicYAAIlxrlLY9ah2fUtR40kBA",
			},
			Caption: prt.Sprintf("<code>%v💥 %v %v%v.\nРеспавн через %d мин.</code>", ricochetText, utils.UserFullName(victim), prependText, utils.GetBless(), duration),
		})
	} else {
		return context.Send(prt.Sprintf("<code>%v💥 %v %v%v.\nРеспавн через %d мин.</code>", ricochetText, utils.UserFullName(victim), prependText, utils.GetBless(), duration))
	}
}
