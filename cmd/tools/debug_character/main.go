package main

import (
	"log"
	"os"

	debugcharacter "Qingyu_backend/internal/app/tools/debugcharacter"
)

func main() {
	if err := debugcharacter.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
