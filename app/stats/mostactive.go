package stats

import (
	"time"

	"github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func MostActiveUsersTodayPieChart(from time.Time, to time.Time, context telebot.Context) *charts.Pie {
	result, _ := utils.DB.
		Table("`messages`, `users`").
		Select("`users`.first_name || ' ' || `users`.last_name AS FullName, COUNT(`messages`.`id`) as Messages").
		Where("`messages`.`user_id`=`users`.`id`").
		Where("`messages`.`chat_id`=?", -1001123405621).
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
		charts.WithTooltipOpts(opts.Tooltip{Show: true}),
		charts.WithTitleOpts(opts.Title{Title: "Most active users"}),
	)

	pie.AddSeries("Messages", Users)
	return pie
}
