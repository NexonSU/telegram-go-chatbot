package welcome

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/valyala/fastjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"io"
	"io/ioutil"
	"log"
	pseudorand "math/rand"
	"net/http"
	"regexp"
	"strconv"
	"time"
)

var Message *tb.Message
var Selector = tb.ReplyMarkup{}
var FirstWrongButton = Selector.Data("Джабир, Латиф и Хиляль", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
var CorrectButton = Selector.Data("Дмитрий, Тимур и Максим", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
var SecondWrongButton = Selector.Data("Бубылда, Чингачгук и Гавкошмыг", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
var ThirdWrongButton = Selector.Data("Мандарин, Оладушек и Эчпочмак", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))

func shuffleButtons(array []tb.Btn) []tb.Btn {
	pseudorand.Seed(time.Now().UnixNano())
	pseudorand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})
	return array
}

var buttons = shuffleButtons([]tb.Btn{CorrectButton, FirstWrongButton, SecondWrongButton, ThirdWrongButton})

var arabicSymbols, _ = regexp.Compile("[\u0600-\u06ff]|[\u0750-\u077f]|[\ufb50-\ufbc1]|[\ufbd3-\ufd3f]|[\ufd50-\ufd8f]|[\ufd92-\ufdc7]|[\ufe70-\ufefc]|[\uFDF0-\uFDFD]")

func OnJoin(m *tb.Message) {
	if Message == nil {
		Message = m
		Message.Unixtime = 0
	}
	if m.Chat.Username != utils.Config.Telegram.Chat {
		return
	}
	err := utils.Bot.Delete(m)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	log.Printf("New user detected in %v (%v)! ID: %v. Login: %v. Name: %v.", m.Chat.Title, m.Chat.ID, m.Sender.ID, utils.UserName(m.Sender), utils.UserFullName(m.Sender))
	User := m.Sender
	Chat := m.Chat
	ChatMember, err := utils.Bot.ChatMemberOf(Chat, User)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	ChatMember.CanSendMessages = false
	err = utils.Bot.Restrict(Chat, ChatMember)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	if arabicSymbols.MatchString(utils.UserFullName(User)) || User.FirstName == "ICSM" {
		err = utils.Bot.Ban(Chat, ChatMember)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
	var httpClient = &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", User.ID))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(httpResponse.Body)
	jsonBytes, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	if fastjson.GetBool(jsonBytes, "ok") {
		err = utils.Bot.Ban(Chat, ChatMember)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
	Selector.Inline(
		Selector.Row(buttons[0]),
		Selector.Row(buttons[1]),
		Selector.Row(buttons[2]),
		Selector.Row(buttons[3]),
	)
	if time.Now().Unix()-Message.Time().Unix() > 10 {
		Message, err = utils.Bot.Send(Chat, fmt.Sprintf("Добро пожаловать, %v!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста?", utils.MentionUser(User)), &Selector)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	} else {
		text := "Добро пожаловать"
		for _, element := range Message.Entities {
			text += ", " + utils.MentionUser(element.User)
		}
		text += ", " + utils.MentionUser(m.Sender) + "!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста? У тебя 2 минуты."
		_, err = utils.Bot.Edit(Message, text, &Selector)
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
	go func() {
		message := Message
		time.Sleep(time.Second * 120)
		ChatMember, err := utils.Bot.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			return
		}
		if ChatMember.Role != "member" {
			err := utils.Bot.Ban(m.Chat, &tb.ChatMember{User: m.Sender})
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
		}
		err = utils.Bot.Delete(message)
		if err != nil {
			return
		}
	}()
}
