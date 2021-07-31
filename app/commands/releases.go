package commands

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	ical "github.com/arran4/golang-ical"
	"gopkg.in/tucnak/telebot.v3"
)

//Send releases of 2 weeks on /releases
func Releases(context telebot.Context) error {
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
	for _, element := range cal.Events() {
		date := element.GetProperty(ical.ComponentPropertyDtStart).Value
		name := element.GetProperty(ical.ComponentPropertySummary).Value
		dateint, _ := strconv.Atoi(date)
		if dateint > today && dateint < twoweeks {
			releases = fmt.Sprintf("<b>%v</b> - %v.%v.%v\n%v", strings.ReplaceAll(name, "\\,", ","), date[6:8], date[4:6], date[0:4], releases)
		}
	}
	return context.Reply(releases)
}
