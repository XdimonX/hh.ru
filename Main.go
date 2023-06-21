package main

import (
	"sync"
	"x.hh.ru/logs"
)

const (
	cfgFile                             = "config.cfg"
	timeOutContextUpdateResumeInSeconds = 300
	teleAdminID                         = 385060683
)

var (
	token               = ""
	loginHHru           = ""
	passwordHHru        = ""
	passwordTeleBot     = ""
	resumeForUpdates    []string
	timeoutResumeUpdate = 0
	lock                = &sync.Mutex{}
	working             = true
	password            = ""
)

func main() {
	logs.PrepareLogger()
	parseCfg()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go startBot()
	go goUpdateMonitor(false)
	wg.Wait()
}
