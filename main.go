package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

type thing struct {
	Message string `json:"message"`
}

const thingTXT = "thing.txt"

func main() {
	s := server{
		db: &dbFile{},
	}
	s.routes()
	http.ListenAndServe(":8080", s.router)
}

type server struct {
	db     db
	router *chi.Mux
}

type db interface {
	getThing() (*thing, error)
	putThing(t *thing) error
}

type dbFile struct{}

func (d *dbFile) getThing() (*thing, error) {
	b, err := os.ReadFile(thingTXT)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}
	var t thing
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling json: %w", err)
	}
	return &t, nil
}

func (d *dbFile) putThing(t *thing) error {
	b, err := json.Marshal(t)
	if err != nil {
		return fmt.Errorf("marshalling json: %w", err)
	}
	err = os.WriteFile(thingTXT, b, 0644)
	if err != nil {
		return fmt.Errorf("writing file: %w", err)
	}
	return nil
}

func (s *server) routes() {
	s.router = chi.NewRouter()
	s.router.Get("/", s.handleGet)
	s.router.Put("/", s.handlePut)
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	t, err := s.db.getThing()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.Marshal(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(b)
}

func (s *server) handlePut(w http.ResponseWriter, r *http.Request) {
	var t thing
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Printf("Got thing: %#v", t)
	err = s.db.putThing(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write([]byte("OK"))
}
