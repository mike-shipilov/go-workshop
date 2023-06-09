package main

import (
	"encoding/json"
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
	s := server{}
	s.routes()
	http.ListenAndServe(":8080", s.router)
}

type server struct {
	router *chi.Mux
}

func (s *server) routes() {
	s.router = chi.NewRouter()
	s.router.Get("/", s.handleGet)
	s.router.Put("/", s.handlePut)
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	b, err := os.ReadFile(thingTXT)
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
	str, err := json.Marshal(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b := []byte(str)
	os.WriteFile(thingTXT, b, 0644)
	log.Printf("Got thing: %#v", t)
}
