package main

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

func TestGSLCurrentChampionIsTY(t *testing.T) {
	// We still have to do this.  I hate my life.
	contents, err := ioutil.ReadFile("./champion.txt")

	if err != nil {
		t.Fatal("Failed to read file:", err)
	}

	expectedWorldChampion := string(contents)

	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	dataStore := NewGSLDataStore("champion.txt")

	handler := gslCurrentChampionHandler(dataStore)

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
	req := httptest.NewRequest("GET", "/champion", nil)
	res := httptest.NewRecorder()

	// We can do this because we've pulled out the data store from inside
	// the handler, so that's nice at least...
	dataStore := NewGSLDataStore("i-dont-exist.txt")

	handler := gslCurrentChampionHandler(dataStore)

	handler(res, req)

	gotCode := res.Code

	if 500 != gotCode {
		t.Errorf("Expected code 500 but got %d", gotCode)
	}
}
