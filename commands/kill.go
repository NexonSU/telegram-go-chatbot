package commands

import (
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

// Kill user on /kill
func Kill(context tele.Context) error {
	// prt will replace fmt package to format text according plurals defined in utils package
	// If no plural rule matched it will be ignored and processed as usual formatting
	prt := message.NewPrinter(language.Russian)

	command := strings.Split(strings.Split(context.Text(), "@")[0], " ")[0]
	if (context.Message().ReplyTo == nil && len(context.Args()) != 1) || (context.Message().ReplyTo != nil && len(context.Args()) != 0) {
		return utils.SendAndRemove(prt.Sprintf("Пример использования: <code>%v {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>%v</code>", command, command), context)
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return err
	}
	ChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return err
	}
	victimText := ""
	if ChatMember.Role == "administrator" || ChatMember.Role == "creator" || context.Sender().ID == 825209730 {
		var victim *tele.ChatMember
		var userID int64
		rows, err := utils.DB.Model(&utils.Stats{}).Where("stat_type = 3").Order("last_update desc").Select("context_id").Limit(100).Rows()
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			rows.Scan(&userID)
			victim, err = utils.Bot.ChatMemberOf(context.Chat(), &tele.User{ID: userID})
			if err != nil {
				continue
			}
			if victim.Role == "member" {
				ChatMember = victim
				victimText = prt.Sprintf("Пуля отскакивает от головы %v и летит в голову %v.\n", utils.MentionUser(&target), utils.MentionUser(victim.User))
				target = *ChatMember.User
				rows.Close()
				break
			}
		}
	} else {
		if context.Message().ReplyTo != nil {
			utils.Bot.Delete(context.Message().ReplyTo)
		}
	}
	var duelist utils.Duelist
	result := utils.DB.Model(utils.Duelist{}).Where(target.ID).First(&duelist)
	if result.RowsAffected == 0 {
		duelist.UserID = target.ID
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
	if utils.RandInt(0, 100) >= 90 {
		duration = duration * 10
		prependText = "критически "
		if command == "/bless" {
			prependText = "очень "
		}
	}
	if victimText != "" {
		duration = 1
	}
	ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(60*duration)).Unix()
	err = utils.Bot.Restrict(context.Chat(), ChatMember)
	if err != nil {
		return err
	}
	text := prt.Sprintf("💥 %v %vпристрелил %v.\n%v отправился на респавн на %d мин.", utils.UserFullName(context.Sender()), prependText, utils.UserFullName(&target), utils.UserFullName(&target), duration)
	if command == "/bless" {
		text = prt.Sprintf("🤫 %v %vпопросил %v помолчать %d минут.", utils.UserFullName(context.Sender()), prependText, utils.UserFullName(&target), duration)
	}
	if command == "/bite" {
		text = prt.Sprintf("😼 %v %vсделал кусь %v.\n%v отправился на респавн на %d мин.", utils.UserFullName(context.Sender()), prependText, utils.UserFullName(&target), utils.UserFullName(&target), duration)
	}
	if victimText != "" {
		text = prt.Sprintf("💥 %v\n%v отправился на респавн на %d мин.", victimText, utils.UserFullName(&target), duration)
	}
	return context.Send(text)
}
