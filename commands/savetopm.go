package commands

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

// Resend post on user request
func SaveToPM(context tele.Context) error {
	if context.Message() == nil || context.Message().ReplyTo == nil {
		return utils.SendAndRemove("Пример использования:\n/topm в ответ на какое-либо сообщение\nБот должен быть запущен и разблокирован в личке.", context)
	}
	link := fmt.Sprintf("https://t.me/c/%v/%v", strings.TrimLeft(strings.TrimLeft(strconv.Itoa(int(context.Chat().ID)), "-1"), "0"), context.Message().ReplyTo.ID)
	var err error
	if context.Message().ReplyTo.Media() != nil {
		msgID, chatID := context.Message().ReplyTo.MessageSig()
		if context.Message().ReplyTo.Caption != "" {
			link = fmt.Sprintf("%v\n\n%v", context.Message().ReplyTo.Caption, link)
		}
		params := map[string]string{
			"chat_id":      context.Sender().Recipient(),
			"from_chat_id": strconv.FormatInt(chatID, 10),
			"message_id":   msgID,
			"caption":      link,
		}
		_, err = utils.Bot.Raw("copyMessage", params)
	} else {
		_, err = utils.Bot.Send(context.Sender(), html.EscapeString(fmt.Sprintf("%v\n\n%v", context.Message().ReplyTo.Text, link)))
	}
	if err != nil {
		return err
	}
	return context.Delete()
}
