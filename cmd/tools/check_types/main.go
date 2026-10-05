package main

import (
	"log"
	"os"

	checktypes "Qingyu_backend/internal/app/tools/checktypes"
)

func main() {
	if err := checktypes.Run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
