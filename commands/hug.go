package commands

import (
	"bytes"

	tele "github.com/NexonSU/telebot"
	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/chai2010/webp"
	"github.com/fogleman/gg"
)

//Write username on hug picture and send to target
func Hug(context tele.Context) error {
	var err error
	if context.Message().ReplyTo == nil {
		return context.Reply("Просто отправь <code>/hug</code> в ответ на чье-либо сообщение.")
	}
	context.Delete()
	im, err := webp.Load("files/hug.webp")
	if err != nil {
		return err
	}
	dc := gg.NewContextForImage(im)
	dc.DrawImage(im, 0, 0)
	dc.Rotate(gg.Radians(15))
	dc.SetRGB(0, 0, 0)
	err = dc.LoadFontFace("files/impact.ttf", 20)
	if err != nil {
		return err
	}
	dc.SetRGB(1, 1, 1)
	s := utils.UserFullName(context.Sender())
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
		return err
	}
	return context.Send(&tele.Sticker{File: tele.FromReader(buf)}, &tele.SendOptions{ReplyTo: context.Message().ReplyTo})
}
