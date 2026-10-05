package main

import (
	"log"
	"os"

	fixduplicateglobaloutline "Qingyu_backend/internal/app/tools/fixduplicateglobaloutline"
)

func main() {
	if err := fixduplicateglobaloutline.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
