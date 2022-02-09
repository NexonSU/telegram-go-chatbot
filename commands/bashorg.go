package commands

import (
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"gopkg.in/telebot.v3"
)

var lastBash string

//Send text in chat on /say
func Bashorg(context telebot.Context) error {
	httpClient := &http.Client{Timeout: 10 * time.Second}
	httpResponse, err := httpClient.Get("https://bash.im/forweb/?u")
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		Body.Close()
	}(httpResponse.Body)
	html, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		return err
	}

	text := string(html)
	text = strings.ReplaceAll(text, "' + '", "")
	text = strings.Split(text, "/header>")[1]
	text = strings.Split(text, "<footer")[0]
	text = strings.ReplaceAll(text, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "<br />", "\n")
	tags := regexp.MustCompile(`<div.*">|</div>`)
	text = tags.ReplaceAllString(text, ``)

	if lastBash == text {
		text += "\n\n<code>https://bash.im/forweb/?u вернул кэшированный баш.\nПодождите ~10 секунд, чтобы получить новый.</code>"
	} else {
		lastBash = text
	}

	return context.Reply(text)
}
