package commands

import (
	tele "gopkg.in/telebot.v3"
)

// Reply with stats link
func Stats(context tele.Context) error {
	return context.Reply("https://grafana.nexon.su/goto/gK6C_fw4R?orgId=1")
}
