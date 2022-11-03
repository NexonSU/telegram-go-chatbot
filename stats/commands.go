package stats

import (
	"bytes"
	cntx "context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"github.com/go-echarts/go-echarts/v2/components"
	tele "gopkg.in/telebot.v3"
)

var ctx cntx.Context

func init() {
	ctx, _ = chromedp.NewContext(
		cntx.Background(),
	)
}

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
	if len(context.Args()) > 0 {
		for _, graph := range graphs {
			if strings.EqualFold(graph, context.Args()[0]) {
				selected = graph
			}
		}
		if selected == "Stats" {
			return context.Reply("Доступные графики:\n<pre>" + strings.Join(graphs, ", ") + "</pre>")
		}
	}
	from := time.Now().AddDate(0, -1, 0)
	to := time.Now().Add(time.Hour * 24)
	filename := fmt.Sprintf("%v_%v_%v-%v.html", selected, context.Chat().ID, from.Format("02.01.2006"), time.Now().Format("02.01.2006"))
	filepath := os.TempDir() + "/" + filename
	f, _ := os.Create(filepath)

	width := int64(900)

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
		width = 1900
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
	var buf []byte
	if err := chromedp.Run(ctx, chromedp.Tasks{
		chromedp.Emulate(device.Reset),
		chromedp.Navigate("file:///" + filepath),
		chromedp.EmulateViewport(width, 0),
		chromedp.Sleep(time.Second),
		chromedp.FullScreenshot(&buf, 100),
	}); err != nil {
		return err
	}

	os.Remove(filepath)

	return context.Reply(&tele.Photo{File: tele.File{FileReader: bytes.NewBuffer(buf)}})
}
