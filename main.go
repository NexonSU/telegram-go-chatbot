package main

import (
	"github.com/NexonSU/telegram-go-chatbot/app/commands"
	"github.com/NexonSU/telegram-go-chatbot/app/roulette"
	"github.com/NexonSU/telegram-go-chatbot/app/services"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/NexonSU/telegram-go-chatbot/app/welcome"
	"gopkg.in/tucnak/telebot.v3"
)

func main() {
	//utils.Bot.OnError = utils.ErrorReporting
	utils.Bot.Handle("/admin", commands.Admin)
	utils.Bot.Handle("/debug", commands.Debug)
	utils.Bot.Handle("/get", commands.Get)
	utils.Bot.Handle("/getall", commands.Getall)
	utils.Bot.Handle("/set", commands.Set)
	utils.Bot.Handle("/del", commands.Del)
	utils.Bot.Handle("/say", commands.Say)
	utils.Bot.Handle("/shrug", commands.Shrug)
	utils.Bot.Handle("/sed", commands.Sed)
	utils.Bot.Handle("/getid", commands.Getid)
	utils.Bot.Handle("/ping", commands.Ping)
	utils.Bot.Handle("/marco", commands.Marco)
	utils.Bot.Handle("/cur", commands.Cur)
	utils.Bot.Handle("/google", commands.Google)
	utils.Bot.Handle("/kick", commands.Kick)
	utils.Bot.Handle("/ban", commands.Ban)
	utils.Bot.Handle("/unban", commands.Unban)
	utils.Bot.Handle("/mute", commands.Mute)
	utils.Bot.Handle("/unmute", commands.Unmute)
	utils.Bot.Handle("/revive", commands.Revive)
	utils.Bot.Handle("/resurrect", commands.Revive)
	utils.Bot.Handle("/me", commands.Me)
	utils.Bot.Handle("/bonk", commands.Bonk)
	utils.Bot.Handle("/hug", commands.Hug)
	utils.Bot.Handle("/slap", commands.Slap)
	utils.Bot.Handle("/releases", commands.Releases)
	utils.Bot.Handle("/warn", commands.Warn)
	utils.Bot.Handle("/mywarns", commands.Mywarns)
	utils.Bot.Handle("/pidorules", commands.Pidorules)
	utils.Bot.Handle("/pidoreg", commands.Pidoreg)
	utils.Bot.Handle("/pidorme", commands.Pidorme)
	utils.Bot.Handle("/pidordel", commands.Pidordel)
	utils.Bot.Handle("/pidorlist", commands.Pidorlist)
	utils.Bot.Handle("/pidorall", commands.Pidorall)
	utils.Bot.Handle("/pidorstats", commands.Pidorstats)
	utils.Bot.Handle("/pidor", commands.Pidor)
	utils.Bot.Handle("/blessing", commands.Blessing)
	utils.Bot.Handle("/suicide", commands.Blessing)
	utils.Bot.Handle("/kill", commands.Kill)
	utils.Bot.Handle("/duelstats", commands.Duelstats)
	utils.Bot.Handle("/restart", commands.Restart)
	utils.Bot.Handle("/update", commands.Update)

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
