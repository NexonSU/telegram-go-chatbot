package commands

import tele "gopkg.in/telebot.v3"

//Reply "Pong!" on /ping
func Ping(context tele.Context) error {
	return context.Reply("Pong!")
}
