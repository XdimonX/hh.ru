package main

import (
	"log"
	"os"
)

func prepareLogger() {
	logFile, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(2)
	}
	log.SetOutput(logFile)
	// log.SetPrefix("")
	log.SetFlags(log.Lshortfile)
}
