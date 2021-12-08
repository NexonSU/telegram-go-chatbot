package utils

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"log"
	"math/big"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf16"

	"gopkg.in/tucnak/telebot.v3"
	"gorm.io/gorm/clause"
)

var LastChatMessageID int
var WelcomeMessageID int
var RestrictedUsers []CheckPointRestrict
var WordStatsExcludes []WordStatsExclude

func UserFullName(user *telebot.User) string {
	fullname := user.FirstName
	if user.LastName != "" {
		fullname = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	}
	return fullname
}

func UserName(user *telebot.User) string {
	username := user.Username
	if user.Username == "" {
		username = UserFullName(user)
	}
	return username
}

func MentionUser(user *telebot.User) string {
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

func ErrorReporting(err error, context telebot.Context) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[%s:%d] %v", fn, line, err)
	text := fmt.Sprintf("<pre>[%s:%d]\n%v</pre>", fn, line, err)
	if strings.Contains(err.Error(), "specified new message content and reply markup are exactly the same") {
		return
	}
	if strings.Contains(err.Error(), "message to delete not found") {
		return
	}
	if context.Message() != nil {
		MarshalledMessage, _ := json.MarshalIndent(context.Message(), "", "    ")
		JsonMessage := html.EscapeString(string(MarshalledMessage))
		text += fmt.Sprintf("\n\nMessage:\n<pre>%v</pre>", JsonMessage)
	}
	Bot.Send(telebot.ChatID(Config.SysAdmin), text)
}

