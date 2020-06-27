package handlers

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

// UserDataStore can access and modify user data
//
// This is a little more broad and we're assuming the same type
// will have all these methods, but that's a totally reasonable
// compromise compared to creating a separate interface for every
// single potential method.  The important thing is that it's in
// the local package here.
type UserDataStore interface {
	GetUserScore(ctx context.Context, id string) (int, error)
	DeleteUser(ctx context.Context, id string) error
}

// GetUserScoreHandler creates an HTTP handler that can get a user's score
func GetUserScoreHandler(userDataStore UserDataStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.Header.Get("x-user-id")

		score, err := userDataStore.GetUserScore(req.Context(), id)

		if err != nil {
			fmt.Println("userDataStore.GetUserScore: ", err)
			res.WriteHeader(500)
			return
		}

		res.Write([]byte(fmt.Sprintf("%d", score)))
	}
}

// DeleteUserHandler creates an HTTP handler that deletes a user from the store
func DeleteUserHandler(userDataStore UserDataStore) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Totally trust the client, this is fine (it's not, don't do this)
		body, err := ioutil.ReadAll(req.Body)

		if err != nil {
			fmt.Println("ioutil.ReadAll(req.Body): ", err)
			res.WriteHeader(500)
			return
		}

		id := string(body)

		err = userDataStore.DeleteUser(req.Context(), id)

		if err != nil {
			fmt.Println("userDataStore.GetUserScore: ", err)
			res.WriteHeader(500)
			return
		}
	}
}
