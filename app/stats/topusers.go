package stats

import (
	"fmt"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"gopkg.in/tucnak/telebot.v3"
)

func TopUsersBarChart(from time.Time, to time.Time, context telebot.Context) *charts.Bar {
	result, _ := utils.DB.
		Table("`messages`, `users`").
		Select("`users`.first_name || ' ' || `users`.last_name AS FullName, COUNT(`messages`.`id`) as Messages").
		Where("`messages`.`user_id`=`users`.`id`").
		Where("`messages`.`chat_id`=?", context.Chat().ID).
		Where("date BETWEEN ? AND ?", from, to).
		Group("`users`.`id`").
		Order("messages DESC").
		Limit(20).
		Rows()
	var FullName string
	var Messages int
	var Users []string
	var UsersData []opts.BarData
	for result.Next() {
		err := result.Scan(&FullName, &Messages)
		if err != nil {
			utils.ErrorReporting(err, context)
			return nil
		}
		Users = append(Users, FullName)
		UsersData = append(UsersData, opts.BarData{Name: FullName, Value: Messages})
	}

	for i, j := 0, len(Users)-1; i < j; i, j = i+1, j-1 {
		Users[i], Users[j] = Users[j], Users[i]
		UsersData[i], UsersData[j] = UsersData[j], UsersData[i]
	}

	// create a new bar instance
	bar := charts.NewBar()
	// set some global options like Title/Legend/ToolTip or anything else
	bar.SetGlobalOptions(
		charts.WithInitializationOpts(opts.Initialization{PageTitle: fmt.Sprintf("%v Top Users since %v to %v", context.Chat().Title, from.Format("02.01.2006"), to.Format("02.01.2006")), Theme: "shine"}),
		charts.WithTitleOpts(opts.Title{
			Title: "Top users",
		}),
		charts.WithTooltipOpts(opts.Tooltip{
			Show:    true,
			Trigger: "axis",
			AxisPointer: &opts.AxisPointer{
				Type: "cross",
				Snap: true,
			},
		}),
		charts.WithXAxisOpts(opts.XAxis{
			Type: "value",
		}),
		charts.WithYAxisOpts(opts.YAxis{
			Type: "category",
			Data: Users,
		}),
	)

	bar.AddSeries("Messages", UsersData)
	return bar
}
