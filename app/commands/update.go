package commands

import (
	"fmt"
	"os/exec"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Restart bot on /restart
func Update(context telebot.Context) error {
	var err error
	if !utils.IsAdminOrModer(context.Sender().Username) {
		if context.Chat().Username != utils.Config.Telegram.Chat {
			return err
		}
		return context.Reply(&telebot.Animation{File: telebot.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
	}
	utils.Bot.Delete(context.Message())
	_, err = utils.Bot.Send(context.Sender(), "Starting go get...")
	if err != nil {
		return err
	}
	cmd, err := exec.Command("bash", "-c", "go get -u -v github.com/NexonSU/telegram-go-chatbot").CombinedOutput()
	if err != nil {
		return err
	}
	_, err = utils.Bot.Send(context.Sender(), fmt.Sprintf("Update finished:\n<pre>%s</pre>", cmd))
	if err != nil {
		return err
	}
	return err
}
