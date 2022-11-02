package main

import (
	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/checkpoint"
	"github.com/NexonSU/telegram-go-chatbot/commands"
	"github.com/NexonSU/telegram-go-chatbot/duel"
	"github.com/NexonSU/telegram-go-chatbot/pidor"
	"github.com/NexonSU/telegram-go-chatbot/stats"
	"github.com/NexonSU/telegram-go-chatbot/utils"
)

func main() {
	//Initializing Bot
	utils.BotInit()

	//limiting bot use
	admin := utils.Bot.Group()
	admin.Use(utils.Whitelist(append(utils.Config.Admins, utils.Config.SysAdmin)...))
	moder := utils.Bot.Group()
	moder.Use(utils.Whitelist(append(append(utils.Config.Admins, utils.Config.Moders...), utils.Config.SysAdmin)...))
	chats := utils.Bot.Group()
	chats.Use(utils.Whitelist(append(append(utils.Config.Admins, utils.Config.Moders...), utils.Config.SysAdmin, utils.Config.Chat, utils.Config.ReserveChat)...))
	chato := utils.Bot.Group()
	chato.Use(utils.Whitelist(utils.Config.ReserveChat))
	comms := utils.Bot.Group()
	comms.Use(utils.Whitelist(utils.Config.CommentChat))
	chann := utils.Bot.Group()
	chann.Use(utils.Whitelist(utils.Config.Channel))

	//commands
	admin.Handle("/restart", commands.Restart)
	moder.Handle("/debug", commands.Debug)
	moder.Handle("/say", commands.Say)
	moder.Handle("/getid", commands.Getid)
	moder.Handle("/kick", commands.Kick)
	moder.Handle("/ban", commands.Ban)
	moder.Handle("/unban", commands.Unban)
	moder.Handle("/mute", commands.Mute)
	moder.Handle("/unmute", commands.Unmute)
	moder.Handle("/revive", commands.Revive)
	moder.Handle("/resurrect", commands.Revive)
	moder.Handle("/warn", commands.Warn)
	moder.Handle("/kill", commands.Kill)
	moder.Handle("/bless", commands.Kill)
	moder.Handle("/addnope", commands.AddNope)
	moder.Handle("/setgetowner", commands.SetGetOwner)
	moder.Handle("/addantispam", checkpoint.AddAntispam)
	moder.Handle("/listantispam", checkpoint.ListAntispam)
	moder.Handle("/delantispam", checkpoint.DelAntispam)
	moder.Handle("/getspamchance", checkpoint.CommandGetSpamChance)
	moder.Handle("/convert", commands.Convert)
	chats.Handle("/admin", commands.Admin)
	chats.Handle("/get", commands.Get)
	chats.Handle("/getall", commands.Getall)
	chats.Handle("/set", commands.Set)
	chats.Handle("/del", commands.Del)
	chats.Handle("/shrug", commands.Shrug)
	chats.Handle("/sed", commands.Sed)
	chats.Handle("/ping", commands.Ping)
	chats.Handle("/marco", commands.Marco)
	chats.Handle("/savetopm", commands.SaveToPM)
	chats.Handle("/topm", commands.SaveToPM)
	chats.Handle("/giveme", commands.SaveToPM)
	chats.Handle("/cur", commands.Cur)
	chats.Handle("/google", commands.Google)
	chats.Handle("/me", commands.Me)
	chats.Handle("/bonk", commands.Bonk)
	chats.Handle("/hug", commands.Hug)
	chats.Handle("/slap", commands.Slap)
	chats.Handle("/releases", commands.Releases)
	chats.Handle("/mywarns", commands.Mywarns)
	chats.Handle("/blessing", commands.Blessing)
	chats.Handle("/suicide", commands.Blessing)
	chats.Handle("/isekai", commands.Blessing)
	chats.Handle("/testrandom", commands.TestRandom)
	moder.Handle("/anekdot", commands.Anekdot)
	moder.Handle("/bashorg", commands.Bashorg)
	moder.Handle("/meow", commands.Meow)

	//stats commands
	chats.Handle("/stats", stats.Stats)
	admin.Handle("/removeword", stats.RemoveWord)

	//pidor of the day commands
	chats.Handle("/pidor", pidor.Pidor)
	chats.Handle("/pidoreg", pidor.Pidoreg)
	chats.Handle("/pidorules", pidor.Pidorules)
	chats.Handle("/pidorme", pidor.Pidorme)
	chats.Handle("/pidorall", pidor.Pidorall)
	chats.Handle("/pidorstats", pidor.Pidorstats)
	moder.Handle("/pidordel", pidor.Pidordel)
	moder.Handle("/pidorlist", pidor.Pidorlist)

	//duel commands and buttons
	chats.Handle("/russianroulette", duel.Request)
	chats.Handle("/duel", duel.Request)
	chats.Handle("/duelstats", duel.Duelstats)
	chats.Handle(&duel.AcceptButton, duel.Accept)
	chats.Handle(&duel.DenyButton, duel.Deny)

	//repost channel post to chat
	chann.Handle(tele.OnChannelPost, utils.ForwardPost)

	//spam filter in comment chat
	comms.Handle(tele.OnText, checkpoint.SpamFilter)
	comms.Handle(tele.OnSticker, checkpoint.SpamFilter)

	//user entry filter
	chato.Handle(tele.OnChatMember, checkpoint.ChatMemberUpdate)
	chato.Handle(tele.OnUserJoined, utils.Remove)
	chato.Handle(tele.OnUserLeft, utils.Remove)

	//inline
	utils.Bot.Handle(tele.OnQuery, commands.GetInline)

	utils.Bot.Start()
}
