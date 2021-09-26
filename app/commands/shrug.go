package commands

import (
	"github.com/NexonSU/telebot"
)

//Send shrug in chat on /shrug
func Shrug(context telebot.Context) error {
	return context.Send("¯\\_(ツ)_/¯")
}
