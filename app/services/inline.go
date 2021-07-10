package services

import (
	"fmt"
	"log"
	"strconv"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
)

//Answer on inline query
func OnInline(q *tb.Query) {
	var count int64
	gets := utils.DB.Limit(50).Model(utils.Get{}).Where("name LIKE ?", "%"+q.Text+"%").Count(&count)
	get_rows, err := gets.Rows()
	if err != nil {
		log.Println(err.Error())
		return
	}
	if count > 50 {
		count = 50
	}
	results := make(tb.Results, count)
	var i int
	for get_rows.Next() {
		var get utils.Get
		err := utils.DB.ScanRows(get_rows, &get)
		if err != nil {
			log.Println(err.Error())
			return
		}
		switch {
		case get.Type == "Animation":
			results[i] = &tb.GifResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Audio":
			results[i] = &tb.AudioResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Photo":
			results[i] = &tb.PhotoResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Video":
			results[i] = &tb.VideoResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Voice":
			results[i] = &tb.VoiceResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Document":
			results[i] = &tb.DocumentResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Text":
			results[i] = &tb.ArticleResult{
				Title:       get.Name,
				Description: get.Data,
			}
			results[i].SetContent(tb.InputMessageContent(&tb.InputTextMessageContent{
				Text:      fmt.Sprintf("<b>%v</b>\n%v", get.Name, get.Data),
				ParseMode: "HTML",
			}))
		default:
			log.Printf("Не удалось отправить гет %v через inline.", get.Name)
		}

		results[i].SetResultID(strconv.Itoa(i))

		i++
	}

	err = utils.Bot.Answer(q, &tb.QueryResponse{
		Results:   results,
		CacheTime: 0,
	})

	if err != nil {
		log.Println(err.Error())
	}
}
