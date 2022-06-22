package commands

import (
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	ical "github.com/arran4/golang-ical"
)

//Send releases of 2 weeks on /releases
func Releases(context tele.Context) error {
	var err error
	if utils.Config.ReleasesUrl == "" {
		return context.Reply("Список ближайших релизов не настроен")
	}
	resp, err := http.Get(utils.Config.ReleasesUrl)
	if err != nil {
		return err
	}
	cal, err := ical.ParseCalendar(resp.Body)
	if err != nil {
		return err
	}
	releases := ""
	today, _ := strconv.Atoi(time.Now().Format("20060102"))
	twoweeks, _ := strconv.Atoi(time.Now().AddDate(0, 0, 14).Format("20060102"))
	events := cal.Events()
	sort.SliceStable(events, func(i, j int) bool {
		return events[i].GetProperty(ical.ComponentPropertyDtStart).Value > events[j].GetProperty(ical.ComponentPropertyDtStart).Value
	})
	for _, element := range events {
		date := element.GetProperty(ical.ComponentPropertyDtStart).Value
		name := element.GetProperty(ical.ComponentPropertySummary).Value
		dateint, _ := strconv.Atoi(date)
		if dateint > today && dateint < twoweeks {
			releases = fmt.Sprintf("<b>%v</b> - %v.%v.%v\n%v", strings.ReplaceAll(name, "\\,", ","), date[6:8], date[4:6], date[0:4], releases)
		}
	}
	return context.Reply(releases)
}
