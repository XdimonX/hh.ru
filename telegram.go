package main

import (
	// "fmt"
	"strconv"
	"strings"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

const (
	teleAdminID = 385060683
)

func startBot() {
	bot, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	checkErr(err)
	teleAdminUser := tb.User{ID: teleAdminID}
	bot.Handle("/help", func(m *tb.Message) {
		if m.Sender.ID == teleAdminID {
			msg := `loginHHru=<Задать логин от сайта hh.ru>
passwordHHru=<Задать пароль от сайта hh.ru>
timeoutResumeUpdate=<Установить частоту обновления резюме (в минутах)>
getResume Получить список резюме
getLoginHHru Получить логин на hh.ru
getTimeoutResumeUpdate Получить тайм-аут`
			bot.Send(m.Sender, msg)
		}
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		if m.Sender.ID == teleAdminID {
			if strings.HasPrefix(strings.ToLower(m.Text), "loginhhru") {
				saveLoginHHru(m, bot)
			} else if strings.HasPrefix(strings.ToLower(m.Text), "passwordhhru") {
				savePasswordHHru(m, bot)
			} else if strings.HasPrefix(strings.ToLower(m.Text), "timeoutresumeupdate") {
				saveTimeoutResumeUpdate(m, bot)
			} else if strings.HasPrefix(strings.ToLower(m.Text), "getresume") {
			} else if strings.HasPrefix(strings.ToLower(m.Text), "getloginhhru") {
				lock.Lock()
				bot.Send(m.Sender, loginHHru)
				lock.Unlock()
			} else if strings.HasPrefix(strings.ToLower(m.Text), "gettimeoutresumeupdate") {
				lock.Lock()
				bot.Send(m.Sender, strconv.Itoa(timeoutResumeUpdate))
				lock.Unlock()
			} else {
				bot.Send(m.Sender, "Не верная команда")
			}
		}
	})
	bot.Send(&teleAdminUser, "Запуск - "+time.Now().Format(time.ANSIC))
	bot.Start()
}

func saveTimeoutResumeUpdate(m *tb.Message, bot *tb.Bot) {
	text := strings.Split(m.Text, "=")
	if len(text) == 2 {
		lock.Lock()
		var err error
		timeoutResumeUpdate, err = strconv.Atoi(strings.TrimSpace(text[1]))
		if err != nil {
			bot.Send(m.Sender, "Ожидается число")
			lock.Unlock()
			return
		}
		saveCfg()
		lock.Unlock()
		bot.Send(m.Sender, "Тайм-аут успешно сохранён")
	} else {
		bot.Send(m.Sender, "Не верная команда")
	}
}

func saveLoginHHru(m *tb.Message, bot *tb.Bot) {
	text := strings.Split(m.Text, "=")
	if len(text) == 2 {
		lock.Lock()
		loginHHru = strings.TrimSpace(text[1])
		saveCfg()
		lock.Unlock()
		bot.Send(m.Sender, "Логин успешно сохранён")
	} else {
		bot.Send(m.Sender, "Не верная команда")
	}
}

func savePasswordHHru(m *tb.Message, bot *tb.Bot) {
	text := strings.Split(m.Text, "=")
	if len(text) == 2 {
		lock.Lock()
		passwordHHru = strings.TrimSpace(text[1])
		saveCfg()
		lock.Unlock()
		bot.Send(m.Sender, "Пароль успешно сохранён")
	} else {
		bot.Send(m.Sender, "Не верная команда")
	}
}

// func tst() {
// 	b, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
// 	checkErr(err)
// 	teleAdminUser := tb.User{ID: teleAdminID}
// 	// var (
// 	// 	menu    = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
// 	// 	btnHelp = menu.Text("ℹ Help")
// 	// )
// 	// menu.Reply(
// 	// 	menu.Row(btnHelp),
// 	// )
// 	// b.Handle(&btnHelp, func(m *tb.Message) {
// 	// 	b.Send(m.Chat, "Тута будет помощь.")
// 	// })
// 	// b.Handle("/help", func(m *tb.Message) {
// 	// 	// b.Send(m.Chat, "test", menu)
// 	// })
// 	b.Send(&teleAdminUser, "Запуск - "+time.Now().Format(time.ANSIC))
// 	b.Start()
// }
