package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

//Answer on inline query
func OnInline(context telebot.Context) error {
	var count int64
	query := context.Query().Text
	if query == "" {
		return context.Answer(&telebot.QueryResponse{})
	}
	gets := utils.DB.Limit(10).Model(utils.Get{}).Where("name LIKE ?", query+"%").Count(&count)
	get_rows, err := gets.Rows()
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if count > 10 {
		count = 10
	}
	results := make(telebot.Results, count)
	var i int
	for get_rows.Next() {
		var get utils.Get
		err := utils.DB.ScanRows(get_rows, &get)
		if err != nil {
			log.Println(err.Error())
			return err
		}
		if get.Caption == "" {
			get.Caption = get.Name
		}
		switch {
		case get.Type == "Animation":
			results[i] = &telebot.GifResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Audio":
			results[i] = &telebot.DocumentResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Photo":
			results[i] = &telebot.PhotoResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Video":
			results[i] = &telebot.VideoResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Voice":
			results[i] = &telebot.VoiceResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Document":
			results[i] = &telebot.DocumentResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Text":
			results[i] = &telebot.ArticleResult{
				Title:       get.Name,
				Description: get.Data,
			}
			results[i].SetContent(telebot.InputMessageContent(&telebot.InputTextMessageContent{
				Text:      fmt.Sprintf("<b>%v</b>\n%v", get.Name, get.Data),
				ParseMode: telebot.ModeHTML,
			}))
		default:
			log.Printf("Не удалось отправить гет %v через inline.", get.Name)
		}

		results[i].SetResultID(strconv.Itoa(i))

		i++
		if i >= int(count) {
			continue
		}
	}

	return context.Answer(&telebot.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})
}
