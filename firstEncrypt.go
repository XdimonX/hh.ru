package main

import (
	"fmt"
	"io/ioutil"
)

func firstEncryptCfg() {
	text, err := ioutil.ReadFile("cfg")
	if err != nil {
		panic("aaaaa")
	}
	fmt.Println(string(text))
	encryptFile(cfgFile, text, password)
}

// func main() {
// 	firstEncryptCfg()
// }
