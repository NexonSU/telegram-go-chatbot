package stats

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/components"
)

func RemoveWord(context telebot.Context) error {
	if len(context.Args()) != 1 {
		return context.Reply("Укажите слово.")
	}
	result := utils.DB.Where("text = ?", context.Data()).Delete(&utils.Word{})
	if result.RowsAffected != 0 {
		return context.Reply("Слово удалено.")
	} else {
		return context.Reply(fmt.Sprintf("Не удалось удалить слово:\n<code>%v</code>", result.Error.Error()))
	}
}

func Stats(context telebot.Context) error {
	selected := "Stats"
	graphs := []string{"Activity", "MostActiveToday", "PopDays", "PopHours", "PopWords", "TopUsers"}
	days := 30
	if len(context.Args()) >= 1 {
		var err error
		days, err = strconv.Atoi(context.Args()[0])
		if err != nil {
			return context.Reply("Ошибка определения дней.")
		}
		if days == 2077 {
			return context.Reply(&telebot.Video{File: telebot.File{FileID: "BAACAgIAAx0CRXO-MQADWWB4LQABzrOqWPkq-JXIi4TIixY4dwACPw4AArBgwUt5sRu-_fDR5x4E"}})
		}
	}
	if len(context.Args()) == 2 {
		for _, graph := range graphs {
			if strings.EqualFold(graph, context.Args()[1]) {
				selected = graph
			}
		}
		if selected == "Stats" {
			return context.Reply("Доступные графики:\n<pre>" + strings.Join(graphs, ", ") + "</pre>")
		}
	}
	days = days * -1
	from := time.Now().AddDate(0, 0, days)
	to := time.Now().Add(time.Hour * 24)
	f := new(bytes.Buffer)

	switch selected {
	case "Activity":
		UserActivityLineChart(from, to, context).Render(f)
	case "MostActiveToday":
		MostActiveUsersTodayPieChart(from, to, context).Render(f)
	case "PopDays":
		PopDaysBarChart(from, to, context).Render(f)
	case "PopHours":
		PopHoursBarChart(from, to, context).Render(f)
	case "PopWords":
		PopWordsWcChart(from, to, context).Render(f)
	case "TopUsers":
		TopUsersBarChart(from, to, context).Render(f)
	case "Stats":
		page := components.NewPage()
		page.SetLayout(components.PageFlexLayout)
		page.AddCharts(
			UserActivityLineChart(from, to, context),
			MostActiveUsersTodayPieChart(from, to, context),
			PopDaysBarChart(from, to, context),
			PopHoursBarChart(from, to, context),
			PopWordsWcChart(from, to, context),
			TopUsersBarChart(from, to, context),
		)
		page.Render(f)
	default:
		return nil
	}
	return context.Reply(&telebot.Document{
		File: telebot.File{
			FileReader: f,
		},
		FileName: fmt.Sprintf("%v %v %v - %v.html", selected, context.Chat().Username, from.Format("02.01.2006"), time.Now().Format("02.01.2006")),
	})
}
