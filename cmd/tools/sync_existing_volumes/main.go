package main

import (
	"log"
	"os"

	syncexistingvolumes "Qingyu_backend/internal/app/tools/syncexistingvolumes"
)

func main() {
	if err := syncexistingvolumes.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
