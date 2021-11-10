package commands

import (
	"gopkg.in/tucnak/telebot.v3"
)

//Reply "Polo!" on "marco"
func Marco(context telebot.Context) error {
	return context.Reply("Polo!")
}
