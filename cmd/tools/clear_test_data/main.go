package main

import (
	"log"
	"os"

	cleartestdata "Qingyu_backend/internal/app/tools/cleartestdata"
)

func main() {
	if err := cleartestdata.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
