package main

import (
	"log"
)

func main() {
	log.Println("Running GSL server on :8080")

	// Now we need to explicitly add our data store, but this is the
	// perfect place to do some configuration!
	dataStore := NewGSLDataStore("./champion.txt")
	err := runServer(":8080", dataStore)

	if err != nil {
		log.Fatal(err)
	}
}
