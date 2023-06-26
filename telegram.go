// Модуль реализцющий работу телеграм бота, через которого происходит управление работой программы.
package main

import (
	// "fmt"
	"context"
	"strconv"
	"strings"
	"time"

	"x.hh.ru/checkErr"

	tb "gopkg.in/tucnak/telebot.v2"
)

func helpAndStart(m *tb.Message, bot *tb.Bot) {
	if m.Sender.ID == teleAdminID {
		msg := `timeoutResumeUpdate=<Установить частоту обновления резюме (в минутах)>

setResume=<Сохранить выбранные резюме для обновления (перечислить номера резюме (результат команды getResume) через запятую)>

setUpdateService=<true|false остановить или запустить службу обновления резюме

startUpdate=<(true, false)  принудительно запустить обновление резюме в видимом или не видимом режиме>

startAuthentication Запустить браузер для авторизации

getResume Получить список резюме

getTimeoutResumeUpdate Получить тайм-аут

getUpdateStatus Получить состояние службы обновления резюме`
		bot.Send(m.Sender, msg)
	}
}

// Запуск бота и обработка команд от пользователя
func startBot() {
	chromeIsRunning := false
	bot, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	checkerr.СheckErr(err)
	teleAdminUser := tb.User{ID: teleAdminID}
	bot.Handle("/help", func(m *tb.Message) {
		helpAndStart(m, bot)
	})
	bot.Handle("/start", func(m *tb.Message) {
		helpAndStart(m, bot)
	})

	bot.Handle(tb.OnText, func(m *tb.Message) {
		if m.Sender.ID == teleAdminID {
			if strings.HasPrefix(strings.ToLower(m.Text), "startupdate") {
				text := strings.Split(m.Text, "=")[1]
				if !chromeIsRunning {
					chromeIsRunning = true
					if strings.ToLower(strings.TrimSpace(text)) == "false" {
						bot.Send(m.Sender, "Обновляем...")
						var cancel context.CancelFunc
						var ctx context.Context
						lock.Lock()
						for _, resume := range resumeForUpdates {
							ctx, cancel = prepareChrome(false)
							updateResume(ctx, resume)
						}
						lock.Unlock()
						cancel()
						bot.Send(m.Sender, "Готово")
					} else if strings.ToLower(strings.TrimSpace(text)) == "true" {
						bot.Send(m.Sender, "Обновляем...")
						var cancel context.CancelFunc
						var ctx context.Context
						lock.Lock()
						for _, resume := range resumeForUpdates {
							ctx, cancel = prepareChrome(true)
							updateResume(ctx, resume)
						}
						lock.Unlock()
						cancel()
						bot.Send(m.Sender, "Готово")
					} else {
						bot.Send(m.Sender, "Неверная команда")
					}
					chromeIsRunning = false
				} else {
					bot.Send(m.Sender, "Процедура уже запущена")
				}
			} else if strings.HasPrefix(strings.ToLower(m.Text), "setupdateservice") {
				setUpdateService(m, bot)
			} else if strings.HasPrefix(strings.ToLower(m.Text), "getupdatestatus") {
				lock.Lock()
				if working {
					bot.Send(m.Sender, "Служба обновления резюме работает")
				} else {
					bot.Send(m.Sender, "Служба обновления резюме не работает")
				}
				lock.Unlock()
			} else if strings.HasPrefix(strings.ToLower(m.Text), "startauthentication") {
				if !chromeIsRunning {
					chromeIsRunning = true
					ctx, cancel := prepareChrome(true)
					firstRunChrome(ctx, cancel)
					chromeIsRunning = false
				} else {
					bot.Send(m.Sender, "Процедура уже запущена")
				}
			} else if strings.HasPrefix(strings.ToLower(m.Text), "timeoutresumeupdate") {
				saveTimeoutResumeUpdate(m, bot)
			} else if strings.HasPrefix(strings.ToLower(m.Text), "setresume") {
				saveResume(m, bot)
			} else if strings.HasPrefix(strings.ToLower(m.Text), "getresume") {
				if !chromeIsRunning {
					chromeIsRunning = true
					bot.Send(m.Sender, "Получаем данные, ожидайте...")
					ctx, cancel := prepareChrome(false)
					resumeList := getResumeList(ctx, cancel)
					msg := ""
					for i, resume := range resumeList {
						msg += strconv.Itoa(i+1) + " - " + resume + "\n"
					}
					bot.Send(m.Sender, msg)
					chromeIsRunning = false
				} else {
					bot.Send(m.Sender, "Процедура уже запущена")
				}
			} else if strings.HasPrefix(strings.ToLower(m.Text), "gettimeoutresumeupdate") {
				lock.Lock()
				bot.Send(m.Sender, strconv.Itoa(timeoutResumeUpdate))
				lock.Unlock()
			} else {
				bot.Send(m.Sender, "Неверная команда")
			}
		}
	})
	bot.Send(&teleAdminUser, "Запуск - "+time.Now().Format(time.ANSIC))
	bot.Start()
}

func saveResume(m *tb.Message, bot *tb.Bot) {
	text := strings.Split(m.Text, "=")
	if len(text) == 2 {
		lock.Lock()
		if len(strings.TrimSpace(text[1])) > 0 {
			resumeTmpList := strings.Split(strings.TrimSpace(text[1]), ",")
			for _, v := range resumeTmpList {
				var err error
				_, err = strconv.Atoi(v)
				if err != nil {
					bot.Send(m.Sender, "Минимум в одном из параметров не числовое значение")
					lock.Unlock()
					return
				}
			}
			resumeForUpdates = []string{}
			for _, v := range resumeTmpList {
				resumeForUpdates = append(resumeForUpdates, strings.TrimSpace(v))
			}
			saveCfg()
			lock.Unlock()
			bot.Send(m.Sender, "Список резюме успешно сохранён")
			return
		} else {
			bot.Send(m.Sender, "Не задано не одного значения")
			lock.Unlock()
			return
		}
	} else {
		bot.Send(m.Sender, "Неверная команда")
	}
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
		bot.Send(m.Sender, "Неверная команда")
	}
}

func setUpdateService(m *tb.Message, bot *tb.Bot) {
	message := strings.Split(m.Text, "=")
	if len(message) == 2 && (strings.ToLower(message[1]) == "true" || strings.ToLower(message[1]) == "false") {
		lock.Lock()
		switch strings.ToLower(message[1]) {
		case "true":
			working = true
			bot.Send(m.Sender, "Служба обновления резюме теперь работает")
		case "false":
			working = false
			bot.Send(m.Sender, "Служба обновления резюме теперь не работает")
		default:
			bot.Send(m.Sender, "Неверная команда")
		}
		lock.Unlock()
	} else {
		bot.Send(m.Sender, "Неверная команда")
	}
}
