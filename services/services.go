package services

import "log"

func checkErr(e error) {
	if e != nil {
		log.Fatal(e)
	}
}
