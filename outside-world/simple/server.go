package main

import (
	"net/http"
)

func gslCurrentChampionHandler(res http.ResponseWriter, req *http.Request) {
	// This is correct as of 2020-07-18
	res.Write([]byte("TY"))
}

// Runs a server that lets us see GSL info
func runServer(address string) error {
	mux := http.NewServeMux()

	// This is our only handler for now, for demonstration's sake.
	mux.HandleFunc("/champion", gslCurrentChampionHandler)

	return http.ListenAndServe(address, mux)
}
