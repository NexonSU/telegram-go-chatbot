package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/dustin/go-humanize"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	"gopkg.in/tucnak/telebot.v3"
)

var CryptoMap []*cmc.MapListing
var FiatMap []*cmc.FiatMapListing

func GenerateMaps() {
	if utils.Config.CurrencyKey == "" {
		return
	}
	var err error
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	CryptoMap, err = client.Cryptocurrency.Map(&cmc.MapOptions{ListingStatus: "active,untracked"})
	if err != nil {
		log.Fatalln(err)
	}
	FiatMap, err = client.Fiat.Map(&cmc.FiatMapOptions{IncludeMetals: true})
	if err != nil {
		log.Fatalln(err)
	}
}

func GetSymbolId(symbol string) (string, error) {
	symbol = strings.ToUpper(symbol)
	if symbol == "BYR" {
		symbol = "BYN"
	}
	if symbol == "COC" {
		symbol = "RUB"
	}
	for _, fiat := range FiatMap {
		if fiat.Symbol == symbol {
			return fmt.Sprintf("%v", int(fiat.ID)), nil
		}
	}
	for _, crypto := range CryptoMap {
		if crypto.Symbol == symbol {
			return fmt.Sprintf("%v", int(crypto.ID)), nil
		}
	}
	return "", fmt.Errorf("не удалось распознать валюту <code>%v</code>", symbol)
}

func GetIdName(ID string) string {
	ID_int, _ := strconv.Atoi(ID)
	for _, fiat := range FiatMap {
		if int(fiat.ID) == ID_int {
			return fiat.Name
		}
	}
	for _, crypto := range CryptoMap {
		if int(crypto.ID) == ID_int {
			return crypto.Name
		}
	}
	return ""
}

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
	symbol, err := GetSymbolId(context.Args()[1])
	if err != nil {
		return context.Reply(err.Error())
	}
	convert, err := GetSymbolId(context.Args()[2])
	if err != nil {
		return context.Reply(err.Error())
	}
	//COC
	if strings.ToUpper(context.Args()[1]) == "COC" {
		amount = amount * 300
	}
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount: amount, ID: symbol, ConvertID: convert})
	if err != nil {
		return context.Reply(fmt.Sprintf("Ошибка при запросе: %v\nОнлайн-версия: https://coinmarketcap.com/ru/converter/", err.Error()), &telebot.SendOptions{DisableWebPagePreview: true})
	}
	resultAmount := conversion.Quote[convert].Price
	resultName := GetIdName(convert)
	//COC
	if strings.ToUpper(context.Args()[1]) == "COC" {
		conversion.Amount = amount / 300
		conversion.Name = "Cup Of Coffee"
	}
	if strings.ToUpper(context.Args()[2]) == "COC" {
		resultAmount = resultAmount / 300
		resultName = "Cup Of Coffee"
	}
	return context.Send(fmt.Sprintf("%v %v = %v %v", conversion.Amount, conversion.Name, strings.Replace(humanize.CommafWithDigits(resultAmount, 2), ",", " ", -1), resultName), &telebot.SendOptions{ReplyTo: context.Message().ReplyTo, AllowWithoutReply: true})
}
