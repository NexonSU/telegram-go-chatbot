package stats

import (
	"fmt"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func PopDaysBarChart(from time.Time, to time.Time, context telebot.Context) *charts.Bar {
	result, _ := utils.DB.
		Model(utils.Message{ChatID: -1001123405621}).
		Select("strftime('%w', `DATE`, 'localtime') AS Weekdays, COUNT(`id`) as Messages").
		Where("date BETWEEN ? AND ?", from, to).
		Group("Weekdays").
		Order("Weekdays").
		Rows()
	var Weekdays []string
	var Messages []opts.BarData
	var Weekday string
	var MessagesCount int
	for result.Next() {
		err := result.Scan(&Weekday, &MessagesCount)
		if err != nil {
			utils.ErrorReporting(err, context)
			return nil
		}
		Weekdays = append(Weekdays, Weekday)
		Messages = append(Messages, opts.BarData{Value: MessagesCount, Name: fmt.Sprintf("%v", MessagesCount)})
	}

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{
			Title: "Popular days of week",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
				Snap: true,
			},
		}),
	)

	// Put data into instance
	bar.SetXAxis(Weekdays).AddSeries("Messages", Messages)
	return bar
}
