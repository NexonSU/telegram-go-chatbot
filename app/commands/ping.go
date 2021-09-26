package commands

import (
	"github.com/NexonSU/telebot"
)

//Reply "Pong!" on /ping
func Ping(context telebot.Context) error {
	return context.Reply("Pong!")
}
