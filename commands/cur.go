package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	cmc "github.com/miguelmota/go-coinmarketcap/pro/v1"
	tele "gopkg.in/telebot.v3"
)

var CryptoMap []*cmc.MapListing
var FiatMap []*cmc.FiatMapListing
var JokeMap = []JokeMapStruct{}
var _ = GenerateMaps()

type JokeMapStruct struct {
	symbol string
	name   string
	amount float64
}

func GenerateMaps() error {
	if utils.Config.CurrencyKey == "" {
		return nil
	}
	var err error
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	CryptoMap, err = client.Cryptocurrency.Map(&cmc.MapOptions{ListingStatus: "active,untracked"})
	if err != nil {
		return err
	}
	FiatMap, err = client.Fiat.Map(&cmc.FiatMapOptions{IncludeMetals: true})
	if err != nil {
		return err
	}
	JokeMap = append(JokeMap, JokeMapStruct{symbol: "COC", name: "Cup Of Coffee", amount: 300.0})
	JokeMap = append(JokeMap, JokeMapStruct{symbol: "DSHK", name: "Doshirak", amount: 71.0})
	JokeMap = append(JokeMap, JokeMapStruct{symbol: "DOSH", name: "Doshirak", amount: 71.0})
	JokeMap = append(JokeMap, JokeMapStruct{symbol: "TBW", name: "Boeing Wing", amount: 178000000.0})
	return nil
}

func GetSymbolId(symbol string) (string, error) {
	symbol = strings.ToUpper(symbol)
	if symbol == "BYR" {
		symbol = "BYN"
	}
	for _, JokeFiat := range JokeMap {
		if symbol == JokeFiat.symbol {
			symbol = "RUB"
		}
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

// Reply currency "cur"
func Cur(context tele.Context) error {
	if utils.Config.CurrencyKey == "" {
		return utils.ReplyAndRemove("Конвертация валют не настроена", context)
	}
	if len(context.Args()) != 3 {
		return utils.ReplyAndRemove("Пример использования:\n/cur 1 USD RUB", context)
	}
	amount, err := strconv.ParseFloat(context.Args()[0], 64)
	if err != nil {
		return err
	}
	symbol, err := GetSymbolId(context.Args()[1])
	if err != nil {
		return err
	}
	convert, err := GetSymbolId(context.Args()[2])
	if err != nil {
		return err
	}
	for _, JokeFiat := range JokeMap {
		if strings.ToUpper(context.Args()[1]) == JokeFiat.symbol {
			amount = amount * JokeFiat.amount
		}
	}
	if symbol == "9911" || convert == "9911" {
		return fmt.Errorf("Невозможно конвертировать тестовую валюту")
	}
	client := cmc.NewClient(&cmc.Config{ProAPIKey: utils.Config.CurrencyKey})
	conversion, err := client.Tools.PriceConversion(&cmc.ConvertOptions{Amount: amount, ID: symbol, ConvertID: convert})
	if err != nil {
		return err
	}
	resultAmount := conversion.Quote[convert].Price
	resultName := GetIdName(convert)
	for _, JokeFiat := range JokeMap {
		if strings.ToUpper(context.Args()[1]) == JokeFiat.symbol {
			conversion.Amount = amount / JokeFiat.amount
			conversion.Name = JokeFiat.name
		}
		if strings.ToUpper(context.Args()[2]) == JokeFiat.symbol {
			resultAmount = resultAmount / JokeFiat.amount
			resultName = JokeFiat.name
		}
	}
	return context.Send(fmt.Sprintf("%.2f %v = %.2f %v", conversion.Amount, conversion.Name, resultAmount, resultName), &tele.SendOptions{ReplyTo: context.Message().ReplyTo, AllowWithoutReply: true})
}
