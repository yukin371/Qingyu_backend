package main

import (
	"log"
	"os"

	seedrealtimeranking "Qingyu_backend/internal/app/tools/seedrealtimeranking"
)

func main() {
	if err := seedrealtimeranking.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
