package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestGSLCurrentChampionIsTY(t *testing.T) {
	// I hate everything about this.  Writing this has caused my keyboard
	// to rebel in anger.  Do not use this.  Do not even think about it
	// for too long or adverse health effects may arise.
	contents, err := ioutil.ReadFile("./champion.txt")

	if err != nil {
		t.Fatal("Failed to read file:", err)
	}

	expectedWorldChampion := string(contents)

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	gslCurrentChampionHandler(res, req)

	gotWorldChampion := string(res.Body.Bytes())
	gotCode := res.Code

	if 200 != gotCode {
		t.Errorf("Expected code 200 but got %d", gotCode)
	}

	if expectedWorldChampion != gotWorldChampion {
		t.Errorf("Expected world champion to be %q but got %q", expectedWorldChampion, gotWorldChampion)
	}
}
