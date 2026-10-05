package main

import (
	"log"

	checkdata "Qingyu_backend/internal/app/tools/checkdata"
)

func main() {
	if err := checkdata.Run(); err != nil {
		log.Fatal(err)
	}
}
