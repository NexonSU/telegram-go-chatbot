package stats

import (
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func PopWordsWcChart(from time.Time, to time.Time, context telebot.Context) *charts.WordCloud {
	result, _ := utils.DB.
		Model(utils.Word{ChatID: -1001123405621}).
		Select("text, COUNT(*) as count").
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
		charts.WithTitleOpts(opts.Title{
			Title: "Popular words",
		}))

	wc.AddSeries("wordcloud", WCData).
		SetSeriesOptions(
			charts.WithWorldCloudChartOpts(
				opts.WordCloudChart{
					SizeRange: []float32{14, 80},
				}),
		)
	return wc
}
