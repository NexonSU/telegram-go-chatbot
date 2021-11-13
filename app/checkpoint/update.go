package checkpoint

import (
	"fmt"
	"log"
	"math/rand"
	"regexp"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

type BorderUser struct {
	User     *telebot.User
	Status   string
	Reason   string
	Role     string
	Checked  bool
	JoinedAt int64
}

type JoinBorder struct {
	Message    *telebot.Message
	Chat       *telebot.Chat
	Users      []BorderUser
	NeedUpdate bool
	NeedCreate bool
}

type Question struct {
	Text              string
	CorrectButton     string
	FirstWrongButton  string
	SecondWrongButton string
	ThirdWrongButton  string
}

var Border JoinBorder
var arabicSymbols, _ = regexp.Compile("[\u0600-\u06ff]|[\u0750-\u077f]|[\ufb50-\ufbc1]|[\ufbd3-\ufd3f]|[\ufd50-\ufd8f]|[\ufd92-\ufdc7]|[\ufe70-\ufefc]|[\uFDF0-\uFDFD]")
var Selector = telebot.ReplyMarkup{}
var CorrectButton = Selector.Data("", "")
var FirstWrongButton = Selector.Data("", "")
var SecondWrongButton = Selector.Data("", "")
var ThirdWrongButton = Selector.Data("", "")
var question = ""
var _ = JoinMessageUpdateService(true)

func GetQuestionWithButtons() (string, []telebot.Btn) {
	var questions [][]string
	questTypes := [][][]string{
		{
			{"Просто нажми на мандарин.", "Мандарин", "Оладушек", "Эчпочмак", "Смешная опция"},
		},
		{
			{"Как зовут одного из ведущих подкаста?", "Дмитрий", "Иван", "Пётр", "Александр"},
			{"Как зовут одного из ведущих подкаста?", "Тимур", "Руслан", "Андрей", "Кирилл"},
			{"Как зовут одного из ведущих подкаста?", "Максим", "Миша", "Паша", "Рома"},
		},
		{
			{"Как зовут кота Тимура?", "Борян", "Барсик", "Вискас", "Чилипиздрик"},
			{"Какой подкаст не имеет отношения к этому чату?", "Радио-Т", "Завтракаст", "ДТКД", "Мама, я в стартапе"},
			{"Какой подкаст не имеет отношения к этому чату?", "BeardyCast", "Сторикаст", "ДТКД", "Мама, я в стартапе"},
		},
	}
	i := 0
	if GetSpamChance(*Border.Users[0].User) < 0 {
		questions = questTypes[0]
	} else if GetSpamChance(*Border.Users[0].User) < 50 {
		questions = questTypes[1]
		i = utils.RandInt(0, len(questions)-1)
	} else {
		questions = questTypes[2]
		i = utils.RandInt(0, len(questions)-1)
	}
	CorrectButton.Text = questions[i][1]
	CorrectButton.Data = fmt.Sprintf("%v", utils.RandInt(10000, 99999))
	FirstWrongButton.Text = questions[i][2]
	FirstWrongButton.Data = fmt.Sprintf("%v", utils.RandInt(10000, 99999))
	SecondWrongButton.Text = questions[i][3]
	SecondWrongButton.Data = fmt.Sprintf("%v", utils.RandInt(10000, 99999))
	ThirdWrongButton.Text = questions[i][4]
	ThirdWrongButton.Data = fmt.Sprintf("%v", utils.RandInt(10000, 99999))
	array := []telebot.Btn{CorrectButton, FirstWrongButton, SecondWrongButton, ThirdWrongButton}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(array), func(i, j int) {
		array[i], array[j] = array[j], array[i]
	})
	return questions[i][0], array
}

func JoinMessageUpdateService(init bool) error {
	if init {
		go JoinMessageUpdateService(false)
		return nil
	}
	Border.Message = &telebot.Message{
		ID:       0,
		Unixtime: 0,
		Chat:     &telebot.Chat{ID: 0},
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

func LastNElements(arr []BorderUser, n int) []BorderUser {
	if len(arr) > n {
		return arr[len(arr)-n:]
	}
	return arr
}

func GetUserNames(list []BorderUser) string {
	var text string
	for i, user := range LastNElements(list, 30) {
		if i != 0 {
			text += ", "
		}
		text += utils.MentionUser(user.User)
	}
	return text
}

func JoinMessageUpdate() error {
	var pending []BorderUser
	var banned []BorderUser
	var accepted []BorderUser
	var text string
	for i, user := range Border.Users {
		user = Check(user)
		Border.Users[i] = user
		switch user.Status {
		case "pending":
			pending = append(pending, user)
		case "banned":
			banned = append(banned, user)
		case "accepted":
			accepted = append(accepted, user)
		}
	}
	if Border.NeedCreate {
		var buttons []telebot.Btn
		question, buttons = GetQuestionWithButtons()
		Selector.Inline(
			Selector.Row(buttons[3], buttons[1]),
			Selector.Row(buttons[2], buttons[0]),
		)
	}
	if len(pending) != 0 {
		text += fmt.Sprintf("Добро пожаловать(%d): ", len(pending))
		text += GetUserNames(pending)
		text += "!\nОтветь на вопрос, чтобы получить доступ в чат, иначе бан.\n"
		text += "<b>" + question + "</b>\n"
	} else {
		Selector = telebot.ReplyMarkup{}
	}
	if len(accepted) != 0 {
		text += fmt.Sprintf("Новые подтвержденные пользователи(%d): ", len(accepted))
		text += GetUserNames(accepted)
		text += ".\n"
	}
	if len(banned) != 0 {
		text += fmt.Sprintf("Заблокированные пользователи(%d): ", len(banned))
		/** Не выводим список пользователей, так как в именах или никнеймах может быть реклама */
		text += ".\n"
	}
	if len(Border.Users) > 10 && Border.NeedCreate && time.Now().Unix()-Border.Message.Time().Unix() < 120 {
		Border.NeedCreate = false
		Border.NeedUpdate = true
	}
	if Border.NeedUpdate && !Border.NeedCreate {
		Border.NeedUpdate = false
		_, err := utils.Bot.Edit(Border.Message, text, &Selector)
		return err
	}
	if Border.NeedCreate {
		Border.NeedCreate = false
		Border.NeedUpdate = false
		newMessage, err := utils.Bot.Send(Border.Chat, text, &Selector)
		if err != nil {
			return err
		}
		utils.Bot.Delete(Border.Message)
		Border.Message = newMessage
		return err
	}
	if len(pending) == 0 && time.Now().Unix()-Border.Message.Time().Unix() > 180 {
		Border.Users = []BorderUser{}
		Border.Message = &telebot.Message{
			ID:       0,
			Unixtime: 0,
			Chat:     &telebot.Chat{ID: 0},
		}
	}
	return nil
}
