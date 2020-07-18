package main

import (
	"log"
)

func main() {
	log.Println("Running GSL server on :8080")

	// This is the same as before, because dataStore matches the CurrentChampionGetter interface
	dataStore := NewGSLDataStore("./champion.txt")
	err := runServer(":8080", dataStore)

	if err != nil {
		log.Fatal(err)
	}
}
