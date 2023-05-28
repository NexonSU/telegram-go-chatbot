package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

var LastNonAdminChatMember = &tele.ChatMember{}

func UserFullName(user *tele.User) string {
	fullname := user.FirstName
	if user.LastName != "" {
		fullname = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	}
	return fullname
}

func UserName(user *tele.User) string {
	username := user.Username
	if user.Username == "" {
		username = UserFullName(user)
	}
	return username
}

func MentionUser(user *tele.User) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a>", user.ID, UserFullName(user))
}

func RandInt(min int, max int) int {
	b, err := rand.Int(rand.Reader, big.NewInt(int64(max+1)))
	if err != nil {
		return 0
	}
	return min + int(b.Int64())
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsAdmin(userid int64) bool {
	for _, b := range Config.Admins {
		if b == userid {
			return true
		}
	}
	return false
}

func IsAdminOrModer(userid int64) bool {
	for _, b := range Config.Admins {
		if b == userid {
			return true
		}
	}
	for _, b := range Config.Moders {
		if b == userid {
			return true
		}
	}
	return false
}

func RestrictionTimeMessage(seconds int64) string {
	var message = ""
	if seconds-30 > time.Now().Unix() {
		message = fmt.Sprintf(" до %v", time.Unix(seconds, 0).Format("02.01.2006 15:04:05"))
	}
	return message
}

func FindUserInMessage(context tele.Context) (tele.User, int64, error) {
	var user tele.User
	var err error = nil
	var untildate = time.Now().Unix() + 86400
	for _, entity := range context.Message().Entities {
		if entity.Type == tele.EntityTMention {
			user = *entity.User
			if len(context.Args()) == 2 {
				addtime, err := strconv.ParseInt(context.Args()[1], 10, 64)
				if err != nil {
					return user, untildate, err
				}
				untildate += addtime - 86400
			}
			return user, untildate, err
		}
	}
	if context.Message().ReplyTo != nil {
		user = *context.Message().ReplyTo.Sender
		if len(context.Args()) == 1 {
			addtime, err := strconv.ParseInt(context.Args()[0], 10, 64)
			if err != nil {
				return user, untildate, errors.New("время указано неверно")
			}
			untildate += addtime - 86400
		}
	} else {
		if len(context.Args()) == 0 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		if user.ID == 0 {
			user, err = GetUserFromDB(context.Args()[0])
			if err != nil {
				return user, untildate, err
			}
		}
		if len(context.Args()) == 2 {
			addtime, err := strconv.ParseInt(context.Args()[1], 10, 64)
			if err != nil {
				return user, untildate, errors.New("время указано неверно")
			}
			untildate += addtime - 86400
		}
	}
	return user, untildate, err
}

func GetUserFromDB(findstring string) (tele.User, error) {
	var user tele.User
	var err error = nil
	if string(findstring[0]) == "@" {
		user.Username = findstring[1:]
	} else {
		user.ID, err = strconv.ParseInt(findstring, 10, 64)
	}
	result := DB.Where("lower(username) = ? OR id = ?", strings.ToLower(user.Username), user.ID).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return user, err
}

// Forward channel post to chat
func ForwardPost(context tele.Context) error {
	if context.Message().Chat.ID != Config.Channel {
		return nil
	}
	_, err := Bot.Forward(&tele.Chat{ID: Config.Chat}, context.Message())
	if Config.StreamChannel != 0 && strings.Contains(context.Text(), "zavtracast/live") {
		_, err = Bot.Forward(&tele.Chat{ID: Config.StreamChannel}, context.Message())
	}
	if err != nil {
		return err
	}
	return nil
}

// Remove message
func Remove(context tele.Context) error {
	return context.Delete()
}

func OnChatMember(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return Bot.Unban(&tele.Chat{ID: context.Chat().ID}, context.ChatMember().NewChatMember.User)
	}
	//User update
	UserResult := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(context.ChatMember().NewChatMember.User)
	if UserResult.Error != nil {
		ErrorReporting(UserResult.Error, nil)
	}
	return nil
}

func OnUserJoined(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return context.Delete()
	}
	return nil
}

func OnUserLeft(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return context.Delete()
	}
	return nil
}

func OnText(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return context.Delete()
	}
	chatMember, err := Bot.ChatMemberOf(context.Chat(), context.Sender())
	if err != nil {
		ErrorReporting(err, nil)
	}
	if chatMember.Role == tele.Member {
		LastNonAdminChatMember = chatMember
	}
	return nil
}

func GetNope() string {
	var nope Nope
	DB.Model(Nope{}).Order("RANDOM()").First(&nope)
	return nope.Text
}

func GetBless() string {
	var bless Bless
	DB.Model(Bless{}).Order("RANDOM()").First(&bless)
	return bless.Text
}

func GetHtmlText(message tele.Message) string {
	type entity struct {
		s string
		i int
	}

	entities := message.Entities
	textString := message.Text

	if len(message.Text) == 0 {
		entities = message.CaptionEntities
		textString = message.Caption
	}

	textString = strings.ReplaceAll(textString, "<", "˂")
	textString = strings.ReplaceAll(textString, ">", "˃")
	text := utf16.Encode([]rune(textString))

	ents := make([]entity, 0, len(entities)*2)

	for _, ent := range entities {
		var a, b string

		switch ent.Type {
		case tele.EntityBold, tele.EntityItalic,
			tele.EntityUnderline, tele.EntityStrikethrough:
			a = fmt.Sprintf("<%c>", ent.Type[0])
			b = a[:1] + "/" + a[1:]
		case tele.EntityCode, tele.EntityCodeBlock:
			a = fmt.Sprintf("<%s>", ent.Type)
			b = a[:1] + "/" + a[1:]
		case tele.EntityTextLink:
			a = fmt.Sprintf("<a href='%s'>", ent.URL)
			b = "</a>"
		case tele.EntityTMention:
			a = fmt.Sprintf("<a href='tg://user?id=%d'>", ent.User.ID)
			b = "</a>"
		default:
			continue
		}

		ents = append(ents, entity{a, ent.Offset})
		ents = append(ents, entity{b, ent.Offset + ent.Length})
	}

	// reverse entities
	for i, j := 0, len(ents)-1; i < j; i, j = i+1, j-1 {
		ents[i], ents[j] = ents[j], ents[i]
	}

	for _, ent := range ents {
		r := utf16.Encode([]rune(ent.s))
		text = append(text[:ent.i], append(r, text[ent.i:]...)...)
	}

	textString = string(utf16.Decode(text))

	if len(message.Entities) != 0 && message.Entities[0].Type == tele.EntityCommand {
		if textString[1:4] == "set" {
			textString = strings.Join(strings.Split(textString, " ")[2:], " ")
		} else {
			textString = textString[message.Entities[0].Length+1:]
		}
	}

	return textString
}
