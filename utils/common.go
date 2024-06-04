package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	cntx "context"

	"github.com/neonxp/StemmerRu"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	tele "gopkg.in/telebot.v3"
	ffprobe "gopkg.in/vansante/go-ffprobe.v2"
	"gorm.io/gorm/clause"
)

var onlyWords = regexp.MustCompile(`[,.!?]+`)

func UserFullName(user *tele.User) string {
	fullname := user.FirstName
	if user.LastName != "" {
		fullname = fmt.Sprintf("%v %v", user.FirstName, user.LastName)
	}
	return fullname
}

func UserName(user *tele.User) string {
	username := user.Username
	if user.Username == "" {
		username = UserFullName(user)
	}
	return username
}

func MentionUser(user *tele.User) string {
	return fmt.Sprintf("<a href=\"tg://user?id=%v\">%v</a>", user.ID, UserFullName(user))
}

func RandInt(min int, max int) int {
	b, err := rand.Int(rand.Reader, big.NewInt(int64(max+1)))
	if err != nil {
		return 0
	}
	return min + int(b.Int64())
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func IsAdmin(userid int64) bool {
	for _, b := range Config.Admins {
		if b == userid {
			return true
		}
	}
	return false
}

func IsAdminOrModer(userid int64) bool {
	for _, b := range Config.Admins {
		if b == userid {
			return true
		}
	}
	for _, b := range Config.Moders {
		if b == userid {
			return true
		}
	}
	return false
}

func RestrictionTimeMessage(seconds int64) string {
	var message = ""
	if seconds-30 > time.Now().Unix() {
		message = fmt.Sprintf(" до %v", time.Unix(seconds, 0).Format("02.01.2006 15:04:05"))
	}
	return message
}

func FindUserInMessage(context tele.Context) (tele.User, int64, error) {
	var user tele.User
	var err error = nil
	var untildate = time.Now().Unix() + 86400
	for _, entity := range context.Message().Entities {
		if entity.Type == tele.EntityTMention {
			user = *entity.User
			if len(context.Args()) == 2 {
				addtime, err := strconv.ParseInt(context.Args()[1], 10, 64)
				if err != nil {
					return user, untildate, err
				}
				untildate += addtime - 86400
			}
			return user, untildate, err
		}
	}
	if context.Message().ReplyTo != nil {
		user = *context.Message().ReplyTo.Sender
		if len(context.Args()) == 1 {
			addtime, err := strconv.ParseInt(context.Args()[0], 10, 64)
			if err != nil {
				return user, untildate, errors.New("время указано неверно")
			}
			untildate += addtime - 86400
		}
	} else {
		if len(context.Args()) == 0 {
			err = errors.New("пользователь не найден")
			return user, untildate, err
		}
		if user.ID == 0 {
			user, err = GetUserFromDB(context.Args()[0])
			if err != nil {
				return user, untildate, err
			}
		}
		if len(context.Args()) == 2 {
			addtime, err := strconv.ParseInt(context.Args()[1], 10, 64)
			if err != nil {
				return user, untildate, errors.New("время указано неверно")
			}
			untildate += addtime - 86400
		}
	}
	return user, untildate, err
}

func GetUserFromDB(findstring string) (tele.User, error) {
	var user tele.User
	var err error = nil
	if string(findstring[0]) == "@" {
		user.Username = findstring[1:]
	} else {
		user.ID, err = strconv.ParseInt(findstring, 10, 64)
	}
	result := DB.Where("lower(username) = ? OR id = ?", strings.ToLower(user.Username), user.ID).First(&user)
	if result.Error != nil {
		err = result.Error
	}
	return user, err
}

// Forward channel post to chat
func ForwardPost(context tele.Context) error {
	if context.Message() == nil || context.Message().Chat.ID != Config.Channel {
		return nil
	}
	var err error
	if context.Message().Text != "" || context.Message().Caption != "" {
		_, err = Bot.Forward(&tele.Chat{ID: Config.Chat}, context.Message())
	}
	if Config.StreamChannel != 0 {
		if strings.Contains(context.Text(), "zavtracast/live") {
			_, err = Bot.Forward(&tele.Chat{ID: Config.StreamChannel}, context.Message())
			return err
		}
		for _, entity := range append(context.Message().CaptionEntities, context.Message().Entities...) {
			if entity.Type == tele.EntityURL || entity.Type == tele.EntityTextLink {
				if strings.Contains(entity.URL, "zavtracast/live") {
					_, err = Bot.Forward(&tele.Chat{ID: Config.StreamChannel}, context.Message())
					return err
				}
			}
		}
	}
	return err
}

// Remove message
func Remove(context tele.Context) error {
	return context.Delete()
}

func OnChatMember(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return Bot.Unban(&tele.Chat{ID: context.Chat().ID}, context.ChatMember().NewChatMember.User)
	}
	//User update
	UserResult := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(context.ChatMember().NewChatMember.User)
	if UserResult.Error != nil {
		ErrorReporting(UserResult.Error, nil)
	}
	return nil
}

func OnUserJoined(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return context.Delete()
	}
	return nil
}

