package utils

import (
	cntx "context"

	gogpt "github.com/sashabaranov/go-gpt3"
	tele "gopkg.in/telebot.v3"
)

var c = gogpt.NewClient(Config.OpenAIKey)
var ctx = cntx.Background()

func ChatGPT(context tele.Context) error {
	if context.Message().ReplyTo == nil || context.Message().ReplyTo.Sender.ID != Bot.Me.ID || context.Message().Text[:1] == "/" {
		return nil
	}
	// if !IsAdminOrModer(context.Message().Sender.ID) {
	// 	return nil
	// }

	req := gogpt.ChatCompletionRequest{
		Model:    gogpt.GPT3Dot5Turbo,
		Messages: []gogpt.ChatCompletionMessage{{Role: "system", Content: "ты чатбот, который отвечает кратко"}, {Role: "assistant", Content: context.Message().ReplyTo.Text}, {Role: "user", Content: context.Message().Text}},
	}
	resp, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		return err
	}
	return context.Reply(resp.Choices[0].Message.Content)
}
