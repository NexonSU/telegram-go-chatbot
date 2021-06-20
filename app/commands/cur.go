package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	tb "gopkg.in/tucnak/telebot.v2"
	"math"
	"regexp"
	"strconv"
	"strings"
)

//Reply currency "cur"
func Cur(m *tb.Message) {
	var target = *m
	var text = strings.Split(m.Text, " ")
	if len(text) != 4 {
		_, err := utils.Bot.Reply(m, "Пример использования:\n/cur {количество} {EUR/USD/RUB} {EUR/USD/RUB}")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	if m.ReplyTo != nil {
		target = *m.ReplyTo
	}
	amount, err := strconv.ParseFloat(text[1], 64)
	if err != nil {
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Ошибка определения количества:\n<code>%v</code>", err))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var symbol = strings.ToUpper(text[2])
	if !regexp.MustCompile(`^[A-Z]{3,4}$`).MatchString(symbol) {
		_, err := utils.Bot.Reply(m, "Имя валюты должно состоять из 3-4 больших латинских символов.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	var convert = strings.ToUpper(text[3])
	if !regexp.MustCompile(`^[A-Z]{3,4}$`).MatchString(convert) {
		_, err := utils.Bot.Reply(m, "Имя валюты должно состоять из 3-4 больших латинских символов.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount: amount, Symbol: symbol, Convert: convert})
	if err != nil {
		_, err := utils.Bot.Reply(m, "Ошибка при запросе. Возможно, одна из валют не найдена.\nОнлайн-версия: https://coinmarketcap.com/ru/converter/", &tb.SendOptions{DisableWebPagePreview: true})
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	_, err = utils.Bot.Reply(&target, fmt.Sprintf("%v %v = %v %v", conversion.Amount, conversion.Name, math.Round(conversion.Quote[convert].Price*100)/100, convert))
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
