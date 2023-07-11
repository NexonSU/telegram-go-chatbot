package commands

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/tidwall/gjson"
	scraper "github.com/yudgnahk/go-cloudflare-scraper"
	tele "gopkg.in/telebot.v3"
)

// Send hentai manga by date or userid
func GetMyHentai(context tele.Context) error {
	if utils.Config.NHentaiCookie == "" {
		return utils.SendAndRemove("Отсутствуют куки в конфиге", context)
	}
	hentaiId := ""
	if len(context.Args()) == 1 {
		date, err := time.Parse("02.01.2006", context.Args()[0])
		if err != nil {
			return utils.SendAndRemove("Ошибка парсинга даты: "+err.Error(), context)
		}
		hentaiId = date.Format("20106")
	} else {
		hentaiId = fmt.Sprint(context.Sender().ID % 1e5)
	}

	resp, err := getHentai(hentaiId)
	if err != nil {
		return err
	}

	if resp.StatusCode == 404 {
		reverseId := ""
		for _, v := range hentaiId {
			reverseId = string(v) + reverseId
		}
		hentaiId = strings.TrimLeft(reverseId, "0")
		resp, err = getHentai(hentaiId)
		if err != nil {
			return err
		}
	}

	if resp.StatusCode == 404 {
		return utils.SendAndRemove("Сорян, для тебя нет хентай-манги. Возможно, её удалили.", context)
	}

	if resp.StatusCode != 200 {
		return utils.SendAndRemove("Ошибка запроса: "+resp.Status, context)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	title := gjson.Get(string(data), "title.pretty").String()
	tagsResult := gjson.Get(string(data), "tags.#.name").Array()
	tagsArray := []string{}

	for _, name := range tagsResult {
		tagsArray = append(tagsArray, name.String())
	}

	if err != nil {
		return err
	}

	return context.Reply(fmt.Sprintf("%v, твоя хентай-манга: %v\nТвои теги: %v\nСсылка для друга: <span class=\"tg-spoiler\">https://nhentai.net/g/%v/</span>", utils.MentionUser(context.Sender()), title, strings.Join(tagsArray, ", "), hentaiId))
}

func getHentai(hentaiId string) (*http.Response, error) {
	scraper, err := scraper.NewTransport(http.DefaultTransport)
	if err != nil {
		return nil, err
	}

	client := http.Client{Transport: scraper}
	req, err := http.NewRequest("GET", "https://nhentai.net/api/gallery/"+hentaiId, strings.NewReader(""))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/114.0.0.0 Safari/537.36")
	req.Header.Set("Cookie", utils.Config.NHentaiCookie)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
