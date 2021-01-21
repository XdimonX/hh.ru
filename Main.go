package main

import (
	"fmt"
	"log"
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
	working             = true
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
	wg.Add(1)
	go startBot()
	go goUpdateMonitor(false)
	wg.Wait()
}
