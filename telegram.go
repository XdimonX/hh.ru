package main

import (
	// "fmt"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func tst() {
	b, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	checkErr(err)
	g := tb.User{ID: 385060683}
	var (
		menu    = &tb.ReplyMarkup{ResizeReplyKeyboard: true}
		btnHelp = menu.Text("ℹ Help")
	)
	menu.Reply(
		menu.Row(btnHelp),
	)
	b.Handle(&btnHelp, func(m *tb.Message) {
		b.Send(m.Chat, "Тута будет помощь.")
	})
	b.Handle("/help", func(m *tb.Message) {
		b.Send(m.Chat, "test", menu)
	})
	b.Send(&g, "Запуск - "+time.Now().Format(time.ANSIC))
	b.Start()
}
