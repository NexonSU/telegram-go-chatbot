package commands

import (
	"encoding/json"
	"fmt"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

//Return message on /debug command
func Debug(context tele.Context) error {
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
