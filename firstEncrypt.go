// Первоначальное создание файла конфигурации.
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
)

func firstEncryptCfg() {
	text, err := ioutil.ReadFile("cfg")
	if err != nil {
		checkErr(errors.New("error open and read 'cfg' file"))
	}
	fmt.Println(string(text))
	encryptFile(cfgFile, text, password)
}

// func main() {
// 	firstEncryptCfg()
// }
