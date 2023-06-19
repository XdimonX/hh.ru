package main

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const (
	cfgFile                             = "config.cfg"
	timeOutContextUpdateResumeInSeconds = 300
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
	password            = "qZ}~zo~f)7rUy<4)\"?/r?OB4№7iKl)Xpo?ypx!M>kls,}xIjF}"
)

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		os.Exit(2)
	}
}

func main() {
	prepareLogger()
	parseCfg()
	wg := sync.WaitGroup{}
	wg.Add(2)
	go startBot()
	go goUpdateMonitor(false)
	wg.Wait()
}
