package main

import (
	"net/http/httptest"
	"testing"
)

func TestGSLCurrentChampionIsTY(t *testing.T) {
	expectedWorldChampion := "TY"

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	gslCurrentChampionHandler(res, req)

	gotWorldChampion := string(res.Body.Bytes())

	if expectedWorldChampion != gotWorldChampion {
		t.Errorf("Expected world champion to be %q but got %q", expectedWorldChampion, gotWorldChampion)
	}
}