func OnUserLeft(context tele.Context) error {
	if context.Chat().ID == Config.ReserveChat {
		return context.Delete()
	}
	return nil
}

func OnText(context tele.Context) error {
	//remove message from reservechat
	if context.Chat().ID == Config.ReserveChat {
		return context.Delete()
	}

	//User update
	UserResult := DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(context.Sender())
	if UserResult.Error != nil {
		ErrorReporting(UserResult.Error, nil)
	}

	//update StatsDays(1), StatsHours(2), StatsUsers(3), StatsWords(4), StatsWeekday(5)
	startOfDay := GetStartOfDay()
	timeNow := time.Now().Local()
	statsIncrease(1, startOfDay, int64(timeNow.Day()))
	statsIncrease(2, startOfDay, int64(timeNow.Hour()))
	statsIncrease(3, startOfDay, context.Sender().ID)
	statsIncrease(5, startOfDay, int64(timeNow.Weekday()))
	if context.Text() == "" || string(context.Text()[0]) == "/" {
		return nil
	}
	text := strings.ToLower(onlyWords.ReplaceAllString(context.Text(), ""))
	for _, word := range strings.Split(text, " ") {
		if len([]rune(word)) > 0 {
			statsIncrease(4, startOfDay, getWordID(word))
		}
	}
	return nil
}

func statsIncrease(statType int64, dayTimestamp int64, contextID int64) {
	if DB.Exec("UPDATE stats SET count = count + 1, last_update = ? WHERE context_id = ? AND stat_type = ? AND day_timestamp = ?;", time.Now().Local().Unix(), contextID, statType, dayTimestamp).RowsAffected == 0 {
		DB.Create(Stats{StatType: statType, DayTimestamp: dayTimestamp, ContextID: contextID, Count: 1, LastUpdate: time.Now().Local().Unix()})
	}
}

func getWordID(searchWord string) int64 {
	shortWord := StemmerRu.Stem(searchWord)
	wordResult := StatsWords{}
	if DB.Model(StatsWords{}).Select("id").Where("short_word = ?", shortWord).Find(&wordResult).RowsAffected == 0 {
		wordResult.ShortWord = shortWord
		wordResult.Word = searchWord
		DB.Create(&wordResult)
	}
	return wordResult.ID
}

func GetStartOfDay() int64 {
	unixTS := time.Now().Local().Unix()
	tm := time.Unix(unixTS, 0).In(time.Local)
	hour, minute, second := tm.Clock()

	return unixTS - int64(hour*3600+minute*60+second)
}

func GetNope() string {
	var nope Nope
	DB.Model(Nope{}).Order("RANDOM()").First(&nope)
	return nope.Text
}

func GetBless() string {
	var bless Bless
	DB.Model(Bless{}).Order("RANDOM()").First(&bless)
	return bless.Text
}

