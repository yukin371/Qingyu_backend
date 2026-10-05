package main

import (
	"log"
	"os"

	testhierarchy "Qingyu_backend/internal/app/tools/testhierarchy"
)

func main() {
	if err := testhierarchy.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
