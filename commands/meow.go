package commands

import (
	tdctx "context"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/gotd/td/bin"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
	tele "gopkg.in/telebot.v3"
)

//Reply with GIF from Pan Kotek's channel
func Meow(context tele.Context) error {
	client := telegram.NewClient(utils.Config.AppID, utils.Config.AppHash, telegram.Options{})
	err := client.Run(tdctx.Background(), func(ctx tdctx.Context) error {
		_, err := client.Auth().Bot(ctx, utils.Bot.Token)
		if err != nil {
			return err
		}
		api := client.API()
		sender := message.NewSender(api)

		channelResolve, err := api.ContactsResolveUsername(ctx, "imacat")
		if err != nil {
			return err
		}
		channel, _ := channelResolve.MapChats().AsChannel().First()
		if err != nil {
			return err
		}
		messageObject := tg.Message{ID: utils.RandInt(20, 16000)}
		buf := bin.Buffer{}
		messageObject.AsInputMessageID().Encode(&buf)
		message, err := tg.DecodeInputMessage(&buf)
		if err != nil {
			return err
		}
		messagesResult, err := api.ChannelsGetMessages(ctx, &tg.ChannelsGetMessagesRequest{
			Channel: channel.AsInput(),
			ID:      []tg.InputMessageClass{message},
		})
		if err != nil {
			return err
		}

		messagesResult.Encode(&buf)
		messageResult := tg.MessagesChannelMessages{}
		messageResult.Decode(&buf)
		messageResult.Messages[0].Encode(&buf)
		messageSend := tg.Message{}
		messageSend.Decode(&buf)
		media, _ := messageSend.GetMedia()
		if media == nil {
			return nil
		}
		media.Encode(&buf)
		messageMediaDocument := &tg.MessageMediaDocument{}
		messageMediaDocument.Decode(&buf)
		documentClass, _ := messageMediaDocument.GetDocument()
		document, _ := documentClass.AsNotEmpty()

		_, err = sender.Resolve(context.Chat().Username).Document(ctx, document.AsInput())
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
