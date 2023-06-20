module x.hh.ru

go 1.16

require (
	github.com/chromedp/cdproto v0.0.0-20230601223857-b9862e644d94
	github.com/chromedp/chromedp v0.9.1
	github.com/gobwas/ws v1.2.1 // indirect
	golang.org/x/sys v0.8.0 // indirect
	gopkg.in/tucnak/telebot.v2 v2.3.5
	x.hh.ru/checkErr v0.0.0-00010101000000-000000000000
	x.hh.ru/crypting v0.0.0-00010101000000-000000000000
	x.hh.ru/logs v0.0.0-00010101000000-000000000000
)

replace x.hh.ru/logs => ./logs

replace x.hh.ru/checkErr => ./checkErr

replace x.hh.ru/crypting => ./crypting
