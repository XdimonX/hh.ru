package checkerr

import (
	"fmt"
	"log"
)

func СheckErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
}