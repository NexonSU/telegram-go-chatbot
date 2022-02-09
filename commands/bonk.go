package commands

import (
	"bytes"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	"github.com/chai2010/webp"
	"github.com/fogleman/gg"
	"gopkg.in/telebot.v3"
)

//Write username on bonk picture and send to target
func Bonk(context telebot.Context) error {
	if context.Message().ReplyTo == nil {
		return context.Reply("Просто отправь <code>/bonk</code> в ответ на чье-либо сообщение.")
	}
	context.Delete()
	im, err := webp.Load("files/bonk.webp")
	if err != nil {
		return err
	}
	dc := gg.NewContextForImage(im)
	dc.DrawImage(im, 0, 0)
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
			x := 140 + float64(dx)
			y := 290 + float64(dy)
			dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
		}
	}
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored(s, 140, 290, 0.5, 0.5)
	buf := new(bytes.Buffer)
	err = webp.Encode(buf, dc.Image(), nil)
	if err != nil {
		return err
	}
	return context.Send(&telebot.Sticker{File: telebot.FromReader(buf)}, &telebot.SendOptions{ReplyTo: context.Message().ReplyTo})
}
