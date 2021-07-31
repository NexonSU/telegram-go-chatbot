package roulette

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

var Message *telebot.Message
var Selector = telebot.ReplyMarkup{}
var AcceptButton = Selector.Data("👍 Принять вызов", "russianroulette_accept")
var DenyButton = Selector.Data("👎 Бежать с позором", "russianroulette_deny")
var busy = make(map[string]bool)

func Request(context telebot.Context) error {
	if Message == nil {
		Message = context.Message()
		Message.Unixtime = 0
	}
	if busy["bot_is_dead"] {
		if time.Now().Unix()-Message.Time().Unix() > 3600 {
			busy["bot_is_dead"] = false
		} else {
			return context.Reply("Я не могу провести игру, т.к. я немного умер. Зайдите позже.")
		}
	}
	if busy["russianroulettePending"] && !busy["russianrouletteInProgress"] && time.Now().Unix()-Message.Time().Unix() > 60 {
		busy["russianroulette"] = false
		busy["russianroulettePending"] = false
		busy["russianrouletteInProgress"] = false
		return context.Edit(fmt.Sprintf("%v не пришел на дуэль.", utils.UserFullName(Message.Entities[0].User)))
	}
	if busy["russianrouletteInProgress"] && time.Now().Unix()-Message.Time().Unix() > 120 {
		busy["russianroulette"] = false
		busy["russianroulettePending"] = false
		busy["russianrouletteInProgress"] = false
	}
	if busy["russianroulette"] || busy["russianroulettePending"] || busy["russianrouletteInProgress"] {
		return context.Reply("Команда занята. Попробуйте позже.")
	}
	busy["russianroulette"] = true
	defer func() { busy["russianroulette"] = false }()
	var text = strings.Split(context.Text(), " ")
	if (context.Message().ReplyTo == nil && len(text) != 2) || (context.Message().ReplyTo != nil && len(text) != 1) {
		return context.Reply("Пример использования: <code>/russianroulette {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/russianroulette</code>")
	}
	target, _, err := utils.FindUserInMessage(context)
	if err != nil {
		return context.Reply(fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
	}
	if target.ID == context.Sender().ID {
		return context.Reply("Как ты себе это представляешь? Нет, нельзя вызвать на дуэль самого себя.")
	}
	if target.IsBot {
		return context.Reply("Бота нельзя вызвать на дуэль.")
	}
	ChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &target)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
	}
	log.Println(ChatMember)
	if false {
		err := context.Reply("Нельзя вызвать на дуэль мертвеца.")
		if err != nil {
			return err
		}
		return err
	}
	err = utils.Bot.Delete(context.Message())
	if err != nil {
		return err
	}
	Selector.Inline(
		Selector.Row(AcceptButton, DenyButton),
	)
	Message, err = utils.Bot.Send(context.Chat(), fmt.Sprintf("%v! %v вызывает тебя на дуэль!", utils.MentionUser(&target), utils.MentionUser(context.Sender())), &Selector)
	if err != nil {
		return err
	}
	busy["russianroulettePending"] = true
	return err
}
