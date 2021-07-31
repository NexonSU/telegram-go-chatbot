package commands

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	"gopkg.in/tucnak/telebot.v3"
)

//Reply currency "cur"
func Cur(context telebot.Context) error {
	if utils.Config.CurrencyKey == "" {
		return context.Reply("Конвертация валют не настроена")
	}
	var text = strings.Split(context.Text(), " ")
	if len(text) != 4 {
		return context.Reply("Пример использования:\n/cur {количество} {EUR/USD/RUB} {EUR/USD/RUB}")
	}
	if context.Message().ReplyTo != nil {
		context.Message().Sender = context.Message().ReplyTo.Sender
	}
	amount, err := strconv.ParseFloat(text[1], 64)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения количества:\n<code>%v</code>", err))
	}
	var symbol = strings.ToUpper(text[2])
	if !regexp.MustCompile(`^[A-Z$]{3,5}$`).MatchString(symbol) {
		return context.Reply("Имя валюты должно состоять из 3-5 латинских символов.")
	}
	var convert = strings.ToUpper(text[3])
	if !regexp.MustCompile(`^[A-Z$]{3,5}$`).MatchString(convert) {
		return context.Reply("Имя валюты должно состоять из 3-5 латинских символов.")
	}
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount: amount, Symbol: symbol, Convert: convert})
	if err != nil {
		return context.Reply("Ошибка при запросе. Возможно, одна из валют не найдена.\nОнлайн-версия: https://coinmarketcap.com/ru/converter/", &telebot.SendOptions{DisableWebPagePreview: true})
	}
	return context.Reply(fmt.Sprintf("%v %v = %v %v", conversion.Amount, conversion.Name, math.Round(conversion.Quote[convert].Price*100)/100, convert))
}
