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
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	if utils.Config.CurrencyKey == "" {
		err := context.Reply("Конвертация валют не настроена")
		if err != nil {
			return err
		}
		return err
	}
	var target = context.Message()
	var text = strings.Split(context.Text(), " ")
	if len(text) != 4 {
		err := context.Reply("Пример использования:\n/cur {количество} {EUR/USD/RUB} {EUR/USD/RUB}")
		if err != nil {
			return err
		}
		return err
	}
	if context.Message().ReplyTo != nil {
		target = context.Message().ReplyTo
	}
	amount, err := strconv.ParseFloat(text[1], 64)
	if err != nil {
		err := context.Reply(fmt.Sprintf("Ошибка определения количества:\n<code>%v</code>", err))
		if err != nil {
			return err
		}
		return err
	}
	var symbol = strings.ToUpper(text[2])
	if !regexp.MustCompile(`^[A-Z$]{3,5}$`).MatchString(symbol) {
		err := context.Reply("Имя валюты должно состоять из 3-5 латинских символов.")
		if err != nil {
			return err
		}
		return err
	}
	var convert = strings.ToUpper(text[3])
	if !regexp.MustCompile(`^[A-Z$]{3,5}$`).MatchString(convert) {
		err := context.Reply("Имя валюты должно состоять из 3-5 латинских символов.")
		if err != nil {
			return err
		}
		return err
	}
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount: amount, Symbol: symbol, Convert: convert})
	if err != nil {
		err := context.Reply("Ошибка при запросе. Возможно, одна из валют не найдена.\nОнлайн-версия: https://coinmarketcap.com/ru/converter/", &telebot.SendOptions{DisableWebPagePreview: true})
		if err != nil {
			return err
		}
		return err
	}
	_, err = utils.Bot.Reply(target, fmt.Sprintf("%v %v = %v %v", conversion.Amount, conversion.Name, math.Round(conversion.Quote[convert].Price*100)/100, convert))
	if err != nil {
		return err
	}
	return err
}
