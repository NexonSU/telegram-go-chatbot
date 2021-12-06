package stats

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"gopkg.in/tucnak/telebot.v3"
)

func PopWordsWcChart(from time.Time, to time.Time, context telebot.Context) *charts.WordCloud {
	result, _ := utils.DB.
		Model(utils.Word{}).
		Select("text, COUNT(*) as count").
		Where("chat_id = ?", context.Chat().ID).
		Where("date BETWEEN ? AND ?", from, to).
		Group("text").
		Order("count DESC").
		Limit(100).
		Rows()
	var Word string
	var Count int
	var WCData []opts.WordCloudData
	for result.Next() {
		err := result.Scan(&Word, &Count)
		if err != nil {
			utils.ErrorReporting(err, context)
			return nil
		}
		WCData = append(WCData, opts.WordCloudData{Name: Word, Value: Count})
	}

	wc := charts.NewWordCloud()
	wc.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: fmt.Sprintf("%v Popular Words since %v to %v", context.Chat().Title, from.Format("02.01.2006"), to.Format("02.01.2006")), Theme: "shine"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithTitleOpts(opts.Title{
			Title: "Popular words",
		}))

	wc.AddSeries("Popular words", WCData).
		SetSeriesOptions(
			charts.WithWorldCloudChartOpts(
				opts.WordCloudChart{
					SizeRange: []float32{14, 80},
				}),
		)
	return wc
}
