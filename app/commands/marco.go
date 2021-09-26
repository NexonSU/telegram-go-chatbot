package commands

import (
	"github.com/NexonSU/telebot"
)

//Reply "Polo!" on "marco"
func Marco(context telebot.Context) error {
	return context.Reply("Polo!")
}