func FindUserInMessage(context telebot.Context) (telebot.User, int64, error) {
	var user telebot.User
	var err error = nil
	var untildate = time.Now().Unix() + 86400
	for _, entity := range context.Message().Entities {
		if entity.Type == telebot.EntityTMention {
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

func CheckPoint(update *telebot.Update) error {
	if update.Message == nil || update.Message.Sender == nil {
		return nil
	}
	if update.Message.SenderChat != nil &&
		(update.Message.Chat.ID == Config.Chat ||
			update.Message.Chat.ID == Config.CommentChat) &&
		update.Message.Sender.ID == 777000 &&
		update.Message.SenderChat.ID != Config.Channel {
		return Bot.Delete(update.Message)
	}
	LastChatMessageID = update.Message.ID
	if update.Message.ReplyTo != nil && update.Message.ReplyTo.ID == WelcomeMessageID {
		delete := DB.Delete(CheckPointRestrict{UserID: update.Message.Sender.ID})
		if delete.Error != nil {
			return delete.Error
		}
		find := DB.Find(&RestrictedUsers)
		if find.Error != nil {
			return find.Error
		}
	}
	for _, user := range RestrictedUsers {
		if update.Message.Sender.ID == user.UserID {
			return Bot.Delete(update.Message)
		}
	}
	return nil
}

func GatherData(update *telebot.Update) error {
	if update.Message == nil || update.Message.Sender == nil {
		return nil
	}
	//User update
	UserResult := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(update.Message.Sender)
	if UserResult.Error != nil {
		return UserResult.Error
	}
	if update.Message.Sender.IsBot || update.Message.Chat.ID != Config.Chat {
		return nil
	}
	//Message insert
	var Message Message
	Message.ID = update.Message.ID
	Message.UserID = update.Message.Sender.ID
	Message.Date = time.Unix(update.Message.Unixtime, 0)
	Message.ChatID = update.Message.Chat.ID
	if update.Message.ReplyTo != nil {
		Message.ReplyTo = update.Message.ReplyTo.ID
	}
	Message.Text = update.Message.Text
	switch {
	case update.Message.Animation != nil:
		Message.FileType = "Animation"
		Message.FileID = update.Message.Animation.FileID
		Message.Text = update.Message.Caption
	case update.Message.Audio != nil:
		Message.FileType = "Audio"
		Message.FileID = update.Message.Audio.FileID
		Message.Text = update.Message.Caption
	case update.Message.Photo != nil:
		Message.FileType = "Photo"
		Message.FileID = update.Message.Photo.FileID
		Message.Text = update.Message.Caption
	case update.Message.Video != nil:
		Message.FileType = "Video"
		Message.FileID = update.Message.Video.FileID
		Message.Text = update.Message.Caption
	case update.Message.Voice != nil:
		Message.FileType = "Voice"
		Message.FileID = update.Message.Voice.FileID
		Message.Text = update.Message.Caption
	case update.Message.Document != nil:
		Message.FileType = "Document"
		Message.FileID = update.Message.Document.FileID
		Message.Text = update.Message.Caption
	}
	MessageResult := DB.Create(&Message)
	if MessageResult.Error != nil {
		return MessageResult.Error
	}
	//Words insert
	if Message.Text == "" || string(Message.Text[0]) == "/" {
		return nil
	}
	var Word Word
	Word.ChatID = Message.ChatID
	Word.UserID = Message.UserID
	Word.Date = Message.Date
	Message.Text = strings.ReplaceAll(Message.Text, ",", "")
	Message.Text = strings.ReplaceAll(Message.Text, ".", "")
	Message.Text = strings.ReplaceAll(Message.Text, "!", "")
	Message.Text = strings.ReplaceAll(Message.Text, "?", "")
words:
	for _, Word.Text = range strings.Fields(strings.ToLower(Message.Text)) {
		for _, exclude := range WordStatsExcludes {
			if Word.Text == exclude.Text {
				continue words
			}
		}
		if _, err := strconv.Atoi(Word.Text); err == nil ||
			Word.Text == "" ||
			len(Word.Text) == 1 {
			continue
		}
		WordResult := DB.Create(&Word)
		if WordResult.Error != nil {
			return WordResult.Error
		}
	}
	return nil
}

func GetUserFromDB(findstring string) (telebot.User, error) {
	var user telebot.User
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
var AlbumMessages []*telebot.Message

//Repost channel post to chat
func Repost(context telebot.Context) error {
	var err error
	var err2 error
	if context.Message().AlbumID != "" {
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
			var Album []telebot.Inputtable
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
			ChatMessage, err := Bot.SendAlbum(&telebot.Chat{ID: Config.Chat}, Album)
			for i, message := range AlbumMessages {
				Forward.History = append(Forward.History, ForwardHistory{message.ID, ChatMessage[i].ID})
			}
			Forward.AlbumID = ""
			AlbumMessages = []*telebot.Message{}
			Forward.Caption = ""
			return err
		}
		return nil
	}

	var ChatMessage *telebot.Message
	ChatMessage, err = Bot.Copy(&telebot.Chat{ID: Config.Chat}, context.Message())
	Forward.History = append(Forward.History, ForwardHistory{context.Message().ID, ChatMessage.ID})
	if Config.StreamChannel != 0 && strings.Contains(context.Text(), "zavtracast/live") {
		ChatMessage, err2 = Bot.Copy(&telebot.Chat{ID: Config.StreamChannel}, context.Message())
		Forward.StreamHistory = append(Forward.StreamHistory, ForwardHistory{context.Message().ID, ChatMessage.ID})
	}
	if err2 != nil {
		err = err2
	}
	return err
}

//Edit reposted post
func EditRepost(context telebot.Context) error {
	var err error
	var err2 error
	for _, ForwardHistory := range Forward.History {
		if ForwardHistory.ChannelMessageID == context.Message().ID {
			if context.Media() != nil {
				_, err = Bot.Edit(&telebot.Message{ID: ForwardHistory.ChatMessageID, Chat: &telebot.Chat{ID: Config.Chat}}, context.Media())
				_, err2 = Bot.EditCaption(&telebot.Message{ID: ForwardHistory.ChatMessageID, Chat: &telebot.Chat{ID: Config.Chat}}, GetHtmlText(*context.Message()))
			} else {
				_, err = Bot.Edit(&telebot.Message{ID: ForwardHistory.ChatMessageID, Chat: &telebot.Chat{ID: Config.Chat}}, GetHtmlText(*context.Message()))
			}
		}
	}
	if Config.StreamChannel != 0 && strings.Contains(context.Text(), "zavtracast/live") {
		forwarded := false
		for _, ForwardHistory := range Forward.StreamHistory {
			if ForwardHistory.ChannelMessageID == context.Message().ID {
				forwarded = true
				if context.Media() != nil {
					_, err = Bot.Edit(&telebot.Message{ID: ForwardHistory.ChatMessageID, Chat: &telebot.Chat{ID: Config.StreamChannel}}, context.Media())
					_, err2 = Bot.EditCaption(&telebot.Message{ID: ForwardHistory.ChatMessageID, Chat: &telebot.Chat{ID: Config.StreamChannel}}, GetHtmlText(*context.Message()))
				} else {
					_, err = Bot.Edit(&telebot.Message{ID: ForwardHistory.ChatMessageID, Chat: &telebot.Chat{ID: Config.StreamChannel}}, GetHtmlText(*context.Message()))
				}
			}
		}
		if !forwarded {
			var ChatMessage *telebot.Message
			ChatMessage, err2 = Bot.Copy(&telebot.Chat{ID: Config.StreamChannel}, context.Message())
			Forward.StreamHistory = append(Forward.StreamHistory, ForwardHistory{context.Message().ID, ChatMessage.ID})
		}
	}
	if err2 != nil {
		err = err2
	}
	return err
}

//Remove message
func Remove(context telebot.Context) error {
	return context.Delete()
}

func GetNope() string {
	var nope Nope
	DB.Model(Nope{}).Order("RANDOM()").First(&nope)
	return nope.Text
}

func GetHtmlText(message telebot.Message) string {
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
		case telebot.EntityBold, telebot.EntityItalic,
			telebot.EntityUnderline, telebot.EntityStrikethrough:
			a = fmt.Sprintf("<%c>", ent.Type[0])
			b = a[:1] + "/" + a[1:]
		case telebot.EntityCode, telebot.EntityCodeBlock:
			a = fmt.Sprintf("<%s>", ent.Type)
			b = a[:1] + "/" + a[1:]
		case telebot.EntityTextLink:
			a = fmt.Sprintf("<a href='%s'>", ent.URL)
			b = "</a>"
		case telebot.EntityTMention:
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

	if len(message.Entities) != 0 && message.Entities[0].Type == telebot.EntityCommand {
		if textString[1:4] == "set" {
			textString = strings.Join(strings.Split(textString, " ")[2:], " ")
		} else {
			textString = textString[message.Entities[0].Length+1:]
		}
	}

	return textString
}
