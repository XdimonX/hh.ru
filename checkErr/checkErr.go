package checkerr

import (
	"fmt"
	"log"
	"os"
)

func СheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
		os.Exit(2)
	}
}