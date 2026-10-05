package main

import (
	"log"
	"os"

	cleanduplicates "Qingyu_backend/internal/app/tools/cleanduplicates"
)

func main() {
	if err := cleanduplicates.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
