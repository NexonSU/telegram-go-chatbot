package welcome

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"log"
	"math/rand"
	"strconv"
	"time"
)

var Selector = tb.ReplyMarkup{}
var CorrectButton = Selector.Data("Дмитрий, Тимур, Максим", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
var FirstWrongButton = Selector.Data("Франсуа, Гюго, Эркюль", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
var SecondWrongButton = Selector.Data("Хасан, Камаран, Хидео", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))
var ThirdWrongButton = Selector.Data("Ганс, Шульц, Фриц", "Button"+strconv.Itoa(utils.RandInt(10000, 99999)))

func shuffleButtons(array []tb.Btn) []tb.Btn {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})
	return array
}

var buttons = shuffleButtons([]tb.Btn{CorrectButton, SecondWrongButton, ThirdWrongButton})

func JoinMessageUpdateService() {
	Border.Message = &tb.Message{
		ID:       0,
		Unixtime: 0,
		Chat:     &tb.Chat{ID: 0},
	}
	for {
		delay := 1
		err := JoinMessageUpdate()
		if err != nil {
			log.Println(err.Error())
		}
		delay = len(Border.Users)
		if delay < 1 {
			delay = 1
		}
		if delay > 4 {
			delay = 4
		}
		time.Sleep(time.Second * time.Duration(delay))
	}
}

func JoinMessageUpdate() error {
	var pending []BorderUser
	var banned []BorderUser
	var accepted []BorderUser
	var text string
	for i, user := range Border.Users {
		switch user.Status {
		case "pending":
			if time.Now().Unix()-user.JoinedAt.Unix() > 120 {
				err := utils.Bot.Ban(Border.Chat, &tb.ChatMember{User: user.User})
				if err != nil {
					continue
				}
				Border.Users[i].Status = "banned"
				user.Status = "banned"
				Border.Users[i].Reason = "не прошел проверку"
				user.Reason = "не прошел проверку"
				banned = append(banned, user)
				Border.NeedUpdate = true
			} else {
				pending = append(pending, user)
			}
		case "banned":
			banned = append(banned, user)
		case "accepted":
			accepted = append(accepted, user)
		}
	}
	Selector.Inline(
		Selector.Row(FirstWrongButton),
		Selector.Row(buttons[0]),
		Selector.Row(buttons[1]),
		Selector.Row(buttons[2]),
	)
	if len(pending) != 0 {
		text += "Добро пожаловать: "
		for i, user := range pending {
			if i != 0 {
				text += ", "
			}
			text += utils.MentionUser(user.User)
		}
		text += "!\nОтветь на вопрос, чтобы получить доступ в чат, иначе бан через 2 минуты.\nКак зовут ведущих подкаста?\n"
	} else {
		Selector = tb.ReplyMarkup{}
	}
	if len(accepted) != 0 {
		text += "Новые подтвержденные пользователи: "
		for i, user := range accepted {
			if i != 0 {
				text += ", "
			}
			text += utils.MentionUser(user.User)
		}
		text += ".\n"
	}
	if len(banned) != 0 {
		text += "Заблокированные пользователи: "
		for i, user := range banned {
			if i != 0 {
				text += ", "
			}
			text += utils.MentionUser(user.User)
			text += " (" + user.Reason + ")"
		}
		text += ".\n"
	}
	if Border.NeedUpdate && !Border.NeedCreate {
		Border.NeedUpdate = false
		_, err := utils.Bot.Edit(Border.Message, text, &Selector)
		if err != nil {
			return err
		}
		return nil
	}
	if Border.NeedCreate {
		Border.NeedCreate = false
		Border.NeedUpdate = false
		newMessage, err := utils.Bot.Send(Border.Chat, text, &Selector)
		if err != nil {
			return err
		}
		_ = utils.Bot.Delete(Border.Message)
		Border.Message = newMessage
		return nil
	}
	if len(pending) == 0 && time.Now().Unix()-Border.Message.Time().Unix() > 60 {
		Border.Users = []BorderUser{}
		Border.Message = &tb.Message{
			ID:       0,
			Unixtime: 0,
			Chat:     &tb.Chat{ID: 0},
		}
	}
	return nil
}
