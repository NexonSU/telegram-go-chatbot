package commands

import tele "github.com/NexonSU/telebot"

//Reply "Polo!" on "marco"
func Marco(context tele.Context) error {
	return context.Reply("Polo!")
}