func FFmpegConvert(context tele.Context, filePath string, targetType string) error {
	var KwArgs ffmpeg.KwArgs
	var extension string
	var height int
	var width int
	var duration float64
	var err error

	videoKwArgs := ffmpeg.KwArgs{}
	defaultKwArgs := ffmpeg.KwArgs{"loglevel": "fatal", "hide_banner": ""}
	name := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	ctx, cancelFn := cntx.WithTimeout(cntx.Background(), 5*time.Second)
	defer cancelFn()

	data, err := ffprobe.ProbeURL(ctx, filePath)
	if err != nil {
		return err
	}
	inputVideoFormat := data.FirstVideoStream()
	inputAudioFormat := data.FirstAudioStream()

	if inputVideoFormat != nil {
		width = inputVideoFormat.Width
		height = inputVideoFormat.Height
		duration, err = strconv.ParseFloat(inputVideoFormat.Duration, 32)
		if err != nil {
			duration = 0
		}
	} else {
		switch targetType {
		case "animation", "gif", "webm", "animation_reverse", "sticker_reverse", "animation_loop", "loop":
			return fmt.Errorf("видео-дорожка не найдена")
		case "video", "mp4":
			targetType = "audio"
		case "video_reverse", "reverse", "invert":
			targetType = "audio_reverse"
		}
	}

	if inputAudioFormat != nil {
		if duration == 0 {
			duration, err = strconv.ParseFloat(inputAudioFormat.Duration, 32)
			if err != nil {
				duration = 0
			}
		}
	} else {
		switch targetType {
		case "audio", "mp3", "voice", "ogg", "audio_reverse", "voice_reverse", "audio_loop", "voice_loop":
			return fmt.Errorf("аудио-дорожка не найдена")
		case "video", "mp4", "webm":
			targetType = "animation"
		case "video_reverse", "reverse", "invert":
			targetType = "animation_reverse"
		case "video_loop", "loop":
			targetType = "animation_loop"
		}
	}

	if (strings.Contains(targetType, "reverse") || strings.Contains(targetType, "loop")) && duration > 120 {
		return fmt.Errorf("слишком длинное видео для эффекта")
	}

	if inputAudioFormat == nil && inputVideoFormat == nil {
		return fmt.Errorf("медиа-дорожек не найдено")
	}

	switch targetType {
	case "audio", "mp3":
		KwArgs = ffmpeg.KwArgs{"vn": "", "c:a": "libmp3lame"}
		extension = "mp3"
		targetType = "audio"
	case "voice", "ogg":
		KwArgs = ffmpeg.KwArgs{"vn": "", "c:a": "libopus"}
		extension = "ogg"
		targetType = "voice"
	case "photo", "jpg":
		KwArgs = ffmpeg.KwArgs{"vf": "select=eq(n\\,0)", "format": "image2"}
		extension = "jpg"
		targetType = "photo"
	case "sticker", "webp":
		KwArgs = ffmpeg.KwArgs{"vf": "select=eq(n\\,0)", "format": "image2"}
		extension = "webp"
		targetType = "sticker"
	case "animation", "gif":
		KwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{videoKwArgs, {"an": ""}})
		extension = "mp4"
		targetType = "animation"
	case "video", "mp4":
		KwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{videoKwArgs, {"c:a": "aac"}})
		extension = "mp4"
		targetType = "video"
	case "video_reverse", "reverse", "invert":
		KwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{videoKwArgs, {"c:a": "aac", "vf": "reverse", "af": "areverse"}})
		extension = "mp4"
		targetType = "video"
	case "animation_reverse":
		KwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{videoKwArgs, {"an": "", "vf": "reverse"}})
		extension = "mp4"
		targetType = "animation"
	case "sticker_reverse", "webm":
		KwArgs = ffmpeg.KwArgs{"c:v": "libvpx-vp9", "an": "", "vf": "reverse"}
		extension = "mp4"
		targetType = "sticker"
	case "audio_reverse":
		KwArgs = ffmpeg.KwArgs{"vn": "", "c:a": "libmp3lame", "af": "areverse"}
		extension = "mp3"
		targetType = "audio"
	case "voice_reverse":
		KwArgs = ffmpeg.KwArgs{"vn": "", "c:a": "libopus", "af": "areverse"}
		extension = "ogg"
		targetType = "voice"
	case "animation_loop", "loop":
		framesInt, err := strconv.Atoi(inputVideoFormat.NbFrames)
		if err != nil {
			return err
		}
		KwArgs = ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{videoKwArgs, {"an": "", "filter_complex": fmt.Sprintf("[0]trim=start_frame=1:end_frame=%d,setpts=PTS-STARTPTS,reverse[r];[0][r]concat,loop=1:2", framesInt-1)}})
		extension = "mp4"
		targetType = "animation"
	default:
		return fmt.Errorf("targetType %v not supported", targetType)
	}

	resultFilePath := fmt.Sprintf("%v/%v_converted.%v", os.TempDir(), name, extension)

	err = ffmpeg.Input(filePath).Output(resultFilePath, ffmpeg.MergeKwArgs([]ffmpeg.KwArgs{defaultKwArgs, KwArgs})).OverWriteOutput().ErrorToStdOut().Run()
	if err != nil {
		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		return err
	}

	defer func(resultFilePath string) {
		os.Remove(resultFilePath)
	}(resultFilePath)

	switch targetType {
	case "video":
		return context.Reply(&tele.Video{
			File:      tele.FromDisk(resultFilePath),
			FileName:  filepath.Base(resultFilePath),
			Streaming: true,
			Width:     width,
			Height:    height,
			Duration:  int(duration),
			MIME:      "video/mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	case "animation":
		return context.Reply(&tele.Animation{
			File:     tele.FromDisk(resultFilePath),
			FileName: filepath.Base(resultFilePath),
			Width:    width,
			Height:   height,
			Duration: int(duration),
			MIME:     "video/mp4",
		}, &tele.SendOptions{AllowWithoutReply: true})
	case "audio":
		return context.Reply(&tele.Audio{
			File:     tele.FromDisk(resultFilePath),
			FileName: filepath.Base(resultFilePath),
			MIME:     "audio/mp3",
		}, &tele.SendOptions{AllowWithoutReply: true})
	case "voice":
		return context.Reply(&tele.Voice{
			File: tele.FromDisk(resultFilePath),
			MIME: "audio/ogg",
		}, &tele.SendOptions{AllowWithoutReply: true})
	default:
		return context.Reply(&tele.Document{
			File:     tele.FromDisk(resultFilePath),
			FileName: filepath.Base(resultFilePath),
		}, &tele.SendOptions{AllowWithoutReply: true})
	}
}

func DownloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func ReplyAndRemove(message string, context tele.Context) error {
	message += "\n\nЭто сообщение самоуничтожится через 30 секунд."
	sentMessage, err := Bot.Send(context.Chat(), message, &tele.SendOptions{ReplyTo: context.Message(), AllowWithoutReply: true})
	if err != nil {
		return err
	}
	if context.Message() == nil || context.Message().ID == 0 {
		return nil
	}
	go func(messages []int) {
		time.Sleep(30 * time.Second)
		for _, message := range messages {
			Bot.Delete(&tele.Message{ID: message, Chat: &tele.Chat{ID: context.Chat().ID}})
		}
	}([]int{context.Message().ID, sentMessage.ID})
	return nil
}
