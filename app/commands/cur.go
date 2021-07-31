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
	if len(context.Args()) != 3 {
		return context.Reply("Пример использования:\n/cur {количество} {EUR/USD/RUB} {EUR/USD/RUB}")
	}
	amount, err := strconv.ParseFloat(context.Args()[0], 64)
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка определения количества:\n<code>%v</code>", err))
	}
	var symbol = strings.ToUpper(context.Args()[1])
	if !regexp.MustCompile(`^[A-Z$]{3,5}$`).MatchString(symbol) {
		return context.Reply("Имя валюты должно состоять из 3-5 латинских символов.")
	}
	var convert = strings.ToUpper(context.Args()[2])
	if !regexp.MustCompile(`^[A-Z$]{3,5}$`).MatchString(convert) {
		return context.Reply("Имя валюты должно состоять из 3-5 латинских символов.")
	}
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount: amount, Symbol: symbol, Convert: convert})
	if err != nil {
		return context.Reply("Ошибка при запросе. Возможно, одна из валют не найдена.\nОнлайн-версия: https://coinmarketcap.com/ru/converter/", &telebot.SendOptions{DisableWebPagePreview: true})
	}
	return context.Send(fmt.Sprintf("%v %v = %v %v", conversion.Amount, conversion.Name, math.Round(conversion.Quote[convert].Price*100)/100, convert), &telebot.SendOptions{ReplyTo: context.Message().ReplyTo, AllowWithoutReply: true})
}
