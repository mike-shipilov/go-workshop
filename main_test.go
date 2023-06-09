package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleIndex(t *testing.T) {
	// setup
	s := server{}
	s.routes()

	// given
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	// when
	s.router.ServeHTTP(w, r)

	//then
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("Expected status code is 200, but got: %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	wantResBody := `{"message":"Hello server"}`
	gotResBody := string(b)
	if gotResBody != wantResBody {
		t.Fatalf("Want: %v, got: %v", wantResBody, gotResBody)
	}
}
