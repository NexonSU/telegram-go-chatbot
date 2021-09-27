package pidor

import (
	"fmt"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
)

var busy = make(map[string]bool)

// Pidor game
func Pidor(context telebot.Context) error {
	if context.Message().Private() {
		return nil
	}
	if busy["pidor"] {
		return context.Reply("Команда занята. Попробуйте позже.")
	}
	busy["pidor"] = true
	defer func() { busy["pidor"] = false }()
	var pidor utils.PidorStats
	var pidorToday utils.PidorList
	result := utils.DB.Model(utils.PidorStats{}).Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local), time.Now()).First(&pidor)
	if result.RowsAffected == 0 {
		utils.DB.Model(utils.PidorList{}).Order("RANDOM()").First(&pidorToday)
		TargetChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), &telebot.User{ID: pidorToday.ID})
		if err != nil {
			utils.DB.Delete(pidorToday)
			return context.Reply(fmt.Sprintf("Я нашел пидора дня, но похоже, что с <a href=\"tg://user?id=%v\">%v</a> что-то не так, так что попробуйте еще раз, пока я удаляю его из игры! Ошибка:\n<code>%v</code>", pidorToday.ID, pidorToday.Username, err.Error()))
		}
		if TargetChatMember.Role == "left" {
			utils.DB.Delete(pidorToday)
			return context.Reply(fmt.Sprintf("Я нашел пидора дня, но похоже, что <a href=\"tg://user?id=%v\">%v</a> вышел из этого чата (вот пидор!), так что попробуйте еще раз, пока я удаляю его из игры!", pidorToday.ID, pidorToday.Username))
		}
		if TargetChatMember.Role == "kicked" {
			utils.DB.Delete(pidorToday)
			return context.Reply(fmt.Sprintf("Я нашел пидора дня, но похоже, что <a href=\"tg://user?id=%v\">%v</a> был забанен в этом чате (получил пидор!), так что попробуйте еще раз, пока я удаляю его из игры!", pidorToday.ID, pidorToday.Username))
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
				context.Send(message)
			}()
		}
	} else {
		utils.DB.Model(utils.PidorList{}).Where(pidor.UserID).First(&pidorToday)
		return context.Reply(fmt.Sprintf("Согласно моей информации, по результатам сегодняшнего розыгрыша пидор дня - %v!", pidorToday.Username))
	}
	return nil
}
