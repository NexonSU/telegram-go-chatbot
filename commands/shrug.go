package commands

import tele "github.com/NexonSU/telebot"

//Send shrug in chat on /shrug
func Shrug(context tele.Context) error {
	return context.Send("¯\\_(ツ)_/¯")
}
