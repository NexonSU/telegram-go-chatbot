package commands

import (
	"encoding/json"
	"fmt"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Return message on /debug command
func Debug(context telebot.Context) error {
	err := utils.Bot.Delete(context.Message())
	if err != nil {
		return err
	}
	var message = context.Message()
	if context.Message().ReplyTo != nil {
		message = context.Message().ReplyTo
	}
	MarshalledMessage, _ := json.MarshalIndent(message, "", "    ")
	_, err = utils.Bot.Send(context.Sender(), fmt.Sprintf("<pre>%v</pre>", string(MarshalledMessage)))
	return err
}
