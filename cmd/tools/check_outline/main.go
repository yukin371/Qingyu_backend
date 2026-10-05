package main

import (
	"log"

	checkoutline "Qingyu_backend/internal/app/tools/checkoutline"
)

func main() {
	if err := checkoutline.Run(); err != nil {
		log.Fatal(err)
	}
}
