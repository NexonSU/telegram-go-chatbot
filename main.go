package main

import (
	"github.com/NexonSU/telegram-go-chatbot/app/commands"
	"github.com/NexonSU/telegram-go-chatbot/app/roulette"
	"github.com/NexonSU/telegram-go-chatbot/app/services"
	"github.com/NexonSU/telegram-go-chatbot/app/userActions"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	pseudorand "math/rand"
	"regexp"
	"strconv"
	"time"
)

func main() {
	var busy = make(map[string]bool)

	// commands
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
	utils.Bot.Handle("/me", commands.Me)
	utils.Bot.Handle("/topic", commands.Topic)
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
	utils.Bot.Handle("/pidor", commands.Pidor(busy))
	utils.Bot.Handle("/blessing", commands.Blessing)
	utils.Bot.Handle("/suicide", commands.Suicide)
	utils.Bot.Handle("/kill", commands.Kill)
	utils.Bot.Handle("/duelstats", commands.Duelstats)

	// Roulette
	var russianRouletteMessage *tb.Message
	russianRouletteSelector := tb.ReplyMarkup{}
	russianRouletteAcceptButton := russianRouletteSelector.Data("👍 Принять вызов", "russianroulette_accept")
	russianRouletteDenyButton := russianRouletteSelector.Data("👎 Бежать с позором", "russianroulette_deny")
	russianRouletteSelector.Inline(
		russianRouletteSelector.Row(russianRouletteAcceptButton, russianRouletteDenyButton),
	)
	utils.Bot.Handle("/russianroulette", roulette.RussianRoulette(busy, russianRouletteMessage, russianRouletteSelector))
	utils.Bot.Handle(&russianRouletteAcceptButton, roulette.RouletteAccept(busy))
	utils.Bot.Handle(&russianRouletteDenyButton, roulette.RouletteDeny(busy))

	// Gather text
	utils.Bot.Handle(tb.OnText, userActions.OnText)
	utils.Bot.Handle(tb.OnChannelPost, userActions.OnPost)

	// User join
	var welcomeMessage *tb.Message
	welcomeSelector := tb.ReplyMarkup{}
	welcomeFirstWrongButton := welcomeSelector.Data("Джабир, Латиф и Хиляль", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
	welcomeRightButton := welcomeSelector.Data("Дмитрий, Тимур и Максим", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
	welcomeSecondWrongButton := welcomeSelector.Data("Бубылда, Чингачгук и Гавкошмыг", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
	welcomeThirdWrongButton := welcomeSelector.Data("Мандарин, Оладушек и Эчпочмак", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
	buttons := []tb.Btn{welcomeRightButton, welcomeFirstWrongButton, welcomeSecondWrongButton, welcomeThirdWrongButton}
	pseudorand.Seed(time.Now().UnixNano())
	pseudorand.Shuffle(len(buttons), func(i, j int) {
		buttons[i], buttons[j] = buttons[j], buttons[i]
	})
	welcomeSelector.Inline(
		welcomeSelector.Row(buttons[0], buttons[1]),
		welcomeSelector.Row(buttons[2], buttons[3]),
	)

	arab, err := regexp.Compile("[\u0600-\u06ff]|[\u0750-\u077f]|[\ufb50-\ufbc1]|[\ufbd3-\ufd3f]|[\ufd50-\ufd8f]|[\ufd92-\ufdc7]|[\ufe70-\ufefc]|[\uFDF0-\uFDFD]")
	if err != nil {
		log.Fatal(err)
		return
	}
	utils.Bot.Handle(tb.OnUserJoined, userActions.OnJoin(welcomeMessage, welcomeSelector, arab))
	utils.Bot.Handle(tb.OnUserLeft, userActions.OnLeft)

	utils.Bot.Handle(&welcomeRightButton, userActions.OnClickRightButton)
	utils.Bot.Handle(&welcomeFirstWrongButton, userActions.OnClickWrongButton)
	utils.Bot.Handle(&welcomeSecondWrongButton, userActions.OnClickWrongButton)
	utils.Bot.Handle(&welcomeThirdWrongButton, userActions.OnClickWrongButton)

	go services.ZavtraStreamCheckService()
	utils.Bot.Start()
}
