package commands

import (
	"bytes"
	"fmt"
	"html"
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
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://300.ya.ru/api/sharing-url",
		bytes.NewBuffer([]byte(`{"article_url": "`+html.EscapeString(context.Data())+`"}`)))
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

	text = regexp.MustCompile(`\n\n\s+|\n\s+\n\s+|\n\n`).ReplaceAllString(text, "\n")

	if utf8.RuneCountInString(text) > 4000 {
		text = string([]rune(text)[:4000])
	}

	//\n          \n

	return context.Reply(text)
}
