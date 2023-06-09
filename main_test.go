package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

type dbMock struct {
	thing *thing
}

func (d *dbMock) getThing() (*thing, error) {
	return d.thing, nil
}

func (d *dbMock) putThing(t *thing) error {
	return nil
}

func Test_server_handleGet(t *testing.T) {
	// setup
	s := server{
		db: &dbMock{thing: &thing{Message: "Hello server"}},
	}
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

func Test_server_handlePut(t *testing.T) {
	// setup
	s := &server{
		db: &dbMock{},
	}
	s.routes()

	// given
	reqBody := strings.NewReader(`{"message":"Hello server"}`)
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, "/", reqBody)

	// when
	s.router.ServeHTTP(w, r)

	// then
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Fatalf("Expected status code is 200, but got: %d", res.StatusCode)
	}
	resBody, _ := io.ReadAll(res.Body)
	wantResBody := "OK"
	gotResBody := string(resBody)
	if !reflect.DeepEqual(gotResBody, wantResBody) {
		t.Fatalf("Want: %v, got: %v", wantResBody, gotResBody)
	}
}
