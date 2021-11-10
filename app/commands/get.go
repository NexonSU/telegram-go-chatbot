package commands

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send Get to user on /get
func Get(context telebot.Context) error {
	var get utils.Get
	if len(context.Args()) != 1 {
		return context.Reply("Пример использования: <code>/get {гет}</code>")
	}
	result := utils.DB.Where(&utils.Get{Name: strings.ToLower(context.Data())}).First(&get)
	if get.Caption == "" {
		get.Caption = get.Name
	}
	if result.RowsAffected != 0 {
		switch {
		case get.Type == "Animation":
			return context.Reply(&telebot.Animation{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Audio":
			return context.Reply(&telebot.Audio{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Photo":
			return context.Reply(&telebot.Photo{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Video":
			return context.Reply(&telebot.Video{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Voice":
			return context.Reply(&telebot.Voice{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Document":
			return context.Reply(&telebot.Document{
				File:    telebot.File{FileID: get.Data},
				Caption: get.Caption,
			})
		case get.Type == "Text":
			return context.Reply(get.Data)
		default:
			return context.Reply(fmt.Sprintf("Ошибка при определении типа гета, я не знаю тип <code>%v</code>.", get.Type))
		}
	} else {
		return context.Reply(fmt.Sprintf("Гет <code>%v</code> не найден.", context.Data()))
	}
}

//Answer on inline get query
func GetInline(context telebot.Context) error {
	var count int64
	query := strings.ToLower(context.Query().Text)
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
				Title:       get.Name,
				Caption:     get.Caption,
				Cache:       get.Data,
				Description: get.Caption,
			}
		case get.Type == "Photo":
			results[i] = &telebot.PhotoResult{
				Title:       get.Name,
				Caption:     get.Caption,
				Cache:       get.Data,
				Description: get.Caption,
			}
		case get.Type == "Video":
			results[i] = &telebot.VideoResult{
				Title:       get.Name,
				Caption:     get.Caption,
				Cache:       get.Data,
				Description: get.Caption,
			}
		case get.Type == "Voice":
			results[i] = &telebot.VoiceResult{
				Title:   get.Name,
				Caption: get.Caption,
				Cache:   get.Data,
			}
		case get.Type == "Document":
			results[i] = &telebot.DocumentResult{
				Title:       get.Name,
				Caption:     get.Caption,
				Cache:       get.Data,
				Description: get.Caption,
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
		Results: results,
	})
}
