package main

import (
	"log"
	"os"

	migrateapp "Qingyu_backend/internal/app/migrate"
)

func main() {
	if err := migrateapp.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
