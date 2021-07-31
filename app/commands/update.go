package commands

import (
	"fmt"
	"os/exec"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Restart bot on /restart
func Update(context telebot.Context) error {
	utils.Bot.Delete(context.Message())
	_, err := utils.Bot.Send(context.Sender(), "Starting go get...")
	if err != nil {
		return err
	}
	cmd, err := exec.Command("bash", "-c", "go get -u -v github.com/NexonSU/telegram-go-chatbot").CombinedOutput()
	if err != nil {
		return err
	}
	return context.Send(fmt.Sprintf("Update finished:\n<pre>%s</pre>", cmd))
}
