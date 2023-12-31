//Настройка логов. Назначение вывода, формат вывода
package logs

import (
	"log"
	"os"
)

func PrepareLogger() {
	logFile, err := os.OpenFile("logs.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		os.Exit(2)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
