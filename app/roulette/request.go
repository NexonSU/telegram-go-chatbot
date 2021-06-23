package roulette

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"time"
)

var Message *tb.Message
var Selector = tb.ReplyMarkup{}
var AcceptButton = Selector.Data("👍 Принять вызов", "russianroulette_accept")
var DenyButton = Selector.Data("👎 Бежать с позором", "russianroulette_deny")
var busy = make(map[string]bool)

func Request(m *tb.Message) {
	if m.Chat.Username != utils.Config.Telegram.Chat {
		return
	}
	if Message == nil {
		Message = m
		Message.Unixtime = 0
	}
	if busy["bot_is_dead"] {
		if time.Now().Unix()-Message.Time().Unix() > 3600 {
			busy["bot_is_dead"] = false
		} else {
			_, err := utils.Bot.Reply(m, "Я не могу провести игру, т.к. я немного умер. Зайдите позже.")
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			return
		}
	}
	if busy["russianroulettePending"] && !busy["russianrouletteInProgress"] && time.Now().Unix()-Message.Time().Unix() > 60 {
		busy["russianroulette"] = false
		busy["russianroulettePending"] = false
		busy["russianrouletteInProgress"] = false
		_, err := utils.Bot.Edit(Message, fmt.Sprintf("%v не пришел на дуэль.", utils.UserFullName(Message.Entities[0].User)))
		if err != nil {
			utils.ErrorReporting(err, Message)
			return
		}
	}
	if busy["russianrouletteInProgress"] && time.Now().Unix()-Message.Time().Unix() > 120 {
		busy["russianroulette"] = false
		busy["russianroulettePending"] = false
		busy["russianrouletteInProgress"] = false
	}
	if busy["russianroulette"] || busy["russianroulettePending"] || busy["russianrouletteInProgress"] {
		_, err := utils.Bot.Reply(m, "Команда занята. Попробуйте позже.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	busy["russianroulette"] = true
	defer func() { busy["russianroulette"] = false }()
	var text = strings.Split(m.Text, " ")
	if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
		_, err := utils.Bot.Reply(m, "Пример использования: <code>/russianroulette {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/russianroulette</code>")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	target, _, err := utils.FindUserInMessage(*m)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	if target.ID == m.Sender.ID {
		_, err := utils.Bot.Reply(m, "Как ты себе это представляешь? Нет, нельзя вызвать на дуэль самого себя.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	if target.IsBot {
		_, err := utils.Bot.Reply(m, "Бота нельзя вызвать на дуэль.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	ChatMember, err := utils.Bot.ChatMemberOf(m.Chat, &target)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	log.Println(ChatMember)
	if false {
		_, err := utils.Bot.Reply(m, "Нельзя вызвать на дуэль мертвеца.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	err = utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	Selector.Inline(
		Selector.Row(AcceptButton, DenyButton),
	)
	Message, err = utils.Bot.Send(m.Chat, fmt.Sprintf("%v! %v вызывает тебя на дуэль!", utils.MentionUser(&target), utils.MentionUser(m.Sender)), &Selector)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	busy["russianroulettePending"] = true
}
