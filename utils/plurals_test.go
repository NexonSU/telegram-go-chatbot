package utils

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func TestPlurals(t *testing.T) {
	prt := message.NewPrinter(language.Russian)
	require.NotNil(t, prt)

	msg := prt.Sprintf("<code>💥 %v %v%v.\nРеспавн через %d мин.</code>", "Fake Sender", "", "Fake Reason", 5)
	assert.Contains(t, msg, "5 минут")
	msg = prt.Sprintf("<code>💥 %v %v%v.\nРеспавн через %d мин.</code>", "Fake Sender", "", "Fake Reason", 42)
	assert.Contains(t, msg, "42 минуты")
	msg = prt.Sprintf("<code>💥 %v %v%v.\nРеспавн через %d мин.</code>", "Fake Sender", "", "Fake Reason", 451)
	assert.Contains(t, msg, "451 минуту")

	msg = prt.Sprintf("💥 %v %vпристрелил %v.\n%v отправился на респавн на %d мин.", "Fake Admin", "", "Fake User", "Fake User", 8)
	assert.Contains(t, msg, "8 минут")
	msg = prt.Sprintf("💥 %v %vпристрелил %v.\n%v отправился на респавн на %d мин.", "Fake Admin", "", "Fake User", "Fake User", 42)
	assert.Contains(t, msg, "42 минуты")
	msg = prt.Sprintf("💥 %v %vпристрелил %v.\n%v отправился на респавн на %d мин.", "Fake Admin", "", "Fake User", "Fake User", 451)
	assert.Contains(t, msg, "451 минуту")

	msg = prt.Sprintf("🤫 %v %vпопросил %v помолчать %d минут.", "Fake Admin", "", "Fake User", 15)
	assert.Contains(t, msg, "15 минут")
	msg = prt.Sprintf("🤫 %v %vпопросил %v помолчать %d минут.", "Fake Admin", "", "Fake User", 42)
	assert.Contains(t, msg, "42 минут")
	msg = prt.Sprintf("🤫 %v %vпопросил %v помолчать %d минут.", "Fake Admin", "", "Fake User", 451)
	assert.Contains(t, msg, "451 минуту")

	msg = prt.Sprintf("У тебя %d предупреждений.", 1)
	assert.Contains(t, msg, "1 предупреждение")
	msg = prt.Sprintf("У тебя %d предупреждений.", 2)
	assert.Contains(t, msg, "2 предупреждения")
	msg = prt.Sprintf("У тебя %d предупреждений.", 7)
	assert.Contains(t, msg, "7 предупреждений")

	msg = prt.Sprintf("%v😈 Наводит револьвер на %v и стреляет.\nЯ хз как это объяснить, но %v победитель!\n%v отправился на респавн на %d мин.",
		"Fake Admin", "Fake User", "Fake Admin", "Fake User", 16)
	assert.Contains(t, msg, "16 минут")
	msg = prt.Sprintf("%v😈 Наводит револьвер на %v и стреляет.\nЯ хз как это объяснить, но %v победитель!\n%v отправился на респавн на %d мин.",
		"Fake Admin", "Fake User", "Fake Admin", "Fake User", 23)
	assert.Contains(t, msg, "23 минуты")
	msg = prt.Sprintf("%v😈 Наводит револьвер на %v и стреляет.\nЯ хз как это объяснить, но %v победитель!\n%v отправился на респавн на %d мин.",
		"Fake Admin", "Fake User", "Fake Admin", "Fake User", 451)
	assert.Contains(t, msg, "451 минуту")

	msg = prt.Sprintf("%v\nПобедитель дуэли: %v.\n%v отправился на респавн на %d мин.",
		"Fake reason", "Fake Admin", "Fake User", 16)
	assert.Contains(t, msg, "16 минут")
	msg = prt.Sprintf("%v\nПобедитель дуэли: %v.\n%v отправился на респавн на %d мин.",
		"Fake reason", "Fake Admin", "Fake User", 23)
	assert.Contains(t, msg, "23 минуты")
	msg = prt.Sprintf("%v\nПобедитель дуэли: %v.\n%v отправился на респавн на %d мин.",
		"Fake reason", "Fake Admin", "Fake User", 451)
	assert.Contains(t, msg, "451 минуту")

	msg = prt.Sprintf("%d смертей", 1)
	assert.Contains(t, msg, "1 смерть")
	msg = prt.Sprintf("%d смертей", 2)
	assert.Contains(t, msg, "2 смерти")
	msg = prt.Sprintf("%d смертей", 7)
	assert.Contains(t, msg, "7 смертей")

	msg = prt.Sprintf("%d побед", 1)
	assert.Contains(t, msg, "1 победа")
	msg = prt.Sprintf("%d побед", 2)
	assert.Contains(t, msg, "2 победы")
	msg = prt.Sprintf("%d побед", 7)
	assert.Contains(t, msg, "7 побед")

	msg = prt.Sprintf("%v. %v - %d раз\n", "something", "something", 16)
	assert.Contains(t, msg, "16 раз")
	msg = prt.Sprintf("%v. %v - %d раз\n", "something", "something", 23)
	assert.Contains(t, msg, "23 раза")
	msg = prt.Sprintf("%v. %v - %d раз\n", "something", "something", 51)
	assert.Contains(t, msg, "51 раз")

	msg = prt.Sprintf("В этом году ты был пидором дня — %d раз", 16)
	assert.Contains(t, msg, "16 раз")
	msg = prt.Sprintf("В этом году ты был пидором дня — %d раз", 23)
	assert.Contains(t, msg, "23 раза")
	msg = prt.Sprintf("В этом году ты был пидором дня — %d раз", 51)
	assert.Contains(t, msg, "51 раз")

	msg = prt.Sprintf("За всё время ты был пидором дня — %d раз!", 16)
	assert.Contains(t, msg, "16 раз")
	msg = prt.Sprintf("За всё время ты был пидором дня — %d раз!", 23)
	assert.Contains(t, msg, "23 раза")
	msg = prt.Sprintf("За всё время ты был пидором дня — %d раз!", 51)
	assert.Contains(t, msg, "51 раз")

	msg = prt.Sprintf("\nВсего участников — %d", 16)
	assert.Contains(t, msg, "16 участников")
	msg = prt.Sprintf("\nВсего участников — %d", 23)
	assert.Contains(t, msg, "23 участника")
	msg = prt.Sprintf("\nВсего участников — %d", 51)
	assert.Contains(t, msg, "51 участник")
}
