package main

import (
	"log"
	"net/http"
)

// CurrentChampionGetter can get the current champion somehow
//
// This is a declaration of intent.  Whoever asks for this interface is going
// to ask for "something" that can get the current champion.  We don't care how.
// We don't need to know how.  We don't want to know.  All we know is that
// it can get the current champion, and we cannot do anything else.
type CurrentChampionGetter interface {
	GetCurrentChampion() (string, error)
}

// Creates a handler that writes the current champion to the client.
//
// Now our handler is saying something very powerful.  It's saying
// "I need something that can get the current champion in order to
// do my job."  This is much more descriptive than before!
func gslCurrentChampionHandler(currentChampionGetter CurrentChampionGetter) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		champion, err := currentChampionGetter.GetCurrentChampion()

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
// Our server is
func runServer(address string, currentChampionGetter CurrentChampionGetter) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/champion", gslCurrentChampionHandler(currentChampionGetter))

	return http.ListenAndServe(address, mux)
}
