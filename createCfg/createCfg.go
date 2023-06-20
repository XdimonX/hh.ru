// Первоначальное создание файла конфигурации.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
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
	text, err := ioutil.ReadFile("cfg")
	if err != nil {
		checkErr(errors.New("error open and read 'cfg' file"))
	}
	fmt.Println(string(text))
	crypting.EncryptFile(cfgFile, text, password)
}

func main() {
	firstEncryptCfg()
}
