package commands

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"unicode/utf8"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
	tele "gopkg.in/telebot.v3"
)

// Send Yandex 300 response on link
func TLDR(context tele.Context) error {
	if utils.Config.YandexSummarizerToken == "" {
		return fmt.Errorf("не задан Yandex Summarizer токен")
	}
	if context.Message().ReplyTo == nil && len(context.Args()) == 0 {
		return utils.ReplyAndRemove("Бот заберёт статью по ссылке и сделает её краткое описание.\nПример использования:\n<code>/tldr ссылка</code>.\nИли отправь в ответ на какое-либо сообщение с ссылкой.", context)
	}

	link := ""
	message := &tele.Message{}

	if context.Message().ReplyTo == nil {
		message = context.Message()
	} else {
		message = context.Message().ReplyTo
	}

	for _, entity := range message.Entities {
		if entity.Type == tele.EntityURL || entity.Type == tele.EntityTextLink {
			link = entity.URL
			if link == "" {
				link = message.EntityText(entity)
			}
		}
	}

	if link == "" {
		for _, entity := range message.CaptionEntities {
			if entity.Type == tele.EntityURL || entity.Type == tele.EntityTextLink {
				link = entity.URL
				if link == "" {
					link = message.EntityText(entity)
				}
			}
		}
	}

	if link == "" {
		return utils.ReplyAndRemove("Бот заберёт статью по ссылке и сделает её краткое описание.\nПример использования:\n<code>/tldr ссылка</code>.\nИли отправь в ответ на какое-либо сообщение с ссылкой.", context)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://300.ya.ru/api/sharing-url",
		bytes.NewBuffer([]byte(`{"article_url": "`+link+`"}`)))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "OAuth "+utils.Config.YandexSummarizerToken)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if gjson.Get(string(body), "status").Str != "success" {
		return fmt.Errorf("ошибка, статус: %v", gjson.Get(string(body), "status").Str)
	}

	res, err := http.Get(gjson.Get(string(body), "sharing_url").Str)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return err
	}

	text := doc.Find(".summary .summary-content .summary-text").Text()

	text = regexp.MustCompile(`\n\s+\n|\n\n`).ReplaceAllString(text, "\n")
	text = regexp.MustCompile(`[ ]+`).ReplaceAllString(text, ` `)

	if utf8.RuneCountInString(text) > 4000 {
		text = string([]rune(text)[:4000])
	}

	//\n          \n

	return context.Reply(text)
}
