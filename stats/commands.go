package stats

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/go-echarts/go-echarts/v2/components"
	tele "gopkg.in/telebot.v3"
)

func RemoveWord(context tele.Context) error {
	if len(context.Args()) != 1 {
		return context.Reply("Укажите слово.")
	}
	word := context.Data()
	//remove word
	delete := utils.DB.Where("text = ? OR text = ?", word, strings.ToLower(word)).Delete(&utils.Word{})
	if delete.Error != nil {
		return context.Reply(fmt.Sprintf("Не удалось удалить слово:\n<code>%v</code>", delete.Error.Error()))
	}
	//add word to DB
	exclude := utils.DB.Create(&utils.WordStatsExclude{Text: strings.ToLower(word)})
	if exclude.Error != nil {
		return context.Reply(fmt.Sprintf("Не удалось запретить слово:\n<code>%v</code>", exclude.Error.Error()))
	}
	//update utils.WordStatsExcludes
	utils.DB.Find(&utils.WordStatsExcludes)
	return context.Reply("Слово запрещено для статистики и удалено из базы.")
}

func Stats(context tele.Context) error {
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
			return context.Reply(&tele.Video{File: tele.File{FileID: "BAACAgIAAx0CRXO-MQADWWB4LQABzrOqWPkq-JXIi4TIixY4dwACPw4AArBgwUt5sRu-_fDR5x4E"}})
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
		page.PageTitle = fmt.Sprintf("%v stats since %v to %v", context.Chat().Title, from.Format("02.01.2006"), to.Format("02.01.2006"))
		page.Theme = "shine"
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
	return context.Reply(&tele.Document{
		File: tele.File{
			FileReader: f,
		},
		FileName: fmt.Sprintf("%v %v %v - %v.html", selected, context.Chat().Username, from.Format("02.01.2006"), time.Now().Format("02.01.2006")),
	})
}
