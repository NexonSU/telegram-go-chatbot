package main

import (
	"github.com/NexonSU/telegram-go-chatbot/app/commands"
	"github.com/NexonSU/telegram-go-chatbot/app/middleware"
	"github.com/NexonSU/telegram-go-chatbot/app/roulette"
	"github.com/NexonSU/telegram-go-chatbot/app/services"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/NexonSU/telegram-go-chatbot/app/welcome"
	"gopkg.in/tucnak/telebot.v3"
)

func main() {
	utils.Bot.OnError = utils.ErrorReporting
	utils.Bot.Handle("/restart", commands.Restart, middleware.AdminLevel)
	utils.Bot.Handle("/update", commands.Update, middleware.AdminLevel)
	utils.Bot.Handle("/debug", commands.Debug, middleware.ModerLevel)
	utils.Bot.Handle("/say", commands.Say, middleware.ModerLevel)
	utils.Bot.Handle("/getid", commands.Getid, middleware.ModerLevel)
	utils.Bot.Handle("/kick", commands.Kick, middleware.ModerLevel)
	utils.Bot.Handle("/ban", commands.Ban, middleware.ModerLevel)
	utils.Bot.Handle("/unban", commands.Unban, middleware.ModerLevel)
	utils.Bot.Handle("/mute", commands.Mute, middleware.ModerLevel)
	utils.Bot.Handle("/unmute", commands.Unmute, middleware.ModerLevel)
	utils.Bot.Handle("/revive", commands.Revive, middleware.ModerLevel)
	utils.Bot.Handle("/resurrect", commands.Revive, middleware.ModerLevel)
	utils.Bot.Handle("/warn", commands.Warn, middleware.ModerLevel)
	utils.Bot.Handle("/pidordel", commands.Pidordel, middleware.ModerLevel)
	utils.Bot.Handle("/pidorlist", commands.Pidorlist, middleware.ModerLevel)
	utils.Bot.Handle("/kill", commands.Kill, middleware.ModerLevel)
	utils.Bot.Handle("/admin", commands.Admin, middleware.ChatLevel)
	utils.Bot.Handle("/get", commands.Get, middleware.ChatLevel)
	utils.Bot.Handle("/getall", commands.Getall, middleware.ChatLevel)
	utils.Bot.Handle("/set", commands.Set, middleware.ChatLevel)
	utils.Bot.Handle("/del", commands.Del, middleware.ChatLevel)
	utils.Bot.Handle("/shrug", commands.Shrug, middleware.ChatLevel)
	utils.Bot.Handle("/sed", commands.Sed, middleware.ChatLevel)
	utils.Bot.Handle("/ping", commands.Ping, middleware.ChatLevel)
	utils.Bot.Handle("/marco", commands.Marco, middleware.ChatLevel)
	utils.Bot.Handle("/cur", commands.Cur, middleware.ChatLevel)
	utils.Bot.Handle("/google", commands.Google, middleware.ChatLevel)
	utils.Bot.Handle("/me", commands.Me, middleware.ChatLevel)
	utils.Bot.Handle("/bonk", commands.Bonk, middleware.ChatLevel)
	utils.Bot.Handle("/hug", commands.Hug, middleware.ChatLevel)
	utils.Bot.Handle("/slap", commands.Slap, middleware.ChatLevel)
	utils.Bot.Handle("/releases", commands.Releases, middleware.ChatLevel)
	utils.Bot.Handle("/mywarns", commands.Mywarns, middleware.ChatLevel)
	utils.Bot.Handle("/pidorules", commands.Pidorules, middleware.ChatLevel)
	utils.Bot.Handle("/pidoreg", commands.Pidoreg, middleware.ChatLevel)
	utils.Bot.Handle("/pidorme", commands.Pidorme, middleware.ChatLevel)
	utils.Bot.Handle("/pidorall", commands.Pidorall, middleware.ChatLevel)
	utils.Bot.Handle("/pidorstats", commands.Pidorstats, middleware.ChatLevel)
	utils.Bot.Handle("/pidor", commands.Pidor, middleware.ChatLevel)
	utils.Bot.Handle("/blessing", commands.Blessing, middleware.ChatLevel)
	utils.Bot.Handle("/suicide", commands.Blessing, middleware.ChatLevel)
	utils.Bot.Handle("/duelstats", commands.Duelstats, middleware.ChatLevel)

	//Inline
	utils.Bot.Handle(telebot.OnQuery, services.OnInline)

	//Russian Roulette game
	utils.Bot.Handle("/russianroulette", roulette.Request)
	utils.Bot.Handle(&roulette.AcceptButton, roulette.Accept)
	utils.Bot.Handle(&roulette.DenyButton, roulette.Deny)

	//Repost channel post to chat
	utils.Bot.Handle(telebot.OnChannelPost, services.OnPost)

	//User join
	utils.Bot.Handle(telebot.OnUserJoined, welcome.OnJoin)
	utils.Bot.Handle(telebot.OnUserLeft, welcome.OnLeft)
	utils.Bot.Handle(&welcome.CorrectButton, welcome.OnClickCorrectButton)
	utils.Bot.Handle(&welcome.FirstWrongButton, welcome.OnClickWrongButton)
	utils.Bot.Handle(&welcome.SecondWrongButton, welcome.OnClickWrongButton)
	utils.Bot.Handle(&welcome.ThirdWrongButton, welcome.OnClickWrongButton)

	//Services
	go services.ZavtraStreamCheckService()
	go welcome.JoinMessageUpdateService()

	utils.Bot.Start()
}
