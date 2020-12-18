package main

import (
	// "fmt"
	"os"
	"sync"
)

const (
	password = "qZ}~zo~f)7rUy<4)\"?/r?OB4№7iKl)Xpo?ypx!M>kls,}xIjF}"
	cfgFile  = "config.cfg"
)

var (
	token               = ""
	loginHHru           = ""
	passwordHHru        = ""
	passwordTeleBot     = ""
	resumeForUpdates    []string
	timeoutResumeUpdate = 0
	lock                = &sync.Mutex{}
)

func checkErr(err error) {
	if err != nil {
		os.Exit(2)
	}
}

func main() {
	parseCfg()
	wg := sync.WaitGroup{}
	wg.Add(1)
	go startBot()
	wg.Wait()
}
