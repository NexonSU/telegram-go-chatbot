package commands

import (
	"github.com/NexonSU/telegram-go-chatbot/app/utils"
	"gopkg.in/tucnak/telebot.v3"
)

//Send pidor rules on /pidorules
func Pidorules(context telebot.Context) error {
	var err error
	if context.Chat().Username != utils.Config.Telegram.Chat && !utils.IsAdminOrModer(context.Sender().Username) {
		return err
	}
	err := context.Reply("Правила игры <b>Пидор Дня</b>:\n<b>1.</b> Зарегистрируйтесь в игру по команде /pidoreg\n<b>2.</b> Подождите пока зарегиструются все (или большинство :)\n<b>3.</b> Запустите розыгрыш по команде /pidor\n<b>4.</b> Просмотр статистики канала по команде /pidorstats, /pidorall\n<b>5.</b> Личная статистика по команде /pidorme\n<b>6. (!!! Только для администраторов чатов)</b>: удалить из игры может только Админ канала, сначала выведя по команде список игроков: /utils.PidorList (список упадёт в личку)\nУдалить же игрока можно по команде (используйте идентификатор пользователя - цифры из списка пользователей): /pidordel {ID или никнейм юзера}\nТак же, удалить можно просто отправив /pidordel в ответ на сообщение пользователя, которого нужно удалить из игры.\n\nВажно, розыгрыш проходит только раз в день, повторная команда выведет <b>результат</b> игры.\n\nСброс розыгрыша происходит каждый день ночью.\n\nПоддержать автора оригинального бота можно по <a href=\"https://www.paypal.me/unicott/2\">ссылке</a> :)")
	if err != nil {
		return err
	}
	return err
}
