package leaderboard

import (
	"context"
	"errors"
	"testing"

	"github.com/Evertras/go-interface-examples/local-interfaces/db"
)

// We mock our interface here; notice how simple and self-contained this
// can be when we don't have to worry about the full database functionality!
type mockTopUserGetter struct {
	pendingError error
	pendingUsers []*db.User
}

func (g *mockTopUserGetter) GetTopUsers(ctx context.Context, count int) ([]*db.User, error) {
	return g.pendingUsers, g.pendingError
}

type mockTopScoreNotifier struct {
	sentToIDs    []string
	pendingError error
}

func (g *mockTopScoreNotifier) NotifyTopScore(ctx context.Context, id string, count int) error {
	if g.pendingError != nil {
		return g.pendingError
	}

	g.sentToIDs = append(g.sentToIDs, id)

	return nil
}

func TestNotifyTopPlayersErrorsWhenGetterFails(t *testing.T) {
	// Make our getter error out
	mockGetter := &mockTopUserGetter{
		pendingError: errors.New("lolnope"),
	}

	// Notifier shouldn't matter
	mockNotifier := &mockTopScoreNotifier{}

	leaderboard := New(mockGetter, mockNotifier)

	err := leaderboard.NotifyTopPlayers(context.Background(), 5)

	if err == nil {
		t.Fatal("Should have gotten an error back, but didn't")
	}

	if len(mockNotifier.sentToIDs) != 0 {
		t.Errorf("Expected to send no notifications but sent %d", len(mockNotifier.sentToIDs))
	}
}

func TestNotifyTopPlayersSendsNotificationsToSpecifiedNumberOfPlayers(t *testing.T) {
	count := 3
	mockGetter := &mockTopUserGetter{
		pendingUsers: []*db.User{
			{
				ID:    "user-1",
				Score: 10,
			},
			{
				ID:    "user-2",
				Score: 7,
			},
			{
				ID:    "user-3",
				Score: 5,
			},
		},
	}
	mockNotifier := &mockTopScoreNotifier{}

	leaderboard := New(mockGetter, mockNotifier)

	err := leaderboard.NotifyTopPlayers(context.Background(), count)

	if err != nil {
		t.Fatal("leaderboard.NotifyTopPlayers: ", err)
	}

	if len(mockNotifier.sentToIDs) != count {
		t.Errorf("Expected %d notifications to be sent but found %d", count, len(mockNotifier.sentToIDs))
	}

	// Could test to make sure it notifies correct user IDs but that's enough for demo
}
