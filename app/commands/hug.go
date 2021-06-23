package commands

import (
	"bytes"
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/chai2010/webp"
	"github.com/fogleman/gg"
	tb "gopkg.in/tucnak/telebot.v2"
	"path/filepath"
	"runtime"
)

//Write username on hug picture and send to target
func Hug(m *tb.Message) {
	if m.ReplyTo == nil {
		_, err := utils.Bot.Reply(m, "Просто отправь <code>/hug</code> в ответ на чье-либо сообщение.")
		if err != nil {
			utils.ErrorReporting(err, m)
			return
		}
		return

	}
	var target = *m.ReplyTo
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	im, err := webp.Load(basepath + "/../../files/hug.webp")
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	dc := gg.NewContextForImage(im)
	dc.DrawImage(im, 0, 0)
	dc.Rotate(gg.Radians(15))
	dc.SetRGB(0, 0, 0)
	err = dc.LoadFontFace(basepath+"/../../files/impact.ttf", 20)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	dc.SetRGB(1, 1, 1)
	s := utils.UserFullName(m.Sender)
	n := 4
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				continue
			}
			x := 400 + float64(dx)
			y := -30 + float64(dy)
			dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
		}
	}
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(s, 400, -30, 0.5, 0.5)
	buf := new(bytes.Buffer)
	err = webp.Encode(buf, dc.Image(), nil)
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
	_, err = utils.Bot.Reply(&target, &tb.Sticker{File: tb.FromReader(buf)})
	if err != nil {
		utils.ErrorReporting(err, m)
		return
	}
}
