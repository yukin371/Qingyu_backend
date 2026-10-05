package main

import (
	"log"
	"os"

	checkduplicates "Qingyu_backend/internal/app/tools/checkduplicates"
)

func main() {
	if err := checkduplicates.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
