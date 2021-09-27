package main

import (
	"fmt"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/checkpoint"
	"github.com/NexonSU/telegram-go-chatbot/app/commands"
	"github.com/NexonSU/telegram-go-chatbot/app/duel"
	"github.com/NexonSU/telegram-go-chatbot/app/pidor"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

func main() {
	utils.Bot.OnError = utils.ErrorReporting
	utils.Bot.Handle("/restart", commands.Restart, utils.AdminLevel)
	utils.Bot.Handle("/update", commands.Update, utils.AdminLevel)
	utils.Bot.Handle("/debug", commands.Debug, utils.ModerLevel)
	utils.Bot.Handle("/say", commands.Say, utils.ModerLevel)
	utils.Bot.Handle("/getid", commands.Getid, utils.ModerLevel)
	utils.Bot.Handle("/kick", commands.Kick, utils.ModerLevel)
	utils.Bot.Handle("/ban", commands.Ban, utils.ModerLevel)
	utils.Bot.Handle("/unban", commands.Unban, utils.ModerLevel)
	utils.Bot.Handle("/mute", commands.Mute, utils.ModerLevel)
	utils.Bot.Handle("/unmute", commands.Unmute, utils.ModerLevel)
	utils.Bot.Handle("/revive", commands.Revive, utils.ModerLevel)
	utils.Bot.Handle("/resurrect", commands.Revive, utils.ModerLevel)
	utils.Bot.Handle("/warn", commands.Warn, utils.ModerLevel)
	utils.Bot.Handle("/pidordel", pidor.Pidordel, utils.ModerLevel)
	utils.Bot.Handle("/pidorlist", pidor.Pidorlist, utils.ModerLevel)
	utils.Bot.Handle("/kill", commands.Kill, utils.ModerLevel)
	utils.Bot.Handle("/addnope", commands.AddNope, utils.ModerLevel)
	utils.Bot.Handle("/setgetowner", commands.SetGetOwner, utils.ModerLevel)
	utils.Bot.Handle("/addantispam", checkpoint.AddAntispam, utils.ModerLevel)
	utils.Bot.Handle("/listantispam", checkpoint.ListAntispam, utils.ModerLevel)
	utils.Bot.Handle("/delantispam", checkpoint.DelAntispam, utils.ModerLevel)
	utils.Bot.Handle("/getspamchance", checkpoint.CommandGetSpamChance, utils.ModerLevel)
	utils.Bot.Handle("/convert", commands.Convert, utils.ModerLevel)
	utils.Bot.Handle("/admin", commands.Admin, utils.ChatLevel)
	utils.Bot.Handle("/get", commands.Get, utils.ChatLevel)
	utils.Bot.Handle(telebot.OnQuery, commands.GetInline)
	utils.Bot.Handle("/getall", commands.Getall, utils.ChatLevel)
	utils.Bot.Handle("/set", commands.Set, utils.GetFilterCreator)
	utils.Bot.Handle("/del", commands.Del, utils.GetFilterCreator)
	utils.Bot.Handle("/shrug", commands.Shrug, utils.ChatLevel)
	utils.Bot.Handle("/sed", commands.Sed, utils.ChatLevel)
	utils.Bot.Handle("/ping", commands.Ping, utils.ChatLevel)
	utils.Bot.Handle("/marco", commands.Marco, utils.ChatLevel)
	utils.Bot.Handle("/cur", commands.Cur, utils.ChatLevel)
	utils.Bot.Handle("/google", commands.Google, utils.ChatLevel)
	utils.Bot.Handle("/me", commands.Me, utils.ChatLevel)
	utils.Bot.Handle("/bonk", commands.Bonk, utils.ChatLevel)
	utils.Bot.Handle("/hug", commands.Hug, utils.ChatLevel)
	utils.Bot.Handle("/slap", commands.Slap, utils.ChatLevel)
	utils.Bot.Handle("/releases", commands.Releases, utils.ChatLevel)
	utils.Bot.Handle("/mywarns", commands.Mywarns, utils.ChatLevel)
	utils.Bot.Handle("/pidorules", pidor.Pidorules, utils.ChatLevel)
	utils.Bot.Handle("/pidoreg", pidor.Pidoreg, utils.ChatLevel)
	utils.Bot.Handle("/pidorme", pidor.Pidorme, utils.ChatLevel)
	utils.Bot.Handle("/pidorall", pidor.Pidorall, utils.ChatLevel)
	utils.Bot.Handle("/pidorstats", pidor.Pidorstats, utils.ChatLevel)
	utils.Bot.Handle("/pidor", pidor.Pidor, utils.ChatLevel)
	utils.Bot.Handle("/blessing", commands.Blessing, utils.ChatLevel)
	utils.Bot.Handle("/suicide", commands.Blessing, utils.ChatLevel)

	//Russian Roulette duels
	utils.Bot.Handle("/russianroulette", duel.Request, utils.ChatOnly)
	utils.Bot.Handle("/duel", duel.Request, utils.ChatOnly)
	utils.Bot.Handle("/duelstats", duel.Duelstats, utils.ChatLevel)
	utils.Bot.Handle(&duel.AcceptButton, duel.Accept, utils.ChatOnly)
	utils.Bot.Handle(&duel.DenyButton, duel.Deny, utils.ChatOnly)

	//Repost channel post to chat
	utils.Bot.Handle(telebot.OnChannelPost, utils.Repost, utils.ChannelOnly)
	utils.Bot.Handle(telebot.OnEditedChannelPost, utils.EditRepost, utils.ChannelOnly)

	//Filter messages in comment chat
	utils.Bot.Handle(telebot.OnText, checkpoint.SpamFilter, utils.CommentChatOnly)
	utils.Bot.Handle(telebot.OnSticker, checkpoint.SpamFilter, utils.CommentChatOnly)

	//User join
	utils.Bot.Handle(telebot.OnChatMember, checkpoint.ChatMemberUpdate, utils.ChatOnly)
	utils.Bot.Handle(telebot.OnUserJoined, utils.Remove, utils.ChatOnly)
	utils.Bot.Handle(telebot.OnUserLeft, utils.Remove, utils.ChatOnly)
	utils.Bot.Handle(telebot.OnCallback, checkpoint.ButtonCallback, utils.ChatOnly)

	utils.Bot.Send(telebot.ChatID(utils.Config.SysAdmin), fmt.Sprintf("%v has finished starting up.", utils.Bot.Me.MentionHTML()))

	utils.Bot.Start()
}
