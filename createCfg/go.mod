module createcfg

go 1.19

replace x.hh.ru/crypting => ./crypting

replace x.hh.ru/checkErr => ./checkErr

require x.hh.ru/crypting v0.0.0-00010101000000-000000000000

require x.hh.ru/checkErr v0.0.0-00010101000000-000000000000 // indirect
