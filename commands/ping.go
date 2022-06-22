package commands

import tele "github.com/NexonSU/telebot"

//Reply "Pong!" on /ping
func Ping(context tele.Context) error {
	return context.Reply("Pong!")
}
