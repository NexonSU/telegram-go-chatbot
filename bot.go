package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	ical "github.com/arran4/golang-ical"
	"github.com/chai2010/webp"
	"github.com/fogleman/gg"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	"github.com/valyala/fastjson"
	tb "gopkg.in/tucnak/telebot.v2"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"html"
	"io"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	pseudorand "math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)
type Configuration struct {
	Telegram struct {
		Token         string   `json:"token,omitempty"`
		Chat          string   `json:"chat,omitempty"`
		StreamChannel string   `json:"stream_channel,omitempty"`
		Channel       string   `json:"channel,omitempty"`
		BotApiUrl     string   `json:"bot_api_url,omitempty"`
		Admins        []string `json:"admins,omitempty"`
		Moders        []string `json:"moders,omitempty"`
		SysAdmin      string   `json:"sysadmin,omitempty"`
	}
	Webhook struct {
		Listen         string   `json:"listen,omitempty"`
		Port           int      `json:"port,omitempty"`
		AllowedUpdates []string `json:"allowed_updates,omitempty"`
	}
	Youtube struct {
		ApiKey      string `json:"api_key,omitempty"`
		ChannelName string `json:"channel_name,omitempty"`
		ChannelID   string `json:"channel_id,omitempty"`
	}
	CurrencyKey string `json:"currency_key,omitempty"`
	ReleasesUrl string `json:"releases_url,omitempty"`
}
type Get struct {
	Name     string `gorm:"primaryKey"`
	Type     string
	Data     string
	Caption  string
}
type Warn struct {
	UserID     int `gorm:"primaryKey"`
	Amount     int
	LastWarn   time.Time
}
type PidorStats struct {
	Date       time.Time `gorm:"primaryKey"`
	UserID     int
}
type PidorList tb.User
type ZavtraStream struct {
	Service     string     `gorm:"primaryKey"`
	LastCheck   time.Time
	VideoID		string
}
type Duelist struct {
	UserID      int     `gorm:"primaryKey"`
	Deaths		int
	Kills		int
}
var ConfigFile, _ = os.Open("config.json")
var Config = new(Configuration)
var _ = json.NewDecoder(ConfigFile).Decode(&Config)
var Bot, _ = tb.NewBot(tb.Settings{
	URL:       Config.Telegram.BotApiUrl,
	Token:     Config.Telegram.Token,
	ParseMode: tb.ModeHTML,
	Poller: &tb.LongPoller{
		Timeout:        10 * time.Second,
		AllowedUpdates: Config.Webhook.AllowedUpdates,
	},
})
var DB, _ = gorm.Open(sqlite.Open("bot.db"), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
var busy = make(map[string]bool)
func ErrorReporting(err error, message *tb.Message)  {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[%s:%d] %v at MessageID \"%v\" in Chat \"%v\"", fn, line, err, message.ID, message.Chat.Username)
	MarshalledMessage, _ := json.MarshalIndent(message, "", "    ")
	JsonMessage := html.EscapeString(string(MarshalledMessage))
	text := fmt.Sprintf("An exception was raised while handling an update\n<pre>%v</pre>\n\nMessage:\n<pre>%v</pre>", err, JsonMessage)
	chat, err := Bot.ChatByID("@" + Config.Telegram.SysAdmin)
	if err != nil {
		return
	}
	_, err = Bot.Send(chat, text)
	if err != nil {
		return
	}
}
func GatherData(user *tb.User) error {
	result := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func UserFullName(user *tb.User) string {
	fullname := user.FirstName
	if user.LastName != "" {
		fullname = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	}
	return fullname
}
func UserName(user *tb.User) string {
	username := user.Username
	if user.Username == "" {
		username = UserFullName(user)
	}
	return username
}
func MentionUser(user *tb.User) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a>", user.ID, UserFullName(user))
}
func RestrictionTimeMessage(seconds int64) string {
	var message = ""
	if seconds-30 > time.Now().Unix() {
		message = fmt.Sprintf(" до %v", time.Unix(seconds, 0).Format("02.01.2006 15:04:05"))
	}
	return message
}
func GetUserFromDB(findstring string) (tb.User, error) {
	var user tb.User
	var err error = nil
	if string(findstring[0]) == "@" {
		user.Username = findstring[1:]
	} else {
		user.ID, err = strconv.Atoi(findstring)
	}
	result := DB.Where(&user).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return user, err
}
func FindUserInMessage(m tb.Message) (tb.User, int64, error) {
	var user tb.User
	var err error = nil
	var untildate = time.Now().Unix()
	var text = strings.Split(m.Text, " ")
	if m.ReplyTo != nil {
		user = *m.ReplyTo.Sender
		if len(text) == 2 {
			addtime, err := strconv.ParseInt(text[1], 10,64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	} else {
		if len(text) == 1 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		user, err = GetUserFromDB(text[1])
		if err != nil {
			return user, untildate, err
		}
		if len(text) == 3 {
			addtime, err := strconv.ParseInt(text[2], 10,64)
			if err != nil {
				return user, untildate, err
			}
			untildate += addtime
		}
	}
	return user, untildate, err
}
func ZavtraStreamCheck(service string) error {
	if service == "youtube" {
		var stream ZavtraStream
		var httpClient = &http.Client{Timeout: 10 * time.Second}
		r, err := httpClient.Get(fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%v&type=video&eventType=live&key=%v", Config.Youtube.ChannelID, Config.Youtube.ApiKey))
		if err != nil {
			return err
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				return
			}
		}(r.Body)
		jsonBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Fatal(err)
		}
		stream.Service = service
		DB.First(&stream)
		results := fastjson.GetInt(jsonBytes, "pageInfo", "totalResults")
		if results != 0 {
			title := fastjson.GetString(jsonBytes, "items", "0", "snippet", "title")
			videoId := fastjson.GetString(jsonBytes, "items", "0", "id", "videoId")
			if stream.VideoID != videoId {
				thumbnail := fmt.Sprintf("https://i.ytimg.com/vi/%v/maxresdefault_live.jpg", videoId)
				caption := fmt.Sprintf("Стрим \"%v\" начался.\nhttps://youtube.com/%v/live", title, Config.Youtube.ChannelName)
				chat, err := Bot.ChatByID("@"+Config.Telegram.StreamChannel)
				if err != nil {
					return err
				}
				_, err = Bot.Send(chat, &tb.Photo{File: tb.File{FileURL: thumbnail}, Caption: caption})
				if err != nil {
					return err
				}
				stream.VideoID = videoId
			}
		}
		stream.LastCheck = time.Now()
		result := DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(stream)
		if result.Error != nil {
			return result.Error
		}
		return nil
	}
	return nil
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

func main() {
	//Create tables, if they not exists in DB
	//DB.AutoMigrate(tb.User{})
	//DB.AutoMigrate(Get{})
	//DB.AutoMigrate(Warn{})
	//DB.AutoMigrate(PidorStats{})
	//DB.AutoMigrate(PidorList{})
	//DB.AutoMigrate(Duelist{})
	//DB.AutoMigrate(ZavtraStream{})

	//Send admin list to user on /admin
	Bot.Handle("/admin", func(m *tb.Message) {
		var get Get
		result := DB.Where(&Get{Name: "admin"}).First(&get)
		if result.RowsAffected != 0 {
			switch {
			case get.Type == "Animation":
				_, err := Bot.Reply(m, &tb.Animation{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Audio":
				_, err := Bot.Reply(m, &tb.Audio{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Photo":
				_, err := Bot.Reply(m, &tb.Photo{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Video":
				_, err := Bot.Reply(m, &tb.Video{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Voice":
				_, err := Bot.Reply(m, &tb.Voice{
					File:      tb.File{FileID: get.Data},
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Document":
				_, err := Bot.Reply(m, &tb.Document{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Text":
				_, err := Bot.Reply(m, get.Data)
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			default:
				_, err := Bot.Reply(m, fmt.Sprintf("Ошибка при определении типа гета, я не знаю тип <code>%v</code>.", get.Type))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			}
		} else {
			_, err := Bot.Reply(m, "Гет <code>admin</code> не найден.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Return message on /debug command
	Bot.Handle("/debug", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var message = m
		if m.ReplyTo != nil {
			message = m.ReplyTo
		}
		MarshalledMessage, _ := json.MarshalIndent(message, "", "    ")
		_, err := Bot.Reply(m, fmt.Sprintf("<pre>%v</pre>", string(MarshalledMessage)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send Get to user on /get
	Bot.Handle("/get", func(m *tb.Message) {
		var get Get
		var text = strings.Split(m.Text, " ")
		if len(text) != 2 {
			_, err := Bot.Reply(m, "Пример использования: <code>/get {гет}</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		result := DB.Where(&Get{Name: strings.ToLower(text[1])}).First(&get)
		if result.RowsAffected != 0 {
			switch {
			case get.Type == "Animation":
				_, err := Bot.Reply(m, &tb.Animation{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Audio":
				_, err := Bot.Reply(m, &tb.Audio{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Photo":
				_, err := Bot.Reply(m, &tb.Photo{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Video":
				_, err := Bot.Reply(m, &tb.Video{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Voice":
				_, err := Bot.Reply(m, &tb.Voice{
					File:      tb.File{FileID: get.Data},
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Document":
				_, err := Bot.Reply(m, &tb.Document{
					File:      tb.File{FileID: get.Data},
					Caption:   get.Caption,
				})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			case get.Type == "Text":
				_, err := Bot.Reply(m, get.Data)
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			default:
				_, err := Bot.Reply(m, fmt.Sprintf("Ошибка при определении типа гета, я не знаю тип <code>%v</code>.", get.Type))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			}
		} else {
			_, err := Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Send list of Gets to user on /getall
	Bot.Handle("/getall", func(m *tb.Message) {
		var getall []string
		var get Get
		result, _ := DB.Model(&Get{}).Rows()
		for result.Next() {
			err := DB.ScanRows(result, &get)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			getall = append(getall, get.Name)
		}
		_, err := Bot.Reply(m, fmt.Sprintf("Доступные геты: %v", strings.Join(getall[:], ", ")))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		return
	})
	//Save Get to DB on /set
	Bot.Handle("/set", func(m *tb.Message) {
		var get Get
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) < 3) || (m.ReplyTo != nil && len(text) != 2) {
			_, err := Bot.Reply(m, "Пример использования: <code>/set {гет} {значение}</code>\nИли отправь в ответ на какое-либо сообщение <code>/set {гет}</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		get.Name = strings.ToLower(text[1])
		if m.ReplyTo == nil && len(text) > 2 {
			get.Type = "Text"
			get.Data = strings.Join(text[2:], " ")
		} else if m.ReplyTo != nil && len(text) == 2 {
			get.Caption = m.ReplyTo.Caption
			switch {
			case m.ReplyTo.Animation != nil:
				get.Type = "Animation"
				get.Data = m.ReplyTo.Animation.FileID
			case m.ReplyTo.Audio != nil:
				get.Type = "Audio"
				get.Data = m.ReplyTo.Audio.FileID
			case m.ReplyTo.Photo != nil:
				get.Type = "Photo"
				get.Data = m.ReplyTo.Photo.FileID
			case m.ReplyTo.Video != nil:
				get.Type = "Video"
				get.Data = m.ReplyTo.Video.FileID
			case m.ReplyTo.Voice != nil:
				get.Type = "Voice"
				get.Data = m.ReplyTo.Voice.FileID
			case m.ReplyTo.Document != nil:
				get.Type = "Document"
				get.Data = m.ReplyTo.Document.FileID
			case m.ReplyTo.Text != "":
				get.Type = "Text"
				get.Data = m.ReplyTo.Text
			default:
				_, err := Bot.Reply(m, "Не удалось распознать файл в сообщении, возможно, он не поддерживается.")
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				return
			}
		}
		result := DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(get)
		if result.Error != nil {
			ErrorReporting(result.Error, m)
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось сохранить гет <code>%v</code>.", get.Name))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err := Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> сохранён как <code>%v</code>.", get.Name, get.Type))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Delete Get in DB on /del
	Bot.Handle("/del", func(m *tb.Message) {
		var text = strings.Split(m.Text, " ")
		if len(text) != 2 {
			_, err := Bot.Reply(m, "Пример использования: <code>/del {гет}</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		result := DB.Delete(&Get{Name: strings.ToLower(text[1])})
		if result.RowsAffected != 0 {
			_, err := Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> удалён.", text[1]))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			_, err := Bot.Reply(m, fmt.Sprintf("Гет <code>%v</code> не найден.", text[1]))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Send text in chat on /say
	Bot.Handle("/say", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var text = strings.Split(m.Text, " ")
		if len(text) > 1 {
			err := Bot.Delete(m)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			_, err = Bot.Send(m.Chat, strings.Join(text[1:], " "))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			_, err := Bot.Reply(m, "Укажите сообщение.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Send shrug in chat on /shrug
	Bot.Handle("/shrug", func(m *tb.Message) {
		_, err := Bot.Send(m.Chat, "¯\\_(ツ)_/¯")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Replace text in target message and send result on /sed
	Bot.Handle("/sed", func(m *tb.Message) {
		var text = strings.Split(m.Text, " ")
		if m.ReplyTo != nil {
			cmd := fmt.Sprintf("echo \"%v\" | sed \"%v\"", strings.ReplaceAll(m.ReplyTo.Text, "\"", "\\\""), strings.ReplaceAll(text[1], "\"", "\\\""))
			out, err := exec.Command("bash","-c", cmd).Output()
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			_, err = Bot.Reply(m, string(out))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			_, err := Bot.Reply(m, "Пример использования:\n/sed {патерн вида s/foo/bar/} в ответ на сообщение.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Send userid on /getid
	Bot.Handle("/getid", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		if m.ReplyTo != nil && m.ReplyTo.OriginalSender != nil {
			_, err := Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", m.ReplyTo.OriginalSender.FirstName, m.ReplyTo.OriginalSender.LastName, m.ReplyTo.OriginalSender.Username, m.ReplyTo.OriginalSender.ID))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else if m.ReplyTo != nil {
			_, err := Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", m.ReplyTo.Sender.FirstName, m.ReplyTo.Sender.LastName, m.ReplyTo.Sender.Username, m.ReplyTo.Sender.ID))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			_, err := Bot.Send(m.Sender, fmt.Sprintf("Firstname: %v\nLastname: %v\nUsername: %v\nUserID: %v", m.Sender.FirstName, m.Sender.LastName, m.Sender.Username, m.Sender.ID))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Reply "Pong!" on "ping"
	Bot.Handle("/ping", func(m *tb.Message) {
		_, err := Bot.Reply(m, "Pong!")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Reply "Polo!" on "marco"
	Bot.Handle("/marco", func(m *tb.Message) {
		_, err := Bot.Reply(m, "Polo!")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Reply currency "cur"
	Bot.Handle("/cur", func(m *tb.Message) {
		var target = *m
		var text = strings.Split(m.Text, " ")
		if len(text) != 4 {
			_, err := Bot.Reply(m, "Пример использования:\n/cur {количество} {EUR/USD/RUB} {EUR/USD/RUB}")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		if m.ReplyTo != nil {
			target = *m.ReplyTo
		}
		amount, err := strconv.ParseFloat(text[1], 64)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения количества:\n<code>%v</code>", err))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var symbol = strings.ToUpper(text[2])
		if !regexp.MustCompile(`^[A-Z]{3,4}$`).MatchString(symbol) {
			_, err := Bot.Reply(m, "Имя валюты должно состоять из 3-4 больших латинских символов.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var convert = strings.ToUpper(text[3])
		if !regexp.MustCompile(`^[A-Z]{3,4}$`).MatchString(convert) {
			_, err := Bot.Reply(m, "Имя валюты должно состоять из 3-4 больших латинских символов.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		client := cmc.NewClient(&cmc.Config{ProAPIKey: Config.CurrencyKey})
		conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount:  amount, Symbol:  symbol, Convert: convert})
		if err != nil {
			_, err := Bot.Reply(m, "Ошибка при запросе. Возможно, одна из валют не найдена.\nОнлайн-версия: https://coinmarketcap.com/ru/converter/", &tb.SendOptions{DisableWebPagePreview: true})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Reply(&target, fmt.Sprintf("%v %v = %v %v", conversion.Amount, conversion.Name, math.Round(conversion.Quote[convert].Price*100)/100, convert))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Reply google URL on "google"
	Bot.Handle("/google", func(m *tb.Message) {
		var target = *m
		var text = strings.Split(m.Text, " ")
		if len(text) == 1 {
			_, err := Bot.Reply(m, fmt.Sprintf("Пример использования:\n<code>/google {запрос}</code>"))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		if m.ReplyTo != nil {
			target = *m.ReplyTo
		}
		_, err := Bot.Reply(&target, fmt.Sprintf("https://www.google.com/search?q=%v", url.QueryEscape(strings.Join(text[1:], " "))))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Kick user on /kick
	Bot.Handle("/kick", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) == 1) || (m.ReplyTo != nil && len(text) != 2) {
			_, err := Bot.Reply(m, "Пример использования: <code>/kick {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/kick</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
		err = Bot.Ban(m.Chat, TargetChatMember)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		err = Bot.Unban(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка исключения пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> исключен.", target.ID, UserFullName(&target)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Ban user on /ban
	Bot.Handle("/ban", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) < 2) || (m.ReplyTo != nil && len(text) > 2) {
			_, err := Bot.Reply(m, "Пример использования: <code>/ban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/ban</code>\nЕсли нужно забанить на время, то добавь время в секундах через пробел.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, untildate, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя или время бана:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember.RestrictedUntil = untildate
		err = Bot.Ban(m.Chat, TargetChatMember)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка бана пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> забанен%v.", target.ID, UserFullName(&target), RestrictionTimeMessage(untildate)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Unban user on /unban
	Bot.Handle("/unban", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var target tb.User
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
			_, err := Bot.Reply(m, "Пример использования: <code>/unban {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unban</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		err = Bot.Unban(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка разбана пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Reply(m, fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> разбанен.", target.ID, UserFullName(&target)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Mute user on /mute
	Bot.Handle("/mute", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) < 2) || (m.ReplyTo != nil && len(text) > 2) {
			_, err := Bot.Reply(m, "Пример использования: <code>/mute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/mute</code>\nЕсли нужно замьютить на время, то добавь время в секундах через пробел.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, untildate, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя или время ограничения:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember.CanSendMessages = false
		TargetChatMember.RestrictedUntil = untildate
		err = Bot.Restrict(m.Chat, TargetChatMember)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка ограничения пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> больше не может отправлять сообщения%v.", target.ID, UserFullName(&target), RestrictionTimeMessage(untildate)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Unmute user on /unmute
	Bot.Handle("/unmute", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var target tb.User
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
			_, err := Bot.Reply(m, "Пример использования: <code>/unmute {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/unmute</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		TargetChatMember.CanSendMessages = true
		TargetChatMember.CanSendMedia = true
		TargetChatMember.CanSendPolls = true
		TargetChatMember.CanSendOther = true
		TargetChatMember.CanAddPreviews = true
		TargetChatMember.RestrictedUntil = time.Now().Unix() + 60
		err = Bot.Restrict(m.Chat, TargetChatMember)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка снятия ограничения пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Reply(m, fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a> снова может отправлять сообщения в чат.", target.ID, UserFullName(&target)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send formatted text on /me
	Bot.Handle("/me", func(m *tb.Message) {
		var text = strings.Split(m.Text, " ")
		if len(text) == 1 {
			_, err := Bot.Reply(m, fmt.Sprintf("Пример использования:\n<code>/me {делает что-то}</code>"))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		err := Bot.Delete(m)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Send(m.Chat, fmt.Sprintf("<code>%v %v</code>", UserFullName(m.Sender), strings.Join(text[1:], " ")))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Change chat name on /topic
	Bot.Handle("/topic", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var text = strings.Split(m.Text, " ")
		if len(text) < 2 {
			_, err := Bot.Reply(m, "Пример использования:\n<code>/topic {новая тема чата}</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		err := Bot.SetGroupTitle(m.Chat, fmt.Sprintf("Zavtrachat | %v", strings.Join(text[1:], " ")))
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка изменения названия чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}

			return
		}
	})
	//Write username on bonk picture and send to target
	Bot.Handle("/bonk", func(m *tb.Message) {
		if m.ReplyTo == nil {
			_, err := Bot.Reply(m, "Просто отправь <code>/bonk</code> в ответ на чье-либо сообщение.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return

		}
		var target = *m.ReplyTo
		im, err := webp.Load("files/bonk.webp")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		dc := gg.NewContextForImage(im)
		dc.DrawImage(im, 0, 0)
		dc.SetRGB(0,0,0)
		err = dc.LoadFontFace("files/impact.ttf", 20)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		dc.SetRGB(1, 1, 1)
		s := UserFullName(m.Sender)
		n := 4
		for dy := -n; dy <= n; dy++ {
			for dx := -n; dx <= n; dx++ {
				if dx*dx+dy*dy >= n*n {
					continue
				}
				x := 140 + float64(dx)
				y := 290 + float64(dy)
				dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
			}
		}
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(s, 140, 290, 0.5, 0.5)
		buf := new(bytes.Buffer)
		err = webp.Encode(buf, dc.Image(), nil)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Reply(&target, &tb.Sticker{File: tb.FromReader(buf)})
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Write username on hug picture and send to target
	Bot.Handle("/hug", func(m *tb.Message) {
		if m.ReplyTo == nil {
			_, err := Bot.Reply(m, "Просто отправь <code>/hug</code> в ответ на чье-либо сообщение.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return

		}
		var target = *m.ReplyTo
		im, err := webp.Load("files/hug.webp")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		dc := gg.NewContextForImage(im)
		dc.DrawImage(im, 0, 0)
		dc.Rotate(gg.Radians(15))
		dc.SetRGB(0,0,0)
		err = dc.LoadFontFace("files/impact.ttf", 20)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		dc.SetRGB(1, 1, 1)
		s := UserFullName(m.Sender)
		n := 4
		for dy := -n; dy <= n; dy++ {
			for dx := -n; dx <= n; dx++ {
				if dx*dx+dy*dy >= n*n {
					continue
				}
				x := 400 + float64(dx)
				y := -30 + float64(dy)
				dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
			}
		}
		dc.SetRGB(0, 0, 0)
		dc.DrawStringAnchored(s, 400, -30, 0.5, 0.5)
		buf := new(bytes.Buffer)
		err = webp.Encode(buf, dc.Image(), nil)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Reply(&target, &tb.Sticker{File: tb.FromReader(buf)})
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send slap message on /slap
	Bot.Handle("/slap", func(m *tb.Message) {
		var action = "дал леща"
		var target tb.User
		ChatMember, err := Bot.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		if ChatMember.CanRestrictMembers || ChatMember.Role == "creator" {
			action = "дал отцовского леща"
		}
		target, _, err = FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		_, err = Bot.Send(m.Chat, fmt.Sprintf("👋 <b>%v</b> %v %v", UserFullName(m.Sender), action, MentionUser(&target)))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send releases of 2 weeks on /releases
	Bot.Handle("/releases", func(m *tb.Message) {
		resp, err := http.Get(Config.ReleasesUrl)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		cal, err := ical.ParseCalendar(resp.Body)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		releases := ""
		today, _ := strconv.Atoi(time.Now().Format("20060102"))
		twoweeks, _ := strconv.Atoi(time.Now().AddDate(0, 0, 14).Format("20060102"))
		for _, element := range cal.Events() {
			date := element.GetProperty(ical.ComponentPropertyDtStart).Value
			name := element.GetProperty(ical.ComponentPropertySummary).Value
			dateint, _ := strconv.Atoi(date)
			if dateint > today && dateint < twoweeks {
				releases = fmt.Sprintf("<b>%v</b> - %v.%v.%v\n%v", strings.ReplaceAll(name, "\\,", ","), date[6:8], date[4:6], date[0:4], releases)
			}
		}
		_, err = Bot.Reply(m, releases)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send warning to user on /warn
	Bot.Handle("/warn", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var warn Warn
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
			_, err := Bot.Reply(m, "Пример использования: <code>/warn {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/warn</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		result := DB.First(&warn, target.ID)
		if result.RowsAffected != 0 {
			warn.Amount = warn.Amount - int(time.Now().Sub(warn.LastWarn).Hours() / 24 / 7)
			if warn.Amount < 0 {
				warn.Amount = 0
			}
			warn.Amount = warn.Amount + 1
		} else {
			warn.Amount = 1
		}
		warn.UserID = target.ID
		warn.LastWarn = time.Now()
		result = DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(warn)
		if result.Error != nil {
			ErrorReporting(result.Error, m)
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось выдать предупреждение:\n<code>%v</code>.", result.Error))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		if warn.Amount == 1 {
			_, err := Bot.Send(m.Chat, fmt.Sprintf("%v, у тебя 1 предупреждение.\nЕсль получишь 3 предупреждения за 2 недели, то будешь исключен из чата.", MentionUser(&target)))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
		if warn.Amount == 2 {
			_, err := Bot.Send(m.Chat, fmt.Sprintf("%v, у тебя 2 предупреждения.\nЕсли в течении недели получишь ещё одно, то будешь исключен из чата.", MentionUser(&target)))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
		if warn.Amount == 3 {
			untildate := time.Now().AddDate(0, 0, 7).Unix()
			TargetChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
			if err != nil {
				_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				return
			}
			TargetChatMember.RestrictedUntil = untildate
			err = Bot.Ban(m.Chat, TargetChatMember)
			if err != nil {
				_, err := Bot.Reply(m, fmt.Sprintf("Ошибка бана пользователя:\n<code>%v</code>", err.Error()))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				return
			}
			_, err = Bot.Reply(m, fmt.Sprintf("Пользователь <a href=\"tg://user?id=%v\">%v</a> забанен%v, т.к. набрал 3 предупреждения.", target.ID, UserFullName(&target), RestrictionTimeMessage(untildate)))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Send warning amount on /mywarns
	Bot.Handle("/mywarns", func(m *tb.Message) {
		var warn Warn
		result := DB.First(&warn, m.Sender.ID)
		if result.RowsAffected != 0 {
			warn.Amount = warn.Amount - int(time.Now().Sub(warn.LastWarn).Hours() / 24 / 7)
			if warn.Amount < 0 {
				warn.Amount = 0
			}
		} else {
			warn.UserID = m.Sender.ID
			warn.LastWarn = time.Unix(0, 0)
			warn.Amount = 0
		}
		warnStrings := []string{"предупреждений", "предупреждение", "предупреждения", "предупреждения"}
		_, err := Bot.Reply(m, fmt.Sprintf("У тебя %v %v.", warn.Amount, warnStrings[warn.Amount]))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send pidor rules on /pidorules
	Bot.Handle("/pidorules", func(m *tb.Message) {
		_, err := Bot.Reply(m, "Правила игры <b>Пидор Дня</b>:\n<b>1.</b> Зарегистрируйтесь в игру по команде /pidoreg\n<b>2.</b> Подождите пока зарегиструются все (или большинство :)\n<b>3.</b> Запустите розыгрыш по команде /pidor\n<b>4.</b> Просмотр статистики канала по команде /pidorstats, /pidorall\n<b>5.</b> Личная статистика по команде /pidorme\n<b>6. (!!! Только для администраторов чатов)</b>: удалить из игры может только Админ канала, сначала выведя по команде список игроков: /pidorlist (список упадёт в личку)\nУдалить же игрока можно по команде (используйте идентификатор пользователя - цифры из списка пользователей): /pidordel {ID или никнейм юзера}\nТак же, удалить можно просто отправив /pidordel в ответ на сообщение пользователя, которого нужно удалить из игры.\n\nВажно, розыгрыш проходит только раз в день, повторная команда выведет <b>результат</b> игры.\n\nСброс розыгрыша происходит каждый день ночью.\n\nПоддержать автора оригинального бота можно по <a href=\"https://www.paypal.me/unicott/2\">ссылке</a> :)")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send DB result on /pidoreg
	Bot.Handle("/pidoreg", func(m *tb.Message) {
		var pidor PidorList
		result := DB.First(&pidor, m.Sender.ID)
		if result.RowsAffected != 0 {
			_, err := Bot.Reply(m, "Эй, ты уже в игре!")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			pidor = PidorList(*m.Sender)
			result = DB.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(pidor)
			if result.Error != nil {
				ErrorReporting(result.Error, m)
				_, err := Bot.Reply(m, fmt.Sprintf("Не удалось зарегистрироваться:\n<code>%v</code>.", result.Error))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				return
			}
			_, err := Bot.Reply(m, "OK! Ты теперь участвуешь в игре <b>Пидор Дня</b>!")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Send DB stats on /pidorme
	Bot.Handle("/pidorme", func(m *tb.Message) {
		var pidor PidorStats
		var countYear int64
		var countAlltime int64
		pidor.UserID = m.Sender.ID
		DB.Model(&PidorStats{}).Where(pidor).Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(),1,1,0,0,0,0,time.Local), time.Now()).Count(&countYear)
		DB.Model(&PidorStats{}).Where(pidor).Count(&countAlltime)
		_, err := Bot.Reply(m, fmt.Sprintf("В этом году ты был пидором дня — %v раз!\nЗа всё время ты был пидором дня — %v раз!", countYear, countAlltime))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Remove user in DB on /pidordel
	Bot.Handle("/pidordel", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var user tb.User
		var pidor PidorList
		user, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		pidor = PidorList(user)
		result := DB.Delete(&pidor)
		if result.RowsAffected != 0 {
			_, err := Bot.Reply(m, fmt.Sprintf("Пользователь %v удалён из игры <b>Пидор Дня</b>!", MentionUser(&user)))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось удалить пользователя:\n<code>%v</code>", result.Error.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//List add pidors from DB on /pidorlist
	Bot.Handle("/pidorlist", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var pidorlist string
		var pidor PidorList
		var i = 0
		result, _ := DB.Model(&PidorList{}).Rows()
		for result.Next() {
			err := DB.ScanRows(result, &pidor)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			i++
			pidorlist += strconv.Itoa(i) + ". @" + pidor.Username + " (" + strconv.Itoa(pidor.ID) + ")\n"
			if len(pidorlist) > 3900 {
				_, err = Bot.Send(m.Sender, pidorlist)
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				pidorlist = ""
			}
		}
		_, err := Bot.Send(m.Sender, pidorlist)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Reply(m, "Список отправлен в личку.\nЕсли список не пришел, то убедитесь, что бот запущен и не заблокирован в личке.")
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		return
	})
	//Send top 10 pidors of all time on /pidorall
	Bot.Handle("/pidorall", func(m *tb.Message) {
		var i = 0
		var username string
		var count int64
		var pidorall = "Топ-10 пидоров за всё время:\n\n"
		result, _ := DB.Select("username, COUNT(*) as count").Table("pidor_stats, pidor_lists").Where("pidor_stats.user_id=pidor_lists.id").Group("user_id").Order("count DESC").Limit(10).Rows()
		for result.Next() {
			err := result.Scan(&username, &count)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			i++
			pidorall += fmt.Sprintf("%v. %v - %v раз(а)\n", i, username, count)
		}
		DB.Model(PidorList{}).Count(&count)
		pidorall += fmt.Sprintf("\nВсего участников — %v", count)
		_, err := Bot.Reply(m, pidorall)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send top 10 pidors of year on /pidorstats
	Bot.Handle("/pidorstats", func(m *tb.Message) {
		var text = strings.Split(m.Text, " ")
		var i = 0
		var year = time.Now().Year()
		var username string
		var count int64
		if len(text) == 2 {
			argYear, err := strconv.Atoi(text[1])
			if err != nil {
				_, err := Bot.Reply(m, "Ошибка определения года.\nУкажите год с 2019 по предыдущий.")
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				return
			}
			if argYear < year && argYear > 2018 {
				year = argYear
			}
		}
		var pidorall = "Топ-10 пидоров за " + strconv.Itoa(year) + " год:\n\n"
		result, _ := DB.Select("username, COUNT(*) as count").Table("pidor_stats, pidor_lists").Where("pidor_stats.user_id=pidor_lists.id").Where("date BETWEEN ? AND ?", time.Date(year,1,1,0,0,0,0,time.Local), time.Date(year+1,1,1,0,0,0,0,time.Local)).Group("user_id").Order("count DESC").Limit(10).Rows()
		for result.Next() {
			err := result.Scan(&username, &count)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			i++
			pidorall += fmt.Sprintf("%v. %v - %v раз(а)\n", i, username, count)
		}
		DB.Model(PidorList{}).Count(&count)
		pidorall += fmt.Sprintf("\nВсего участников — %v", count)
		_, err := Bot.Reply(m, pidorall)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Launch pidorday game
	Bot.Handle("/pidor", func(m *tb.Message) {
		if busy["pidor"] {
			_, err := Bot.Reply(m, "Команда занята. Попробуйте позже.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		busy["pidor"] = true
		defer func() {busy["pidor"] = false}()
		var pidor PidorStats
		var pidorToday PidorList
		result := DB.Model(PidorStats{}).Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local), time.Now()).First(&pidor)
		if result.RowsAffected == 0 {
			DB.Model(PidorList{}).Order("RANDOM()").First(&pidorToday)
			TargetChatMember, err := Bot.ChatMemberOf(m.Chat, &tb.User{ID: pidorToday.ID})
			if err != nil {
				_, err := Bot.Reply(m, fmt.Sprintf("Я нашел пидора дня, но похоже, что с <a href=\"tg://user?id=%v\">%v</a> что-то не так, так что попробуйте еще раз, пока я удаляю его из игры! Ошибка:\n<code>%v</code>", pidorToday.ID, pidorToday.Username, err.Error()))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				DB.Delete(pidorToday)
				return
			}
			if TargetChatMember.Role == "left" {
				_, err := Bot.Reply(m, fmt.Sprintf("Я нашел пидора дня, но похоже, что <a href=\"tg://user?id=%v\">%v</a> вышел из этого чата (вот пидор!), так что попробуйте еще раз, пока я удаляю его из игры!", pidorToday.ID, pidorToday.Username))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				DB.Delete(pidorToday)
				return
			}
			if TargetChatMember.Role == "kicked" {
				_, err := Bot.Reply(m, fmt.Sprintf("Я нашел пидора дня, но похоже, что <a href=\"tg://user?id=%v\">%v</a> был забанен в этом чате (получил пидор!), так что попробуйте еще раз, пока я удаляю его из игры!", pidorToday.ID, pidorToday.Username))
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				DB.Delete(pidorToday)
				return
			}
			DB.Create(pidorToday)
			messages := [][]string{
				{"Инициирую поиск пидора дня...", "Опять в эти ваши игрульки играете? Ну ладно...", "Woop-woop! That's the sound of da pidor-police!", "Система взломана. Нанесён урон. Запущено планирование контрмер.", "Сейчас поколдуем...", "Инициирую поиск пидора дня...", "Зачем вы меня разбудили...", "Кто сегодня счастливчик?"},
				{"Хм...", "Сканирую...", "Ведётся поиск в базе данных", "Сонно смотрит на бумаги", "(Ворчит) А могли бы на работе делом заниматься", "Военный спутник запущен, коды доступа внутри...", "Ну давай, посмотрим кто тут классный..."},
				{"Высокий приоритет мобильному юниту.", "Ох...", "Ого-го...", "Так, что тут у нас?", "В этом совершенно нет смысла...", "Что с нами стало...", "Тысяча чертей!", "Ведётся захват подозреваемого..."},
				{"Стоять! Не двигаться! Ты объявлен пидором дня, ", "Ого, вы посмотрите только! А пидор дня то - ", "Пидор дня обыкновенный, 1шт. - ", ".∧＿∧ \n( ･ω･｡)つ━☆・*。 \n⊂  ノ    ・゜+. \nしーＪ   °。+ *´¨) \n         .· ´¸.·*´¨) \n          (¸.·´ (¸.·'* ☆ ВЖУХ И ТЫ ПИДОР, ", "Ага! Поздравляю! Сегодня ты пидор - ", "Кажется, пидор дня - ", "Анализ завершен. Ты пидор, "},
			}
			for i := 0; i <= 3; i++ {
				duration := time.Second * time.Duration(i * 2)
				message := messages[i][RandInt(0, len(messages[i])-1)]
				if i == 3 {
					message += fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a>", pidorToday.ID, pidorToday.Username)
				}
				go func() {
					time.Sleep(duration)
					_, err := Bot.Send(m.Chat, message)
					if err != nil {
						ErrorReporting(err, m)
						return
					}
				}()
			}
		} else {
			DB.Model(PidorList{}).Where(pidor.UserID).First(&pidorToday)
			_, err := Bot.Reply(m, fmt.Sprintf("Согласно моей информации, по результатам сегодняшнего розыгрыша пидор дня - %v!", pidorToday.Username))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})
	//Kill user on /blessing, /suicide
	Bot.Handle("/blessing", func(m *tb.Message) {
		err := Bot.Delete(m)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		ChatMember, err := Bot.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
			_, err := Bot.Reply(m, fmt.Sprintf("<code>👻 %v возродился у костра.</code>", UserFullName(m.Sender)))
			if err != nil {
				ErrorReporting(err, m)
			}
			return
		}
		var duelist Duelist
		result := DB.Model(Duelist{}).Where(m.Sender.ID).First(&duelist)
		if result.RowsAffected == 0 {
			duelist.UserID = m.Sender.ID
			duelist.Kills = 0
			duelist.Deaths = 0
		}
		duelist.Deaths++
		result = DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(duelist)
		if result.Error != nil {
			ErrorReporting(result.Error, m)
			return
		}
		ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(600*duelist.Deaths)).Unix()
		err = Bot.Restrict(m.Chat, ChatMember)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Send(m.Chat, fmt.Sprintf("<code>💥 %v выбрал лёгкий путь.\nРеспавн через %v0 минут.</code>", UserFullName(m.Sender), duelist.Deaths))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	Bot.Handle("/suicide", func(m *tb.Message) {
		err := Bot.Delete(m)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		ChatMember, err := Bot.ChatMemberOf(m.Chat, m.Sender)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
			_, err := Bot.Reply(m, fmt.Sprintf("<code>👻 %v возродился у костра.</code>", UserFullName(m.Sender)))
			if err != nil {
				ErrorReporting(err, m)
			}
			return
		}
		var duelist Duelist
		result := DB.Model(Duelist{}).Where(m.Sender.ID).First(&duelist)
		if result.RowsAffected == 0 {
			duelist.UserID = m.Sender.ID
			duelist.Kills = 0
			duelist.Deaths = 0
		}
		duelist.Deaths++
		result = DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(duelist)
		if result.Error != nil {
			ErrorReporting(result.Error, m)
			return
		}
		ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(600*duelist.Deaths)).Unix()
		err = Bot.Restrict(m.Chat, ChatMember)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Send(m.Chat, fmt.Sprintf("<code>💥 %v выбрал лёгкий путь.\nРеспавн через %v0 минут.</code>", UserFullName(m.Sender), duelist.Deaths))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Kill user on /kill
	Bot.Handle("/kill", func(m *tb.Message) {
		if !StringInSlice(m.Sender.Username, Config.Telegram.Admins) && !StringInSlice(m.Sender.Username, Config.Telegram.Moders) {
			_, err := Bot.Reply(m, &tb.Animation{File: tb.File{FileID: "CgACAgIAAx0CQvXPNQABHGrDYIBIvDLiVV6ZMPypWMi_NVDkoFQAAq4LAAIwqQlIQT82LRwIpmoeBA"}})
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
			_, err := Bot.Reply(m, "Пример использования: <code>/kill {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/kill</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		ChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		var duelist Duelist
		result := DB.Model(Duelist{}).Where(target.ID).First(&duelist)
		if result.RowsAffected == 0 {
			duelist.UserID = target.ID
			duelist.Kills = 0
			duelist.Deaths = 0
		}
		duelist.Deaths++
		result = DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(duelist)
		if result.Error != nil {
			ErrorReporting(result.Error, m)
			return
		}
		ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(600*duelist.Deaths)).Unix()
		err = Bot.Restrict(m.Chat, ChatMember)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		_, err = Bot.Send(m.Chat, fmt.Sprintf("💥 %v пристрелил %v.\n%v отправился на респавн на %v0 минут.", UserFullName(m.Sender), UserFullName(&target), UserFullName(&target), duelist.Deaths))
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	//Send user duelist stats on /duelstats
	Bot.Handle("/duelstats", func(m *tb.Message) {
		var duelist Duelist
		result := DB.Model(Duelist{}).Where(m.Sender.ID).First(&duelist)
		if result.RowsAffected == 0 {
			_, err := Bot.Reply(m, "У тебя нет статистики.")
			if err != nil {
				ErrorReporting(err, m)
			}
			return
		}
		_, err := Bot.Reply(m, fmt.Sprintf("Побед: %v\nСмертей: %v", duelist.Kills, duelist.Deaths))
		if err != nil {
			ErrorReporting(err, m)
		}
	})
	//Russianroulette game
	var russianrouletteMessage *tb.Message
	russianrouletteSelector := tb.ReplyMarkup{}
	russianrouletteAcceptButton := russianrouletteSelector.Data("👍 Принять вызов", "russianroulette_accept")
	russianrouletteDenyButton := russianrouletteSelector.Data("👎 Бежать с позором", "russianroulette_deny")
	russianrouletteSelector.Inline(
		russianrouletteSelector.Row(russianrouletteAcceptButton, russianrouletteDenyButton),
	)
	Bot.Handle("/russianroulette", func(m *tb.Message) {
		if russianrouletteMessage == nil {
			russianrouletteMessage = m
			russianrouletteMessage.Unixtime = 0
		}
		if busy["bot_is_dead"] {
			if time.Now().Unix() - russianrouletteMessage.Time().Unix() > 3600 {
				busy["bot_is_dead"] = false
			} else {
				_, err := Bot.Reply(m, "Я не могу провести игру, т.к. я немного умер. Зайдите позже.")
				if err != nil {
					ErrorReporting(err, m)
					return
				}
				return
			}
		}
		if busy["russianroulettePending"] && !busy["russianrouletteInProgress"] && time.Now().Unix() - russianrouletteMessage.Time().Unix() > 60 {
			busy["russianroulette"] = false
			busy["russianroulettePending"] = false
			busy["russianrouletteInProgress"] = false
			_, err := Bot.Edit(russianrouletteMessage, fmt.Sprintf("%v не пришел на дуэль.", UserFullName(russianrouletteMessage.Entities[0].User)))
			if err != nil {
				ErrorReporting(err, russianrouletteMessage)
				return
			}
		}
		if busy["russianrouletteInProgress"] && time.Now().Unix() - russianrouletteMessage.Time().Unix() > 120 {
			busy["russianroulette"] = false
			busy["russianroulettePending"] = false
			busy["russianrouletteInProgress"] = false
		}
		if busy["russianroulette"] || busy["russianroulettePending"] || busy["russianrouletteInProgress"]  {
			_, err := Bot.Reply(m, "Команда занята. Попробуйте позже.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		busy["russianroulette"] = true
		defer func() {busy["russianroulette"] = false}()
		var text = strings.Split(m.Text, " ")
		if (m.ReplyTo == nil && len(text) != 2) || (m.ReplyTo != nil && len(text) != 1) {
			_, err := Bot.Reply(m, "Пример использования: <code>/russianroulette {ID или никнейм}</code>\nИли отправь в ответ на какое-либо сообщение <code>/russianroulette</code>")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		target, _, err := FindUserInMessage(*m)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Не удалось определить пользователя:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		if target.ID == m.Sender.ID {
			_, err := Bot.Reply(m, "Как ты себе это представляешь? Нет, нельзя вызвать на дуэль самого себя.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		if target.IsBot {
			_, err := Bot.Reply(m, "Бота нельзя вызвать на дуэль.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		ChatMember, err := Bot.ChatMemberOf(m.Chat, &target)
		if err != nil {
			_, err := Bot.Reply(m, fmt.Sprintf("Ошибка определения пользователя чата:\n<code>%v</code>", err.Error()))
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		log.Println(ChatMember)
		if false {
			_, err := Bot.Reply(m, "Нельзя вызвать на дуэль мертвеца.")
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			return
		}
		err = Bot.Delete(m)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		russianrouletteMessage, err = Bot.Send(m.Chat, fmt.Sprintf("%v! %v вызывает тебя на дуэль!", MentionUser(&target), MentionUser(m.Sender)), &russianrouletteSelector)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		busy["russianroulettePending"] = true
	})
	Bot.Handle(&russianrouletteAcceptButton, func(c *tb.Callback) {
		err := Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		message := c.Message
		victim := c.Message.Entities[0].User
		if victim.ID != c.Sender.ID {
			err := Bot.Respond(c, &tb.CallbackResponse{})
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			return
		}
		player := c.Message.Entities[1].User
		busy["russianroulette"] = false
		busy["russianroulettePending"] = false
		busy["russianrouletteInProgress"] = true
		defer func() {busy["russianrouletteInProgress"] = false}()
		success := []string{"%v остаётся в живых. Хм... может порох отсырел?", "В воздухе повисла тишина. %v остаётся в живых.", "%v сегодня заново родился.", "%v остаётся в живых. Хм... я ведь зарядил его?", "%v остаётся в живых. Прикольно, а давай проверим на ком-нибудь другом?"}
		invincible := []string{"пуля отскочила от головы %v и улетела в другой чат.", "%v похмурил брови и отклеил расплющенную пулю со своей головы.", "но ничего не произошло. %v взглянул на револьвер, он был неисправен.", "пуля прошла навылет, но не оставила каких-либо следов на %v."}
		fail := []string{"мозги %v разлетелись по чату!", "%v упал со стула и его кровь растеклась по месседжу.", "%v замер и спустя секунду упал на стол.", "пуля едва не задела кого-то из участников чата! А? Что? А, %v мёртв, да.", "и в воздухе повисла тишина. Все начали оглядываться, когда %v уже был мёртв."}
		prefix := fmt.Sprintf("Дуэль! %v против %v!\n", MentionUser(player), MentionUser(victim))
		_, err = Bot.Edit(message, fmt.Sprintf("%vЗаряжаю один патрон в револьвер и прокручиваю барабан.", prefix), &tb.SendOptions{ReplyMarkup: nil})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		time.Sleep(time.Second * 2)
		_, err = Bot.Edit(message, fmt.Sprintf("%vКладу револьвер на стол и раскручиваю его.", prefix))
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		time.Sleep(time.Second * 2)
		if RandInt(1,360)%2 == 0 {
			player, victim = victim, player
		}
		_, err = Bot.Edit(message, fmt.Sprintf("%vРевольвер останавливается на %v, первый ход за ним.", prefix, MentionUser(victim)))
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		bullet := RandInt(1,6)
		for i := 1; i <= bullet; i++ {
			time.Sleep(time.Second * 2)
			prefix = fmt.Sprintf("Дуэль! %v против %v, раунд %v:\n%v берёт револьвер, приставляет его к голове и...\n", MentionUser(player), MentionUser(victim), i, MentionUser(victim))
			_, err := Bot.Edit(message, prefix)
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			if bullet != i {
				time.Sleep(time.Second * 2)
				_, err := Bot.Edit(message, fmt.Sprintf("%v🍾 %v", prefix, fmt.Sprintf(success[RandInt(0, len(success)-1)], MentionUser(victim))))
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				player, victim = victim, player
			}
		}
		time.Sleep(time.Second * 2)
		PlayerChatMember, err := Bot.ChatMemberOf(c.Message.Chat, player)
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		VictimChatMember, err := Bot.ChatMemberOf(c.Message.Chat, victim)
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		if (PlayerChatMember.Role == "creator" || PlayerChatMember.Role == "administrator") && (VictimChatMember.Role == "creator" || VictimChatMember.Role == "administrator") {
			_, err = Bot.Edit(message, fmt.Sprintf("%vПуля отскакивает от головы %v и летит в голову %v.", prefix, MentionUser(victim), MentionUser(player)))
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			time.Sleep(time.Second * 2)
			_, err = Bot.Edit(message, fmt.Sprintf("%vПуля отскакивает от головы %v и летит в голову %v.", prefix, MentionUser(player), MentionUser(victim)))
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			time.Sleep(time.Second * 2)
			_, err = Bot.Edit(message, fmt.Sprintf("%vПуля отскакивает от головы %v и летит в мою голову... блять.", prefix, MentionUser(victim)))
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			busy["bot_is_dead"] = true
			return
		}
		if StringInSlice(victim.Username, Config.Telegram.Admins) {
			_, err = Bot.Edit(message, fmt.Sprintf("%v😈 Наводит револьвер на %v и стреляет.", prefix, MentionUser(player)))
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			time.Sleep(time.Second * 3)
			var duelist Duelist
			result := DB.Model(Duelist{}).Where(player.ID).First(&duelist)
			if result.RowsAffected == 0 {
				duelist.UserID = player.ID
				duelist.Kills = 0
				duelist.Deaths = 0
			}
			duelist.Deaths++
			result = DB.Clauses(clause.OnConflict{
				UpdateAll: true,
			}).Create(duelist)
			if result.Error != nil {
				ErrorReporting(result.Error, c.Message)
				return
			}
			PlayerChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(600*duelist.Deaths)).Unix()
			err = Bot.Restrict(c.Message.Chat, PlayerChatMember)
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			_, err = Bot.Edit(message, fmt.Sprintf("%v😈 Наводит револьвер на %v и стреляет.\nЯ хз как это объяснить, но %v победитель!\n%v отправился на респавн на %v0 минут.", prefix, MentionUser(player), MentionUser(victim), MentionUser(player), duelist.Deaths))
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			return
		}
		if VictimChatMember.Role == "creator" || VictimChatMember.Role == "administrator" {
			prefix = fmt.Sprintf("%v💥 %v", prefix, fmt.Sprintf(invincible[RandInt(0, len(invincible)-1)], MentionUser(victim)))
			_, err := Bot.Edit(message, prefix)
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			time.Sleep(time.Second * 2)
			_, err = Bot.Edit(message, fmt.Sprintf("%v\nПохоже, у нас ничья.", prefix))
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			return
		}
		prefix = fmt.Sprintf("%v💥 %v", prefix, fmt.Sprintf(fail[RandInt(0, len(fail)-1)], MentionUser(victim)))
		_, err = Bot.Edit(message, prefix)
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		time.Sleep(time.Second * 2)
		var VictimDuelist Duelist
		result := DB.Model(Duelist{}).Where(victim.ID).First(&VictimDuelist)
		if result.RowsAffected == 0 {
			VictimDuelist.UserID = victim.ID
			VictimDuelist.Kills = 0
			VictimDuelist.Deaths = 0
		}
		VictimDuelist.Deaths++
		result = DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(VictimDuelist)
		if result.Error != nil {
			ErrorReporting(result.Error, c.Message)
			return
		}
		VictimChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(600*VictimDuelist.Deaths)).Unix()
		err = Bot.Restrict(c.Message.Chat, VictimChatMember)
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		_, err = Bot.Edit(message, fmt.Sprintf("%v\nПобедитель дуэли: %v.\n%v отправился на респавн на %v0 минут.", prefix, MentionUser(player), MentionUser(victim), VictimDuelist.Deaths))
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		var PlayerDuelist Duelist
		result = DB.Model(Duelist{}).Where(victim.ID).First(&PlayerDuelist)
		if result.RowsAffected == 0 {
			PlayerDuelist.UserID = victim.ID
			PlayerDuelist.Kills = 0
			PlayerDuelist.Deaths = 0
		}
		PlayerDuelist.Kills++
		result = DB.Clauses(clause.OnConflict{
			UpdateAll: true,
		}).Create(PlayerDuelist)
		if result.Error != nil {
			ErrorReporting(result.Error, c.Message)
			return
		}
	})
	Bot.Handle(&russianrouletteDenyButton, func(c *tb.Callback) {
		err := Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
		victim := c.Message.Entities[0].User
		if victim.ID != c.Sender.ID {
			err := Bot.Respond(c, &tb.CallbackResponse{})
			if err != nil {
				ErrorReporting(err, c.Message)
				return
			}
			return
		}
		busy["russianroulette"] = false
		busy["russianroulettePending"] = false
		busy["russianrouletteInProgress"] = false
		_, err = Bot.Edit(c.Message, fmt.Sprintf("%v отказался от дуэли.", UserFullName(c.Sender)))
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
	})

	//Gather user data on incoming text message
	Bot.Handle(tb.OnText, func(m *tb.Message) {
		err := GatherData(m.Sender)
		if err != nil {
			ErrorReporting(err, m)
		}
	})
	//Repost channel post to chat
	Bot.Handle(tb.OnChannelPost, func(m *tb.Message) {
		if m.Chat.Username == Config.Telegram.Channel {
			chat, err := Bot.ChatByID("@"+Config.Telegram.Chat)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
			_, err = Bot.Forward(chat, m)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
	})

	//User join
	var welcomeMessage *tb.Message
	welcomeSelector := tb.ReplyMarkup{}
	welcomeFirstWrongButton := welcomeSelector.Data("Джабир, Латиф и Хиляль", "Button"+strconv.Itoa(RandInt(10000,99999)))
	welcomeRightButton := welcomeSelector.Data("Дмитрий, Тимур и Максим", "Button"+strconv.Itoa(RandInt(10000,99999)))
	welcomeSecondWrongButton := welcomeSelector.Data("Бубылда, Чингачгук и Гавкошмыг", "Button"+strconv.Itoa(RandInt(10000,99999)))
	welcomeThirdWrongButton := welcomeSelector.Data("Мандарин, Оладушек и Эчпочмак", "Button"+strconv.Itoa(RandInt(10000,99999)))
	buttons := []tb.Btn {welcomeRightButton, welcomeFirstWrongButton, welcomeSecondWrongButton, welcomeThirdWrongButton}
	pseudorand.Seed(time.Now().UnixNano())
	pseudorand.Shuffle(len(buttons), func(i, j int) {
		buttons[i], buttons[j] = buttons[j], buttons[i]
	})
	welcomeSelector.Inline(
		welcomeSelector.Row(buttons[0], buttons[1]),
		welcomeSelector.Row(buttons[2], buttons[3]),
	)
	nopes := []string{"неа", "не", "нет", "не то", "не попал"}
	arab, err := regexp.Compile("[\u0600-\u06ff]|[\u0750-\u077f]|[\ufb50-\ufbc1]|[\ufbd3-\ufd3f]|[\ufd50-\ufd8f]|[\ufd92-\ufdc7]|[\ufe70-\ufefc]|[\uFDF0-\uFDFD]")
	if err != nil {
		log.Fatal(err)
		return
	}
	Bot.Handle(tb.OnUserJoined, func(m *tb.Message) {
		if welcomeMessage == nil {
			welcomeMessage = m
			welcomeMessage.Unixtime = 0
		}
		if m.Chat.Username != Config.Telegram.Chat {
			return
		}
		err := Bot.Delete(m)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		log.Printf("New user detected in %v (%v)! ID: %v. Login: %v. Name: %v.", m.Chat.Title, m.Chat.ID, m.Sender.ID, UserName(m.Sender), UserFullName(m.Sender))
		User := m.Sender
		Chat := m.Chat
		ChatMember, err := Bot.ChatMemberOf(Chat, User)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		ChatMember.CanSendMessages = false
		err = Bot.Restrict(Chat, ChatMember)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
		if arab.MatchString(UserFullName(User)) || User.FirstName == "ICSM" {
			err = Bot.Ban(Chat, ChatMember)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
		var httpClient = &http.Client{Timeout: 10 * time.Second}
		httpResponse, err := httpClient.Get(fmt.Sprintf("https://api.cas.chat/check?user_id=%v", User.ID))
		if err != nil {
			ErrorReporting(err, m)
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
			ErrorReporting(err, m)
			return
		}
		if fastjson.GetBool(jsonBytes, "ok") {
			err = Bot.Ban(Chat, ChatMember)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
		if time.Now().Unix() - welcomeMessage.Time().Unix() > 10 {
			welcomeMessage, err = Bot.Send(Chat, fmt.Sprintf("Добро пожаловать, %v!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста?", MentionUser(User)), &welcomeSelector)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		} else {
			text := "Добро пожаловать"
			for _, element := range welcomeMessage.Entities {
				text += ", " + MentionUser(element.User)
			}
			text += ", " + MentionUser(m.Sender) + "!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста? У тебя 2 минуты."
			_, err = Bot.Edit(welcomeMessage, text, &welcomeSelector)
			if err != nil {
				ErrorReporting(err, m)
				return
			}
		}
		go func() {
			time.Sleep(time.Second * 120)
			ChatMember, err := Bot.ChatMemberOf(m.Chat, m.Sender)
			if err != nil {
				return
			}
			if ChatMember.Role != "member" {
				err := Bot.Ban(m.Chat, &tb.ChatMember{User: m.Sender})
				if err != nil {
					ErrorReporting(err, m)
					return
				}
			}
			err = Bot.Delete(m)
			if err != nil {
				return
			}
		}()
	})
	Bot.Handle(tb.OnUserLeft, func(m *tb.Message) {
		err := Bot.Delete(m)
		if err != nil {
			ErrorReporting(err, m)
			return
		}
	})
	Bot.Handle(&welcomeRightButton, func(c *tb.Callback) {
		for _, element := range c.Message.Entities {
			if element.User.ID == c.Sender.ID {
				err = Bot.Respond(c, &tb.CallbackResponse{Text: fmt.Sprintf("Добро пожаловать, %v!\nТеперь у тебя есть доступ к чату.", UserFullName(c.Sender)), ShowAlert: true})
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				ChatMember, err := Bot.ChatMemberOf(c.Message.Chat, c.Sender)
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				ChatMember.CanSendMessages = true
				ChatMember.RestrictedUntil = time.Now().Add(time.Hour).Unix()
				err = Bot.Promote(c.Message.Chat, ChatMember)
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				if len(c.Message.Entities) == 1 {
					if welcomeMessage.ID == c.Message.ID {
						welcomeMessage.Unixtime = 0
					}
					err = Bot.Delete(c.Message)
					if err != nil {
						ErrorReporting(err, c.Message)
						return
					}
				} else {
					text := "Добро пожаловать"
					for _, element := range c.Message.Entities {
						if element.User.ID != c.Sender.ID {
							text += ", " + MentionUser(element.User)
						}
					}
					text += "!\nЧтобы получить доступ в чат, ответь на вопрос.\nКак зовут ведущих подкаста?"
					_, err = Bot.Edit(c.Message, text, &welcomeSelector)
					if err != nil {
						ErrorReporting(err, c.Message)
						return
					}
				}
				return
			}
		}
		err := Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
	})
	Bot.Handle(&welcomeFirstWrongButton, func(c *tb.Callback) {
		for _, element := range c.Message.Entities {
			if element.User.ID == c.Sender.ID {
				err := Bot.Respond(c, &tb.CallbackResponse{Text: nopes[RandInt(0, len(nopes)-1)]})
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				return
			}
		}
		err := Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
	})
	Bot.Handle(&welcomeSecondWrongButton, func(c *tb.Callback) {
		for _, element := range c.Message.Entities {
			if element.User.ID == c.Sender.ID {
				err := Bot.Respond(c, &tb.CallbackResponse{Text: nopes[RandInt(0, len(nopes)-1)]})
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				return
			}
		}
		err := Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
	})
	Bot.Handle(&welcomeThirdWrongButton, func(c *tb.Callback) {
		for _, element := range c.Message.Entities {
			if element.User.ID == c.Sender.ID {
				err := Bot.Respond(c, &tb.CallbackResponse{Text: nopes[RandInt(0, len(nopes)-1)]})
				if err != nil {
					ErrorReporting(err, c.Message)
					return
				}
				return
			}
		}
		err := Bot.Respond(c, &tb.CallbackResponse{})
		if err != nil {
			ErrorReporting(err, c.Message)
			return
		}
	})

	//ZavtraStreamCheck Loop
	go func() {
		for {
			delay := 240
			if time.Now().Hour() < 24 && time.Now().Hour() >= 18 {
				delay = 30
			}
			time.Sleep(time.Duration(delay) * time.Second)
			err := ZavtraStreamCheck("youtube")
			if err != nil {
				log.Println(err.Error())
				chat, _ := Bot.ChatByID("@"+Config.Telegram.SysAdmin)
				_, _ = Bot.Send(chat, fmt.Sprintf("ZavtraStreamCheck error:\n<code>%v</code>", err.Error()))
			}
		}
	}()

	Bot.Start()
}