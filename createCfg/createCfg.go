// Первоначальное создание файла конфигурации.
package main

import (
	"errors"
	"log"
	"os"
	"strings"

	"x.hh.ru/crypting"
)

const (
	cfgFile = "config.cfg"
)

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func firstEncryptCfg() {
	password, isexist := os.LookupEnv("HHSalt")
	if !isexist && len(strings.TrimSpace(password)) == 0 {
		checkErr(errors.New("environment variable 'HHSalt' is not exist"))
	}
	text, err := os.ReadFile("cfg")
	if err != nil {
		checkErr(errors.New("error open and read 'cfg' file"))
	}
	var preparedText string
	for _, v := range strings.Split(string(text), "\r\n") {
		preparedText += v + "\n"
	}
	preparedText = strings.TrimRight(preparedText, "\n")
	crypting.EncryptFile(cfgFile, []byte(preparedText), password)
}

func main() {
	firstEncryptCfg()
}
