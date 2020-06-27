package handlers

import (
	"bytes"
	"context"
	"fmt"
	"net/http/httptest"
	"testing"
)

type mockUserDataStore struct {
	pendingError error
	pendingScore int

	deletedUsers []string
}

func (m *mockUserDataStore) GetUserScore(ctx context.Context, id string) (int, error) {
	return m.pendingScore, m.pendingError
}

func (m *mockUserDataStore) DeleteUser(ctx context.Context, id string) error {
	if m.pendingError != nil {
		return m.pendingError
	}

	m.deletedUsers = append(m.deletedUsers, id)

	return nil
}

func TestGetUserScoreHandlerReturnsScore(t *testing.T) {
	req := httptest.NewRequest("GET", "/idk", nil)
	res := httptest.NewRecorder()

	userDataStore := &mockUserDataStore{
		pendingScore: 3,
	}

	handler := GetUserScoreHandler(userDataStore)

	handler(res, req)

	resultStr := string(res.Body.Bytes())
	expected := fmt.Sprintf("%d", userDataStore.pendingScore)

	if res.Code != 200 {
		t.Errorf("Expected HTTP response 200 but got %d", res.Code)
	}

	if resultStr != expected {
		t.Errorf("Expected body to contain value %q but got %q", expected, resultStr)
	}
}

func TestDeleteUserDeletesUserIDFromBody(t *testing.T) {
	id := "fakeusersomething"
	req := httptest.NewRequest("DELETE", "/user/idk", bytes.NewBufferString(id))
	res := httptest.NewRecorder()

	userDataStore := &mockUserDataStore{
		pendingScore: 3,
	}

	handler := DeleteUserHandler(userDataStore)

	handler(res, req)

	if res.Code != 200 {
		t.Errorf("Expected HTTP response 200 but got %d", res.Code)
	}

	if len(userDataStore.deletedUsers) != 1 {
		t.Fatalf("Expected %d deletions but saw %d", 1, len(userDataStore.deletedUsers))
	}

	if userDataStore.deletedUsers[0] != id {
		t.Errorf("Expected to delete id %q but deleted %q", id, userDataStore.deletedUsers[0])
	}
}
