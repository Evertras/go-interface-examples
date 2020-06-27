package leaderboard

import (
	"context"
	"fmt"

	"github.com/Evertras/go-interface-examples/local-interfaces/db"
)

// TopUserGetter gets top users from somewhere
//
// Notice that we're tying ourselves to the database package here,
// but we'll live with it.  The important part here is that this
// is just a tiny subsection of the full database functionality,
// because this is all a leaderboard needs to care about.
//
// Think about how it feels to read this and what it tells you about
// what the code in this package does compared to what it would
// be like seeing the full Database functionality at Leaderboard's
// disposal.  When we declare our interface locally like this,
// it makes the intent MUCH clearer to the reader/maintainer!
type TopUserGetter interface {
	GetTopUsers(ctx context.Context, count int) ([]*db.User, error)
}

// TopScoreNotifier notifies users they have a high score
//
// As above, focus on how this feels to read in terms of understanding
// the code in this package.
type TopScoreNotifier interface {
	NotifyTopScore(ctx context.Context, id string, score int) error
}

// Leaderboard knows how to interact with top users
type Leaderboard struct {
	topUserGetter    TopUserGetter
	topScoreNotifier TopScoreNotifier
}

// New creates a new Leaderboard ready to do leaderboard things
//
// Accept interfaces, return implementations!
func New(topUserGetter TopUserGetter, topScoreNotifier TopScoreNotifier) *Leaderboard {
	return &Leaderboard{
		topUserGetter:    topUserGetter,
		topScoreNotifier: topScoreNotifier,
	}
}

// NotifyTopPlayers will send a notification to the top X players
func (l *Leaderboard) NotifyTopPlayers(ctx context.Context, top int) error {
	users, err := l.topUserGetter.GetTopUsers(ctx, top)

	if err != nil {
		return fmt.Errorf("topUserGetter.GetTopUsers: %w", err)
	}

	for _, user := range users {
		err = l.topScoreNotifier.NotifyTopScore(ctx, user.ID, user.Score)

		// Error out early.  If this was real maybe we want to at least try
		// to notify everyone, but that's not important now.
		if err != nil {
			return fmt.Errorf("topScoreNotifier.NotifyTopScore: %w", err)
		}
	}

	return nil
}
