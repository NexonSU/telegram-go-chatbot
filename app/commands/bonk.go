package commands

import (
	"bytes"
	"path/filepath"
	"runtime"

	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"github.com/chai2010/webp"
	"github.com/fogleman/gg"
	"gopkg.in/tucnak/telebot.v3"
)

//Write username on bonk picture and send to target
func Bonk(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	if context.Message().ReplyTo == nil {
		err := context.Reply("Просто отправь <code>/bonk</code> в ответ на чье-либо сообщение.")
		if err != nil {
			return err
		}
		return err

	}
	var target = *context.Message().ReplyTo
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	im, err := webp.Load(basepath + "/../../files/bonk.webp")
	if err != nil {
		return err
	}
	dc := gg.NewContextForImage(im)
	dc.DrawImage(im, 0, 0)
	dc.SetRGB(0, 0, 0)
	err = dc.LoadFontFace(basepath+"/../../files/impact.ttf", 20)
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
	_, err = utils.Bot.Reply(&target, &telebot.Sticker{File: telebot.FromReader(buf)})
	if err != nil {
		return err
	}
	return err
}
