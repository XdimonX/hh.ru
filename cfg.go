package main

import (
	"strconv"
	"strings"
)

func parseCfg() {
	decryptByte := decryptFile(cfgFile, password)
	decryptText := string(decryptByte)
	decryptStrArr := strings.Split(decryptText, "\n")
	for _, v := range decryptStrArr {
		str := strings.Split(v, "&&||")
		param := strings.ToLower(str[0])
		if param == "token" {
			token = str[1]
		} else if param == "loginhhru" {
			loginHHru = str[1]
		} else if param == "passwordhhru" {
			passwordHHru = str[1]
		} else if param == "passwordtelebot" {
			passwordTeleBot = str[1]
		} else if param == "resumeforupdates" {
			str := strings.Split(str[1], "|")
			for _, v := range str {
				if v != "" {
					resumeForUpdates = append(resumeForUpdates, v)
				}
			}
		} else if param == "timeoutresumeupdates" {
			var err error
			timeoutResumeUpdate, err = strconv.Atoi(str[1])
			checkErr(err)
			println(timeoutResumeUpdate)
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
	outCfg := "Token&&||" + token + "\nloginHHru&&||" + loginHHru + "\npasswordHHru&&||" + passwordHHru +
		"\npasswordTeleBot&&||" + passwordTeleBot + "\nresumeForUpdates&&||" + resumeStr + "\ntimeoutResumeUpdates&&||" + strconv.Itoa(timeoutResumeUpdate)
	// _ = outCfg
	encryptFile(cfgFile, []byte(outCfg), password)

}
