package main

import (
	"log"
	"sort"

	"github.com/NexonSU/telegram-go-chatbot/bets"
	"github.com/NexonSU/telegram-go-chatbot/commands"
	"github.com/NexonSU/telegram-go-chatbot/duel"
	"github.com/NexonSU/telegram-go-chatbot/pidor"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
)

type commandList struct {
	command    tele.Command
	function   tele.HandlerFunc
	middleware tele.MiddlewareFunc
}

func main() {
	utils.BotInit()

	//middlewares
	admin := utils.Whitelist(append(utils.Config.Admins, utils.Config.SysAdmin)...)
	moder := utils.Whitelist(append(append(utils.Config.Admins, utils.Config.Moders...), utils.Config.SysAdmin)...)
	chats := utils.Whitelist(append(append(utils.Config.Admins, utils.Config.Moders...), utils.Config.SysAdmin, utils.Config.Chat, utils.Config.ReserveChat)...)
	//chatr := utils.Whitelist(utils.Config.ReserveChat)
	//chann := Whitelist(utils.Config.Channel)

	commandMemberList := []commandList{
		{tele.Command{Text: "releases", Description: "список релизов"}, commands.Releases, chats},
		{tele.Command{Text: "russianroulette", Description: "вызвать на дуэль кого-нибудь"}, duel.Request, chats},
		{tele.Command{Text: "savetopm", Description: "сохранить пост в личку"}, commands.SaveToPM, chats},
		{tele.Command{Text: "sed", Description: "заменить текст типо как в sed"}, commands.Sed, chats},
		{tele.Command{Text: "set", Description: "сохранить гет"}, commands.Set, chats},
		{tele.Command{Text: "shrug", Description: "¯\\_(ツ)_/¯"}, commands.Shrug, chats},
		{tele.Command{Text: "stats", Description: "статистика чата"}, commands.Stats, chats},
		{tele.Command{Text: "get", Description: "получить гет"}, commands.Get, chats},
		{tele.Command{Text: "getall", Description: "получить список гетов"}, commands.Getall, chats},
		{tele.Command{Text: "giveme", Description: "сохранить пост в личку"}, commands.SaveToPM, chats},
		{tele.Command{Text: "google", Description: "загуглить что-нибудь"}, commands.Google, chats},
		{tele.Command{Text: "hug", Description: "обнять кого-нибудь"}, commands.Hug, chats},
		{tele.Command{Text: "isekai", Description: "устроиться в роскомнадзор"}, commands.Blessing, chats},
		{tele.Command{Text: "marco", Description: "поло"}, commands.Marco, chats},
		{tele.Command{Text: "me", Description: "аналог команды /me из IRC (/me пошел спать)"}, commands.Me, chats},
		{tele.Command{Text: "meow", Description: "получить гифку с котиком"}, commands.Meow, chats},
		{tele.Command{Text: "mlem", Description: "получить гифку с котиком"}, commands.Meow, chats},
		{tele.Command{Text: "mywarns", Description: "посмотреть количество своих предупреждений"}, commands.Mywarns, chats},
		{tele.Command{Text: "pidor", Description: "запустить игру \"Пидор Дня!\""}, pidor.Pidor, chats},
		{tele.Command{Text: "pidorall", Description: "статистика \"Пидор Дня!\" за всё время"}, pidor.Pidorall, chats},
		{tele.Command{Text: "pidoreg", Description: "зарегистрироваться в \"Пидор Дня!\""}, pidor.Pidoreg, chats},
		{tele.Command{Text: "pidorme", Description: "личная статистика \"Пидор Дня!\""}, pidor.Pidorme, chats},
		{tele.Command{Text: "pidorstats", Description: "статистика \"Пидор Дня!\" за год"}, pidor.Pidorstats, chats},
		{tele.Command{Text: "pidorules", Description: "правила \"Пидор Дня!\""}, pidor.Pidorules, chats},
		{tele.Command{Text: "anekdot", Description: "получить рандомный анекдот с anekdot.ru"}, commands.Anekdot, chats},
		{tele.Command{Text: "blessing", Description: "устроиться в роскомнадзор"}, commands.Blessing, chats},
		{tele.Command{Text: "bonk", Description: "бонкнуть кого-нибудь"}, commands.Bonk, chats},
		{tele.Command{Text: "cur", Description: "посмотреть курс валют"}, commands.Cur, chats},
		{tele.Command{Text: "del", Description: "удалить гет"}, commands.Del, chats},
		{tele.Command{Text: "distort", Description: "переебать медиа"}, commands.Distort, chats},
		{tele.Command{Text: "invert", Description: "инвертировать медиа"}, commands.Invert, chats},
		{tele.Command{Text: "reverse", Description: "инвертировать медиа"}, commands.Invert, chats},
		{tele.Command{Text: "loop", Description: "залупить гифку"}, commands.Loop, chats},
		{tele.Command{Text: "duel", Description: "вызвать на дуэль кого-нибудь"}, duel.Request, chats},
		{tele.Command{Text: "duelstats", Description: "посмотреть свою статистику дуэли"}, duel.Duelstats, chats},
		{tele.Command{Text: "ping", Description: "понг"}, commands.Ping, chats},
		{tele.Command{Text: "slap", Description: "дать леща кому-нибудь"}, commands.Slap, chats},
		{tele.Command{Text: "suicide", Description: "устроиться в роскомнадзор"}, commands.Blessing, chats},
		{tele.Command{Text: "topm", Description: "сохранить пост в личку"}, commands.SaveToPM, chats},
		{tele.Command{Text: "advice", Description: "получить совет"}, commands.Advice, chats},
		{tele.Command{Text: "bet", Description: "поставить ставку"}, bets.Bet, chats},
		{tele.Command{Text: "allbets", Description: "список актуальных ставок"}, bets.AllBets, chats},
		{tele.Command{Text: "delbet", Description: "удалить ставку"}, bets.DelBet, chats},
		{tele.Command{Text: "convert", Description: "конвертировать файл, доппараметры: mp3,ogg,gif,audio,voice,animation"}, commands.Convert, chats},
		{tele.Command{Text: "download", Description: "скачать файл"}, commands.Download, chats},
	}

	commandAdminList := []commandList{
		{tele.Command{Text: "getid", Description: "получить ID юзера"}, commands.Getid, moder},
		{tele.Command{Text: "kick", Description: "кикнуть кого-нибудь"}, commands.Kick, moder},
		{tele.Command{Text: "bite", Description: "укусить кого-нибудь"}, commands.Kill, moder},
		{tele.Command{Text: "kill", Description: "пристрелить кого-нибудь"}, commands.Kill, moder},
		//{tele.Command{Text: "listantispam", Description: "список антиспама"}, checkpoint.ListAntispam, moder},
		{tele.Command{Text: "mute", Description: "заглушить кого-нибудь"}, commands.Mute, moder},
		{tele.Command{Text: "pidordel", Description: "удалить игрока из \"Пидор Дня!\""}, pidor.Pidordel, moder},
		{tele.Command{Text: "pidorlist", Description: "список всех игроков \"Пидор Дня!\""}, pidor.Pidorlist, moder},
		{tele.Command{Text: "restart", Description: "перезапуск бота"}, commands.Restart, admin},
		{tele.Command{Text: "resurrect", Description: "возродить кого-нибудь"}, commands.Revive, moder},
		{tele.Command{Text: "revive", Description: "возродить кого-нибудь"}, commands.Revive, moder},
		{tele.Command{Text: "addbless", Description: "добавить причину блесса"}, commands.AddBless, moder},
		{tele.Command{Text: "addnope", Description: "добавить сообщение отказа по кнопке"}, commands.AddNope, moder},
		{tele.Command{Text: "ban", Description: "забанить кого-нибудь"}, commands.Ban, moder},
		{tele.Command{Text: "bless", Description: "попросить помолчать кого-нибудь"}, commands.Kill, moder},
		{tele.Command{Text: "debug", Description: "получить сообщение в виде JSON"}, commands.Debug, moder},
		//{tele.Command{Text: "delantispam", Description: "удалить из антиспама"}, checkpoint.DelAntispam, moder},
		{tele.Command{Text: "say", Description: "заставить бота сказать что-нибудь"}, commands.Say, moder},
		{tele.Command{Text: "setgetowner", Description: "задать владельца гета"}, commands.SetGetOwner, moder},
		{tele.Command{Text: "unban", Description: "разбанить кого-нибудь"}, commands.Unban, moder},
		{tele.Command{Text: "unmute", Description: "разглушить кого-нибудь"}, commands.Unmute, moder},
		{tele.Command{Text: "warn", Description: "предупредить кого-нибудь"}, commands.Warn, moder},
		//{tele.Command{Text: "testrandom", Description: "протестировать рандом бота "}, commands.TestRandom, moder}
	}

	commandMemberArray := []tele.Command{}
	for i := range commandMemberList {
		utils.Bot.Handle("/"+commandMemberList[i].command.Text, commandMemberList[i].function, commandMemberList[i].middleware)
		commandMemberArray = append(commandMemberArray, commandMemberList[i].command)
	}
	sort.Slice(commandMemberArray, func(i, j int) bool {
		return commandMemberArray[i].Text < commandMemberArray[j].Text
	})
	err := utils.Bot.SetCommands(commandMemberArray, tele.CommandScope{Type: tele.CommandScopeAllGroupChats})
	if err != nil {
		log.Fatal(err)
	}

	commandAdminArray := []tele.Command{}
	for i := range commandAdminList {
		utils.Bot.Handle("/"+commandAdminList[i].command.Text, commandAdminList[i].function, commandAdminList[i].middleware)
		commandAdminArray = append(commandAdminArray, commandAdminList[i].command)
	}
	commandAdminArray = append(commandMemberArray, commandAdminArray...)
	sort.Slice(commandAdminArray, func(i, j int) bool {
		return commandAdminArray[i].Text < commandAdminArray[j].Text
	})
	err = utils.Bot.SetCommands(commandAdminArray, tele.CommandScope{Type: tele.CommandScopeAllChatAdmin})
	if err != nil {
		log.Fatal(err)
	}

	//non-command handles
	utils.Bot.Handle(&duel.AcceptButton, duel.Accept, chats)
	utils.Bot.Handle(&duel.DenyButton, duel.Deny, chats)
	utils.Bot.Handle(tele.OnChatMember, utils.OnChatMember, chats)
	utils.Bot.Handle(tele.OnUserJoined, utils.OnUserJoined, chats)
	utils.Bot.Handle(tele.OnUserLeft, utils.OnUserLeft, chats)
	utils.Bot.Handle(tele.OnText, utils.OnText, chats)
	utils.Bot.Handle(tele.OnQuery, commands.GetInline)
	utils.Bot.Handle(tele.OnChannelPost, utils.ForwardPost)

	utils.Bot.Start()
}
