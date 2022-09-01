package utils

import (
	tdctx "context"
	"encoding/json"
	"fmt"
	"html"
	"log"
	"runtime"
	"strconv"
	"strings"
	"time"

	tele "github.com/NexonSU/telebot"
	"github.com/gotd/contrib/bg"
	"github.com/gotd/td/telegram"
	"gorm.io/gorm/clause"
)

var Bot tele.Bot
var GotdClient *telegram.Client
var GotdContext tdctx.Context

// BotInit initializes Telegram Bot
// Moved from auto init to manual init to make the code in utils package testable
func BotInit() {
	if Config.Token == "" {
		log.Fatal("Telegram bot token not found in config.json")
	}
	if Config.Chat == 0 {
		log.Fatal("Chat username not found in config.json")
	}
	Settings := tele.Settings{
		URL:       Config.BotApiUrl,
		Token:     Config.Token,
		ParseMode: tele.ModeHTML,
		Poller: &tele.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.AllowedUpdates,
		},
	}
	if Config.EndpointPublicURL != "" || Config.Listen != "" {
		Settings.Poller = &tele.Webhook{
			Listen: Config.Listen,
			Endpoint: &tele.WebhookEndpoint{
				PublicURL: Config.EndpointPublicURL,
			},
			MaxConnections: Config.MaxConnections,
			AllowedUpdates: Config.AllowedUpdates,
		}
	} else {
		Settings.Poller = &tele.LongPoller{
			Timeout:        10 * time.Second,
			AllowedUpdates: Config.AllowedUpdates,
		}
	}
	var bot, err = tele.NewBot(Settings)
	if err != nil {
		log.Println(Config.BotApiUrl)
		log.Fatal(err)
	}
	if Config.SysAdmin != 0 {
		bot.Send(tele.ChatID(Config.SysAdmin), fmt.Sprintf("%v has finished starting up.", MentionUser(bot.Me)))
	}

	Bot = *bot
}

func gotdClientInit() error {
	if Config.AppID == 0 || Config.AppHash == "" {
		return nil
	}
	client := telegram.NewClient(Config.AppID, Config.AppHash, telegram.Options{})
	return client.Run(tdctx.Background(), func(ctx tdctx.Context) error {
		stop, err := bg.Connect(client)
		if err != nil {
			return err
		}
		defer func() { _ = stop() }()

		_, err = client.Auth().Bot(ctx, Bot.Token)
		if err != nil {
			return err
		}

		GotdClient = client
		GotdContext = ctx

		for {
			time.Sleep(time.Second * time.Duration(60))
		}
	})
}

func ErrorReporting(err error, context tele.Context) {
	_, fn, line, _ := runtime.Caller(1)
	log.Printf("[%s:%d] %v", fn, line, err)
	text := fmt.Sprintf("<pre>[%s:%d]\n%v</pre>", fn, line, strings.ReplaceAll(err.Error(), Config.Token, ""))
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
	Bot.Send(tele.ChatID(Config.SysAdmin), text)
}

func gatherData(update *tele.Update) error {
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
	if update.Message.Sender.IsBot || update.Message.Chat.ID != Config.Chat && update.Message.Chat.ID != Config.ReserveChat {
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

func checkPoint(update *tele.Update) error {
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
	for _, user := range RestrictedUsers {
		if update.Message.Sender.ID != user.UserID {
			continue
		}
		if update.Message.ReplyTo != nil && update.Message.ReplyTo.ID == WelcomeMessageID {
			delete := DB.Delete(&CheckPointRestrict{UserID: update.Message.Sender.ID})
			if delete.Error != nil {
				return delete.Error
			}
		} else {
			return Bot.Delete(update.Message)
		}
	}
	return nil
}

func init() {
	Bot.OnError = ErrorReporting

	Bot.Poller = tele.NewMiddlewarePoller(Bot.Poller, func(upd *tele.Update) bool {
		gatherData(upd)
		checkPoint(upd)

		return true
	})

	go gotdClientInit()
}
