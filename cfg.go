//Модуль работы с конфигурацией. Получение и парсинг конфигурации и сохранение изменений.

package main

import (
	"errors"
	"os"
	"strconv"
	"strings"

	"x.hh.ru/checkErr"
	"x.hh.ru/crypting"
)

func parseCfg() {
	var isexist bool
	password, isexist = os.LookupEnv("HHSalt")
	if !isexist && len(strings.TrimSpace(password)) == 0 {
		checkerr.СheckErr(errors.New("environment variable 'HHSalt' is not exist"))
	}

	decryptByte := crypting.DecryptFile(cfgFile, password)
	decryptText := string(decryptByte)
	decryptStrArr := strings.Split(decryptText, "\n")

	for _, v := range decryptStrArr {
		str := strings.Split(v, "&&||")
		param := strings.ToLower(str[0])
		switch param {
		case "token":
			token = str[1]
		case "resumeforupdates":
			str := strings.Split(str[1], "|")
			for _, v := range str {
				if v != "" {
					resumeForUpdates = append(resumeForUpdates, v)
				}
			}
		case "timeoutresumeupdates":
			var err error
			timeoutResumeUpdate, err = strconv.Atoi(str[1])
			checkerr.СheckErr(err)
		case "teleadminid":
			var err error
			teleAdminID, err = strconv.Atoi(str[1])
			checkerr.СheckErr(err)
		default:

		}
	}
}

func saveCfg() {
	resumeStr := ""
	for i, v := range resumeForUpdates {
		if i == 0 {
			resumeStr = v + "|"
		} else {
			resumeStr += v + "|"
		}
	}
	outCfg := "token&&||" + token +
		"\nresumeForUpdates&&||" + resumeStr +
		"\ntimeoutResumeUpdates&&||" + strconv.Itoa(timeoutResumeUpdate) +
		"\nteleAdminID&&||" + strconv.Itoa(teleAdminID)
	crypting.EncryptFile(cfgFile, []byte(outCfg), password)

}
