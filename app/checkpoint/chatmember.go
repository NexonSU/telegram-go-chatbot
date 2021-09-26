package checkpoint

import (
	"github.com/NexonSU/telebot"
)

func ChatMemberUpdate(context telebot.Context) error {
	if Border.Chat == nil {
		Border.Chat = context.Chat()
	}
	Old := context.ChatMember().OldChatMember
	New := context.ChatMember().NewChatMember
	if Old.Role == "left" && New.Role == "member" {
		return UserJoin(context)
	}
	if Old.Role == "member" && New.Role == "left" || Old.Role == "restricted" && New.Role == "left" {
		return UserLeft(context)
	}
	return nil
}
