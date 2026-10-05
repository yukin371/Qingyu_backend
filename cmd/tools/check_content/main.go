package main

import (
	"log"

	checkcontent "Qingyu_backend/internal/app/tools/checkcontent"
)

func main() {
	if err := checkcontent.Run(); err != nil {
		log.Fatal(err)
	}
}
