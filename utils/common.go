package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	tele "gopkg.in/telebot.v3"
)

var WelcomeMessageID int
var RestrictedUsers []CheckPointRestrict
var WordStatsExcludes []WordStatsExclude

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
	b, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
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

type ForwardHistory struct {
	ChannelMessageID int
	ChatMessageID    int
}
type ForwardMesssage struct {
	AlbumID       string
	Caption       string
	History       []ForwardHistory
	StreamHistory []ForwardHistory
}

var Forward ForwardMesssage
var AlbumMessages []*tele.Message

//Repost channel post to chat
func Repost(context tele.Context) error {
	var err error
	var err2 error
	if context.Message().AlbumID != "" && Config.Chat != -1001597398983 {
		AlbumMessages = append(AlbumMessages, context.Message())
		if context.Message().Caption != "" {
			Forward.Caption = context.Message().Caption
		}
		if context.Message().AlbumID != Forward.AlbumID {
			Forward.AlbumID = context.Message().AlbumID
			time.Sleep(5 * time.Second)
			sort.SliceStable(AlbumMessages, func(i, j int) bool {
				return AlbumMessages[i].ID < AlbumMessages[j].ID
			})
			var Album []tele.Inputtable
			for i, message := range AlbumMessages {
				switch {
				case context.Message().Audio != nil:
					message.Audio.Caption = ""
					if i == 0 {
						message.Audio.Caption = Forward.Caption
					}
					Album = append(Album, message.Audio)
				case context.Message().Document != nil:
					message.Document.Caption = ""
					if i == 0 {
						message.Document.Caption = Forward.Caption
					}
					Album = append(Album, message.Document)
				case context.Message().Photo != nil:
					message.Photo.Caption = ""
					if i == 0 {
						message.Photo.Caption = Forward.Caption
					}
					Album = append(Album, message.Photo)
				case context.Message().Video != nil:
					message.Video.Caption = ""
					if i == 0 {
						message.Video.Caption = Forward.Caption
					}
					Album = append(Album, message.Video)
				}
			}
			ChatMessage, err := Bot.SendAlbum(&tele.Chat{ID: Config.Chat}, Album)
			for i, message := range AlbumMessages {
				Forward.History = append(Forward.History, ForwardHistory{message.ID, ChatMessage[i].ID})
			}
			Forward.AlbumID = ""
			AlbumMessages = []*tele.Message{}
			Forward.Caption = ""
			return err
		}
		return nil
	}

	var ChatMessage *tele.Message
	if Config.Chat != -1001597398983 {
		ChatMessage, err = Bot.Copy(&tele.Chat{ID: Config.Chat}, context.Message())
		Forward.History = append(Forward.History, ForwardHistory{context.Message().ID, ChatMessage.ID})
	}
	if Config.StreamChannel != 0 && strings.Contains(context.Text(), "zavtracast/live") {
		ChatMessage, err2 = Bot.Copy(&tele.Chat{ID: Config.StreamChannel}, context.Message())
		Forward.StreamHistory = append(Forward.StreamHistory, ForwardHistory{context.Message().ID, ChatMessage.ID})
	}
	if err2 != nil {
		err = err2
	}
	return err
}

//Edit reposted post
func EditRepost(context tele.Context) error {
	var err error
	var err2 error
	for _, ForwardHistory := range Forward.History {
		if ForwardHistory.ChannelMessageID == context.Message().ID {
			if context.Message().Media() != nil {
				_, err = Bot.Edit(&tele.Message{ID: ForwardHistory.ChatMessageID, Chat: &tele.Chat{ID: Config.Chat}}, context.Message().Media())
				_, err2 = Bot.EditCaption(&tele.Message{ID: ForwardHistory.ChatMessageID, Chat: &tele.Chat{ID: Config.Chat}}, GetHtmlText(*context.Message()))
			} else {
				_, err = Bot.Edit(&tele.Message{ID: ForwardHistory.ChatMessageID, Chat: &tele.Chat{ID: Config.Chat}}, GetHtmlText(*context.Message()))
			}
		}
	}
	if Config.StreamChannel != 0 && strings.Contains(context.Text(), "zavtracast/live") {
		forwarded := false
		for _, ForwardHistory := range Forward.StreamHistory {
			if ForwardHistory.ChannelMessageID == context.Message().ID {
				forwarded = true
				if context.Message().Media() != nil {
					_, err = Bot.Edit(&tele.Message{ID: ForwardHistory.ChatMessageID, Chat: &tele.Chat{ID: Config.StreamChannel}}, context.Message().Media())
					_, err2 = Bot.EditCaption(&tele.Message{ID: ForwardHistory.ChatMessageID, Chat: &tele.Chat{ID: Config.StreamChannel}}, GetHtmlText(*context.Message()))
				} else {
					_, err = Bot.Edit(&tele.Message{ID: ForwardHistory.ChatMessageID, Chat: &tele.Chat{ID: Config.StreamChannel}}, GetHtmlText(*context.Message()))
				}
			}
		}
		if !forwarded {
			var ChatMessage *tele.Message
			ChatMessage, err2 = Bot.Copy(&tele.Chat{ID: Config.StreamChannel}, context.Message())
			Forward.StreamHistory = append(Forward.StreamHistory, ForwardHistory{context.Message().ID, ChatMessage.ID})
		}
	}
	if err2 != nil {
		err = err2
	}
	return err
}

//Remove message
func Remove(context tele.Context) error {
	return context.Delete()
}

func GetNope() string {
	var nope Nope
	DB.Model(Nope{}).Order("RANDOM()").First(&nope)
	return nope.Text
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

func init() {
	//Word stats exclusion list
	DB.Find(&WordStatsExcludes)
}
