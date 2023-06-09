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
	r := chi.NewRouter()
	r.Get("/", handleIndex)
	r.Put("/", handleIndex)
	http.ListenAndServe(":8080", r)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	var t thing
	if r.Method == http.MethodGet {
		b, err := os.ReadFile(thingTXT)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(b)
		return
	}
	if r.Method == http.MethodPut {
		err := json.NewDecoder(r.Body).Decode(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		s, err := json.Marshal(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		b := []byte(s)
		os.WriteFile(thingTXT, b, 0644)
		log.Printf("Got thing: %#v", t)
		return
	}
	http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
}
