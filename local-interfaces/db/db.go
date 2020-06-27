package db

import (
	"context"
	"errors"
)

// Db does some IO with a database, mocked here for simplicity
type Db struct {
	// Would normally have all sorts of connection info here
}

// New returns a new Db ready to do Db things.
//
// Notice we don't provide an interface here!  We only provide the concrete
// implementation.  Accept interfaces, return implementations.
func New() *Db {
	return &Db{}
}

// User represents a user as it's stored in the database
type User struct {
	ID    string
	Score int
}

// GetUser returns a user's full information from the database
//
// Notice this returns a User, which is part of this package.  This means
// that any interface that wants to use GetUser will tie itself to this
// package.  This is often unavoidable for any non-trivial returns, but
// it's a tradeoff to be aware of.
func (d *Db) GetUser(ctx context.Context, id string) (*User, error) {
	// Just fake a result for funsies
	return &User{
		ID:    id,
		Score: 7,
	}, nil
}

// GetUserScore returns a user's score from their ID
//
// Notice this function signature doesn't contain any package-specific types,
// which means any interfaces that want to implement this do -not- need to
// tie themselves to this package.  This is great when you only need single
// fields at a time, but that won't always be the case.
func (d *Db) GetUserScore(ctx context.Context, id string) (int, error) {
	// Just fake a result for funsies
	return 7, nil
}

// CreateUser creates a user starting with a score of 0
func (d *Db) CreateUser(ctx context.Context, id string) error {
	// Just pretend we did the thing, it's fine
	return nil
}

// DeleteUser deletes a user with the given ID
func (d *Db) DeleteUser(ctx context.Context, id string) error {
	// Just pretend we did the thing, it's fine
	return nil
}

// GetTopUsers returns the top X users ranked by score
//
// Note that it makes sense here to return full user information, or at least
// some struct that contains user IDs and scores combined.  So we're tying
// interfaces to this package.  Tradeoffs.
func (d *Db) GetTopUsers(ctx context.Context, count int) ([]*User, error) {
	// I'm too lazy to build a slice, just error out if we actually try to run this
	return nil, errors.New("Not implemented, Evertras is lazy")
}

// AwardPoints gives points to all the users in the ids array
func (d *Db) AwardPoints(ctx context.Context, ids []string, score int) error {
	// Just pretend we did the thing, it's fine
	return nil
}
