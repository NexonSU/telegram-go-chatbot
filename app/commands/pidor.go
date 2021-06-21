package commands

import (
	"fmt"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

var busy = make(map[string]bool)

// Pidor game
func Pidor(m *tb.Message) {
	if m.Private() {
		return
	}
	if busy["pidor"] {
		_, err := utils.Bot.Reply(m, "Команда занята. Попробуйте позже.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return
	}
	busy["pidor"] = true
	defer func() { busy["pidor"] = false }()
	var pidor utils.PidorStats
	var pidorToday utils.PidorList
	result := utils.DB.Model(utils.PidorStats{}).Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local), time.Now()).First(&pidor)
	if result.RowsAffected == 0 {
		utils.DB.Model(utils.PidorList{}).Order("RANDOM()").First(&pidorToday)
		TargetChatMember, err := utils.Bot.ChatMemberOf(m.Chat, &tb.User{ID: pidorToday.ID})
		if err != nil {
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Я нашел пидора дня, но похоже, что с <a href=\"tg://user?id=%v\">%v</a> что-то не так, так что попробуйте еще раз, пока я удаляю его из игры! Ошибка:\n<code>%v</code>", pidorToday.ID, pidorToday.Username, err.Error()))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			utils.DB.Delete(pidorToday)
			return
		}
		if TargetChatMember.Role == "left" {
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Я нашел пидора дня, но похоже, что <a href=\"tg://user?id=%v\">%v</a> вышел из этого чата (вот пидор!), так что попробуйте еще раз, пока я удаляю его из игры!", pidorToday.ID, pidorToday.Username))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			utils.DB.Delete(pidorToday)
			return
		}
		if TargetChatMember.Role == "kicked" {
			_, err := utils.Bot.Reply(m, fmt.Sprintf("Я нашел пидора дня, но похоже, что <a href=\"tg://user?id=%v\">%v</a> был забанен в этом чате (получил пидор!), так что попробуйте еще раз, пока я удаляю его из игры!", pidorToday.ID, pidorToday.Username))
			if err != nil {
				utils.ErrorReporting(err, m)
				return
			}
			utils.DB.Delete(pidorToday)
			return
		}
		pidor.UserID = pidorToday.ID
		pidor.Date = time.Now()
		utils.DB.Create(pidor)
		messages := [][]string{
			{"Инициирую поиск пидора дня...", "Опять в эти ваши игрульки играете? Ну ладно...", "Woop-woop! That's the sound of da pidor-police!", "Система взломана. Нанесён урон. Запущено планирование контрмер.", "Сейчас поколдуем...", "Инициирую поиск пидора дня...", "Зачем вы меня разбудили...", "Кто сегодня счастливчик?"},
			{"Хм...", "Сканирую...", "Ведётся поиск в базе данных", "Сонно смотрит на бумаги", "(Ворчит) А могли бы на работе делом заниматься", "Военный спутник запущен, коды доступа внутри...", "Ну давай, посмотрим кто тут классный..."},
			{"Высокий приоритет мобильному юниту.", "Ох...", "Ого-го...", "Так, что тут у нас?", "В этом совершенно нет смысла...", "Что с нами стало...", "Тысяча чертей!", "Ведётся захват подозреваемого..."},
			{"Стоять! Не двигаться! Ты объявлен пидором дня, ", "Ого, вы посмотрите только! А пидор дня то - ", "Пидор дня обыкновенный, 1шт. - ", ".∧＿∧ \n( ･ω･｡)つ━☆・*。 \n⊂  ノ    ・゜+. \nしーＪ   °。+ *´¨) \n         .· ´¸.·*´¨) \n          (¸.·´ (¸.·'* ☆ ВЖУХ И ТЫ ПИДОР, ", "Ага! Поздравляю! Сегодня ты пидор - ", "Кажется, пидор дня - ", "Анализ завершен. Ты пидор, "},
		}
		for i := 0; i <= 3; i++ {
			duration := time.Second * time.Duration(i*2)
			message := messages[i][utils.RandInt(0, len(messages[i])-1)]
			if i == 3 {
				message += fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a>", pidorToday.ID, pidorToday.Username)
			}
			go func() {
				time.Sleep(duration)
				_, err := utils.Bot.Send(m.Chat, message)
				if err != nil {
					utils.ErrorReporting(err, m)
					return
				}
			}()
		}
	} else {
		utils.DB.Model(utils.PidorList{}).Where(pidor.UserID).First(&pidorToday)
		_, err := utils.Bot.Reply(m, fmt.Sprintf("Согласно моей информации, по результатам сегодняшнего розыгрыша пидор дня - %v!", pidorToday.Username))
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
	}
}
