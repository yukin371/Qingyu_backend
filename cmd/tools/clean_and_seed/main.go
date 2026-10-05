package main

import (
	"log"
	"os"

	cleanandseed "Qingyu_backend/internal/app/tools/cleanandseed"
)

func main() {
	if err := cleanandseed.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
