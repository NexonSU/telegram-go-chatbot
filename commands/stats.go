package commands

import (
	tele "gopkg.in/telebot.v3"
)

// Reply with stats link
func Stats(context tele.Context) error {
	return context.Reply("Статистика чата", &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{{
			{
				Text:   "Открыть",
				WebApp: &tele.WebApp{URL: "https://grafana.nexon.su/d/aef7a25c-3824-4046-8ed3-53ccb5850c9d/zavtrachat?orgId=1&kiosk"},
			},
		}},
	})
}
