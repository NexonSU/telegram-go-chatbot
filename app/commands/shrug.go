package commands

import (
	"gopkg.in/tucnak/telebot.v3"
)

//Send shrug in chat on /shrug
func Shrug(context telebot.Context) error {
	return context.Send("¯\\_(ツ)_/¯")
}
