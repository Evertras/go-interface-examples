package main

import (
	"errors"
	"net/http/httptest"
	"testing"
)

// A simple mock that lets us specify who the current champion is, or even
// what error to return to test our error handling!
type mockCurrentChampionGetter struct {
	current      string
	pendingError error
}

// We match the interface that our handler requires, so we can pass it in
func (g *mockCurrentChampionGetter) GetCurrentChampion() (string, error) {
	if g.pendingError != nil {
		return "", g.pendingError
	}

	return g.current, nil
}

func TestGSLCurrentChampionIsTY(t *testing.T) {
	expectedWorldChampion := "TY"

	// No more files!  We can specify the exact scenario we want.
	championGetter := &mockCurrentChampionGetter{
		current: expectedWorldChampion,
	}

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	handler := gslCurrentChampionHandler(championGetter)

	handler(res, req)

	gotWorldChampion := string(res.Body.Bytes())
	gotCode := res.Code

	if 200 != gotCode {
		t.Errorf("Expected code 200 but got %d", gotCode)
	}

	if expectedWorldChampion != gotWorldChampion {
		t.Errorf("Expected world champion to be %q but got %q", expectedWorldChampion, gotWorldChampion)
	}
}

// On the bright side, we can add this test now!
func TestGSLCurrentChampionReturns500WhenFileDoesntExist(t *testing.T) {
	// Now we can trigger error handling in a very explicit, simple way,
	// without relying on weird underlying implementation... much less
	// fragile on all fronts!
	championGetter := &mockCurrentChampionGetter{
		pendingError: errors.New("something bad happened"),
	}

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	handler := gslCurrentChampionHandler(championGetter)

	handler(res, req)

	gotCode := res.Code

	if 500 != gotCode {
		t.Errorf("Expected code 500 but got %d", gotCode)
	}
}
