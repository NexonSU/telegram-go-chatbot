package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Send slap message on /slap
func Slap(m *tb.Message) {
	var action = "дал леща"
	var target tb.User
	ChatMember, err := utils.Bot.ChatMemberOf(m.Chat, m.Sender)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	if ChatMember.CanRestrictMembers || ChatMember.Role == "creator" {
		action = "дал отцовского леща"
	}
	target, _, err = utils.FindUserInMessage(*m)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err = utils.Bot.Send(m.Chat, fmt.Sprintf("👋 <b>%v</b> %v %v", utils.UserFullName(m.Sender), action, utils.MentionUser(&target)))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
