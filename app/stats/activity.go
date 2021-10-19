package stats

import (
	"fmt"
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func UserActivityLineChart(from time.Time, to time.Time, context telebot.Context) *charts.Line {
	result, _ := utils.DB.
		Model(utils.Message{ChatID: context.Chat().ID}).
		Select("strftime('%d.%m',`date`, 'localtime') as Day, COUNT(DISTINCT `user_id`) AS Users, COUNT(`id`) as Messages").
		Where("date BETWEEN ? AND ?", from, to).
		Group("Day").
		Order("date").
		Rows()
	var Days []string
	var Users []opts.LineData
	var Messages []opts.LineData
	var Day string
	var UsersCount int
	var MessagesCount int
	var UsersMax int
	var MessagesMax int
	for result.Next() {
		err := result.Scan(&Day, &UsersCount, &MessagesCount)
		if err != nil {
			utils.ErrorReporting(err, context)
			return nil
		}
		Days = append(Days, Day)
		Users = append(Users, opts.LineData{Value: UsersCount, Name: fmt.Sprintf("%v", UsersCount), YAxisIndex: 1})
		Messages = append(Messages, opts.LineData{Value: MessagesCount, Name: fmt.Sprintf("%v", MessagesCount)})
		if UsersCount > UsersMax {
			UsersMax = UsersCount
		}
		if MessagesCount > MessagesMax {
			MessagesMax = MessagesCount
		}
	}

	// create a new line instance
	line := charts.NewLine()
	// set some global options like Title/Legend/ToolTip or anything else
	line.SetGlobalOptions(
		charts.WithLegendOpts(opts.Legend{Show: true}),
		charts.WithInitializationOpts(opts.Initialization{PageTitle: context.Chat().Title + " Chat Activity", Theme: "shine"}),
		charts.WithTitleOpts(opts.Title{
			Title: "User activity",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Name: "Messages",
			Type: "value",
			Show: true,
			Min:  0,
			Max:  MessagesMax + 100,
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Name: "Days",
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
	line.ExtendYAxis(opts.YAxis{
		Name: "Users",
		Type: "value",
		Show: true,
		Min:  0,
		Max:  UsersMax + 10,
	})

	// Put data into instance
	line.SetXAxis(Days).
		AddSeries("Users", Users, charts.WithLineChartOpts(opts.LineChart{Smooth: true, YAxisIndex: 1})).
		AddSeries("Messages", Messages, charts.WithLineChartOpts(opts.LineChart{Smooth: true})).
		SetSeriesOptions(
			charts.WithLabelOpts(opts.Label{Show: false}),
		)
	return line
}
