package main

import (
	"log"
	"os"

	debugoutlinequery "Qingyu_backend/internal/app/tools/debugoutlinequery"
)

func main() {
	if err := debugoutlinequery.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
