package commands

import (
	"fmt"
	"os/exec"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Restart bot on /restart
func Update(m *tb.Message) {
	if !utils.IsAdminOrModer(m.Sender.Username) {
		if m.Chat.Username != utils.Config.Telegram.Chat {
			return
		}
		_, err := utils.Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	utils.Bot.Delete(m)
	out, err := exec.Command("go", "get -u -v github.com/NexonSU/telegram-go-chatbot").Output()
	if err != nil {
		utils.ErrorReporting(err, m)
	}
	_, err = utils.Bot.Send(m.Sender, fmt.Sprintf("Update finished:\n<pre>%v</pre>", out))
	if err != nil {
		utils.ErrorReporting(err, m)
	}
}
