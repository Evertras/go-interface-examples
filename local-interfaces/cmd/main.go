package main

import (
	"context"

	"github.com/Evertras/go-interface-examples/local-interfaces/db"
	"github.com/Evertras/go-interface-examples/local-interfaces/handlers"
	"github.com/Evertras/go-interface-examples/local-interfaces/leaderboard"
	"github.com/Evertras/go-interface-examples/local-interfaces/notifications"
)

func main() {
	database := db.New()
	notifier := notifications.New()

	// Our database and notifier match the local interfaces in leaderboard,
	// so we can use them fine
	leaderboard := leaderboard.New(database, notifier)

	leaderboard.NotifyTopPlayers(context.Background(), 3)

	// Similarly, our handlers expect a certain interface which is also fulfilled
	// by our database, so we could create an HTTP server here
	handlers.DeleteUserHandler(database)

	// Don't actually do anything useful, we're just showing how the interfaces
	// properly match!
}
