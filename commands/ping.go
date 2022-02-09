package commands

import "gopkg.in/telebot.v3"

//Reply "Pong!" on /ping
func Ping(context telebot.Context) error {
	return context.Reply("Pong!")
}
