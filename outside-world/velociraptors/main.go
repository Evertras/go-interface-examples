package main

import (
	"log"
)

func main() {
	log.Println("Running GSL server on :8080")

	err := runServer(":8080")

	if err != nil {
		log.Fatal(err)
	}
}
