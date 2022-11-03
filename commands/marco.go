package commands

import tele "gopkg.in/telebot.v3"

//Reply "Polo!" on "marco"
func Marco(context tele.Context) error {
	return context.Reply("Polo!")
}
