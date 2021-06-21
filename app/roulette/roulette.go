package roulette

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"strings"
	"time"
)

func RussianRoulette(busy map[string]bool, russianRouletteMessage *tb.Message, russianRouletteSelector tb.ReplyMarkup) func(*tb.Message) {
	return func(m *tb.Message) {
		if russianRouletteMessage == nil {
			russianRouletteMessage = m
			russianRouletteMessage.Unixtime = 0
		}
		if busy["bot_is_dead"] {
			if time.Now().Unix()-russianRouletteMessage.Time().Unix() > 3600 {
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
		if busy["russianroulettePending"] && !busy["russianrouletteInProgress"] && time.Now().Unix()-russianRouletteMessage.Time().Unix() > 60 {
			busy["russianroulette"] = false
			busy["russianroulettePending"] = false
			busy["russianrouletteInProgress"] = false
			_, err := utils.Bot.Edit(russianRouletteMessage, fmt.Sprintf("%v не пришел на дуэль.", utils.UserFullName(russianRouletteMessage.Entities[0].User)))
			if err != nil {
				utils.ErrorReporting(err, russianRouletteMessage)
				return
			}
		}
		if busy["russianrouletteInProgress"] && time.Now().Unix()-russianRouletteMessage.Time().Unix() > 120 {
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
		russianRouletteMessage, err = utils.Bot.Send(m.Chat, fmt.Sprintf("%v! %v вызывает тебя на дуэль!", utils.MentionUser(&target), utils.MentionUser(m.Sender)), &russianRouletteSelector)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		busy["russianroulettePending"] = true
	}
}
