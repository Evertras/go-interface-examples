package main

import (
	"log"
	"net/http"
)

// Creates a handler that writes the current champion to the client.
//
// Now we're creating a handler and 'injecting' the data store that it will
// use to get the current champion.  The benefit here is that we no longer
// need to know that it's stored in a file, or what format that file is,
// or even what a 'file' is at all.  We just know we need this thing that
// can tell us who our current champion is.
func gslCurrentChampionHandler(dataStore *GSLDataStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		champion, err := dataStore.GetCurrentChampion()

		if err != nil {
			log.Println("Failed to get current champion:", err)
			res.WriteHeader(500)
			return
		}

		res.Write([]byte(champion))
	}
}

// Runs the server on the specified address with the given data store
//
// Notice we added the data store here as a parameter.  To run our server,
// we must have a data store.  We are NOT creating it here!  That's
// someone else's problem.  Always make your dependencies explicit.
func runServer(address string, dataStore *GSLDataStore) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/champion", gslCurrentChampionHandler(dataStore))

	return http.ListenAndServe(address, mux)
}
