package stats

import (
	"time"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	tele "gopkg.in/telebot.v3"
)

func MostActiveUsersTodayPieChart(from time.Time, to time.Time, context tele.Context) *charts.Pie {
	result, _ := utils.DB.
		Table("`messages`, `users`").
		Select("`users`.first_name || ' ' || `users`.last_name AS FullName, COUNT(`messages`.`id`) as Messages").
		Where("`messages`.`user_id`=`users`.`id`").
		Where("`messages`.`chat_id`=?", context.Chat().ID).
		Where("date BETWEEN ? AND ?", time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local), time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 0, time.Local)).
		Group("`users`.`id`").
		Order("messages DESC").
		Limit(20).
		Rows()
	var Users []opts.PieData
	var FullName string
	var MessagesCount int
	for result.Next() {
		err := result.Scan(&FullName, &MessagesCount)
		if err != nil {
			utils.ErrorReporting(err, context)
			return nil
		}
		Users = append(Users, opts.PieData{Value: MessagesCount, Name: FullName})
	}

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: context.Chat().Title + " Most Active Users of Day", Theme: "shine"}),
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithTitleOpts(opts.Title{Title: "Most active users today", Left: "center"}),
	)

	pie.AddSeries("Messages", Users)
	return pie
}
