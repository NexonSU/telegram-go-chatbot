package commands

import (
	"bytes"
	"reflect"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/gotd/td/telegram/downloader"
	"github.com/gotd/td/tg"
	tele "gopkg.in/telebot.v3"
)

var channel *tg.Channel

//Reply with GIF from Pan Kotek's channel
func Meow(context tele.Context) error {
	api := utils.GotdClient.API()

	//get channel object
	if channel == nil {
		channelResolve, err := api.ContactsResolveUsername(utils.GotdContext, "imacat")
		if err != nil {
			return err
		}
		channel = channelResolve.GetChats()[0].(*tg.Channel)
	}
	//prepare message query
	messagesQuery := []tg.InputMessageClass{}
	firstMessageId := utils.RandInt(15, 16000)
	for message_id := firstMessageId; message_id < firstMessageId+10; message_id++ {
		messageObject := tg.Message{ID: message_id}
		messagesQuery = append(messagesQuery, messageObject.AsInputMessageID())
	}
	//query messages
	messagesResult, err := api.ChannelsGetMessages(utils.GotdContext, &tg.ChannelsGetMessagesRequest{
		Channel: channel.AsInput(),
		ID:      messagesQuery,
	})
	if err != nil {
		return err
	}
	//search and download gif
	buf := bytes.Buffer{}
	for _, mc := range messagesResult.(*tg.MessagesChannelMessages).Messages {
		if reflect.TypeOf(mc) != reflect.TypeOf(&tg.Message{}) {
			continue
		}
		messageMediaClass, check := mc.(*tg.Message).GetMedia()
		if check && reflect.TypeOf(messageMediaClass) == reflect.TypeOf(&tg.MessageMediaDocument{}) {
			document, _ := messageMediaClass.(*tg.MessageMediaDocument).GetDocument()
			docFile, _ := document.AsNotEmpty()
			fileName, _ := docFile.MapAttributes().AsDocumentAttributeFilename().First()
			downloader.NewDownloader().Download(api, docFile.AsInputDocumentFileLocation()).Stream(utils.GotdContext, &buf)
			return context.Reply(&tele.Document{FileName: fileName.FileName, File: tele.File{FileReader: &buf}})
		} else {
			continue
		}
	}

	return nil
}
