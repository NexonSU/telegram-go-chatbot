package commands

import (
	"math/rand"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/NexonSU/telegram-go-chatbot/utils"
	tele "gopkg.in/telebot.v3"
	"gorm.io/gorm/clause"
)

var firstSuicide int64
var lastSuicide int64
var burst int
var lastVideoSent int64

//Kill user on /blessing, /suicide
func Blessing(context tele.Context) error {
	// prt will replace fmt package to format text according plurals defined in utils package
	// If no plural rule matched it will be ignored and processed as usual formatting
	prt := message.NewPrinter(language.Russian)

	err := context.Delete()
	if err != nil {
		return err
	}
	ChatMember, err := utils.Bot.ChatMemberOf(context.Chat(), context.Sender())
	if err != nil {
		return err
	}
	if ChatMember.Role == "administrator" || ChatMember.Role == "creator" {
		return context.Send(prt.Sprintf("<code>👻 %v возродился у костра.</code>", utils.UserFullName(context.Sender())))
	}
	var duelist utils.Duelist
	result := utils.DB.Model(utils.Duelist{}).Where(context.Sender().ID).First(&duelist)
	if result.RowsAffected == 0 {
		duelist.UserID = context.Sender().ID
		duelist.Kills = 0
		duelist.Deaths = 0
	}
	duelist.Deaths++
	result = utils.DB.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&duelist)
	if result.Error != nil {
		return result.Error
	}
	duration := utils.RandInt(1, duelist.Deaths+1)
	duration += 10
	prependText := ""
	additionalChance := int(time.Now().Unix() - lastSuicide)
	if additionalChance > 3600 {
		additionalChance = 3600
	}
	additionalChance = (3600 - additionalChance) / 360
	if utils.RandInt(0, 100) >= 90-additionalChance {
		duration = duration * 10
		prependText = "критически "
	}
	ChatMember.RestrictedUntil = time.Now().Add(time.Second * time.Duration(60*duration)).Unix()
	err = utils.Bot.Restrict(context.Chat(), ChatMember)
	if err != nil {
		return err
	}
	reason := []string{
		"выбрал лёгкий путь",
		"сыграл в ящик",
		"слил своё HP до нуля",
		"приказал долго жить",
		"покинул этот скорбный мир",
		"пагиб",
		"разбежавшись прыгнул со скалы",
		"разогнал RTX 4090 Ti",
		"принял ислам",
		"пьёт чай и кушоет конфеты, никакова суецыда",
		"намотался на столб",
		"помер від крінжі",
		"здох",
		"заплатил, а было бесплатно",
		"уехал в дурку",
		"нашёл себя в прошмандовках завтрачата",
		"разочаровал партию, минус 20 социальный кредит и кошкажена",
		"донёс на самого себя",
		"выпил йаду",
		"папил геймпасу",
		"отправился на цыганскую свадьбу",
		"отменил себя",
		"посмотрел на уточку",
		"погасил ебало",
		"сыграл в сабнавтику",
		"ушёл пить комфеты и кушоть чай",
		"хряпнул вишневой балтики",
		"поиграл в леммингов",
		"стал единым с обелиском",
		"встретил Орнштейна и Смоуга",
		"сел в поезд, а поезд сделал бум",
		"стоял в луже АОЕ",
		"получил привет от мистера Сальери",
		"в сделку не входил",
		"не заметил Сефирота",
		"молодец, не воспользовался Бехелитом",
		"не выплатил вовремя долг Нуку",
		"был пойман велоцираптором",
		"был раздавлен Metal Gear REX",
		"исекайнулся",
		"стал целью Агента 47",
		"обнял крипера",
		"разбил пробирку с Т-вирусом",
		"заблудился в туманном городе",
		"забыл, что двойного прыжка в жизни нет",
		"разозлил Кирю",
		"провалился под мир",
		"застрял в геометрии",
		"встретил геймбрейкинг баг",
		"жрал капусту, когда есть картошка",
		"спросил \"А что случилось?\"",
		"наступил на лего",
		"не попал в QTE",
		"был пойман конторой пидорасов",
		"пошел с Романом в боулинг",
		"поверил, что GLaDOS даст тортик",
		"забыл основы CQC",
		"осознал весь сюжет Kingdom Hearts",
		"был прибит самым слабеньким и глупеньким мобом",
		"показал, что может без рук",
		"ушел бастурмировать",
		"был намотан на катамари",
		"охладил траханье",
		"попал в межсезонье",
		"застрял в вентиляции",
		"получил стрелу в колено",
		"совершил равноценный обмен",
		"перепутал красный и синий провод",
		"ушёл смотреть Free!",
		"приставил пистолет к виску и крикнул PERUSONA",
		"приставил пистолет к виску и попытался призвать персону",
		"ушёл искать 228922",
		"Ⓘ Данное сообщение доступно только для пользователей с подпиской Telegram Premium",
		"пил под вишнями компот, лишь на миг он отвернулся, на него упал дроп под",
		"добухтелся",
		"повернул на ультраправо",
		"превратился в дакимакуру",
		"сказал что религия - самый скучный фандом",
		"получил пизды от Олега Тинькова",
		"оказался фанатом Феррари",
		"попытался убраться дома, а потом понял что the biggest garbage - он сам",
	}
	burst++
	if time.Now().Unix() > firstSuicide+120 {
		firstSuicide = time.Now().Unix()
		burst = 1
	}
	lastSuicide = time.Now().Unix()
	if burst > 3 && time.Now().Unix() > lastVideoSent+3600 {
		lastVideoSent = time.Now().Unix()
		return context.Send(&tele.Video{
			File: tele.File{
				FileID: "BAACAgIAAx0CReJGYgABAlMuYnagTilFaB8ke8Rw-dYLbfJ6iF8AAicYAAIlxrlLY9ah2fUtR40kBA",
			},
			Caption: prt.Sprintf("<code>💥 %v %v%v.\nРеспавн через %d мин.</code>", utils.UserFullName(context.Sender()), prependText, reason[rand.Intn(len(reason))], duration),
		})
	} else {
		return context.Send(prt.Sprintf("<code>💥 %v %v%v.\nРеспавн через %d мин.</code>", utils.UserFullName(context.Sender()), prependText, reason[rand.Intn(len(reason))], duration))
	}
}
