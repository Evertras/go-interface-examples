package notifications

import (
	"context"
	"fmt"
)

// Notifier sends notifications to the user
//
// Actual implementation doesn't matter here, we're just
// including this to have something other than a database
// as an example.
type Notifier struct{}

// New returns a new Notifier ready to send notifications
//
// Doesn't actually do anything, but real code will do stuff here...
func New() *Notifier {
	return &Notifier{}
}

// NotifyTopScore sends a notification to a user about their top score
func (n *Notifier) NotifyTopScore(ctx context.Context, id string, score int) error {
	fmt.Printf("Sending notification to ID %q about their high score of %d\n", id, score)

	// Don't actually do anything useful...

	return nil
}

// NotifyPasswordUpdate notifies a user that their password has been updated
func (n *Notifier) NotifyPasswordUpdate(ctx context.Context, id string) error {
	fmt.Printf("Sending notification to ID %q about their password being updated\n", id)

	// Don't actually do anything useful...

	return nil
}
