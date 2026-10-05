package main

import (
	"log"
	"os"

	benchmarkapp "Qingyu_backend/internal/app/benchmark"
)

func main() {
	if err := benchmarkapp.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
