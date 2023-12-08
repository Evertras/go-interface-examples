package main

import (
	"os"
	"log"
	"net/http"
)

func gslCurrentChampionHandler(res http.ResponseWriter, req *http.Request) {
	// We magically know it's in champion.txt and we magically know it's a plaintext file
	// that only contains the name with no line break at the end.  This is terrible.
	contents, err := os.ReadFile("./champion.txt")

	if err != nil {
		log.Println("Failed to read file:", err)
		res.WriteHeader(500)
		return
	}

	res.Write(contents)
}

func runServer(address string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/champion", gslCurrentChampionHandler)

	return http.ListenAndServe(address, mux)
}
