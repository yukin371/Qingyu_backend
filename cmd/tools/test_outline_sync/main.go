package main

import (
	"log"
	"os"

	testoutlinesync "Qingyu_backend/internal/app/tools/testoutlinesync"
)

func main() {
	if err := testoutlinesync.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
