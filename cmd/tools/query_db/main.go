package main

import (
	"log"
	"os"

	querydb "Qingyu_backend/internal/app/tools/querydb"
)

func main() {
	if err := querydb.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
