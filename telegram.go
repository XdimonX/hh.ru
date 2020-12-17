package main

import (
	"fmt"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)

func tst() {
	b, err := tb.NewBot(tb.Settings{Token: token, Poller: &tb.LongPoller{Timeout: 10 * time.Second}})
	checkErr(err)
	_ = b
	fmt.Println(b.Me.Username)
}
