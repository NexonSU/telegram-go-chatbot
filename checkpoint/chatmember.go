package checkpoint

import tele "gopkg.in/telebot.v3"

func ChatMemberUpdate(context tele.Context) error {
	Old := context.ChatMember().OldChatMember
	New := context.ChatMember().NewChatMember
	if Old.Role == "left" && New.Role == "member" {
		return UserJoin(context)
	}
	if (Old.Role == "member" || Old.Role == "restricted") &&
		(New.Role == "left" || New.Role == "kicked") {
		return UserLeft(context)
	}
	return nil
}
