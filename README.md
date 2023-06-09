# go-workshop

## part 1

### install

- Download and install Go https://go.dev/

- Install VS Code extension for Go

### print hello

```go
package main

import "fmt"

func main() {
	fmt.Print("Hello world\n")
	fmt.Println("Hello world")
	fmt.Printf("Hello %s", "world")
}
```

Run with `go run main.go`

### hello function

```go
package main

import "fmt"

func main() {
	hello()
}

func hello() {
	fmt.Println("Hello world")
}
```

### go mod init

Run `go mod init github.com/mike-shipilov/go-workshop`

Now it's possible to start the app with `go run .`

### hello server

```go
package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		b := []byte("Hello world")
		w.Write(b)
	})
	http.ListenAndServe(":8080", nil)
}
```

### get json

```go
package main

import (
	"encoding/json"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		j := map[string]any{"message": "Hello world"}
		s, err := json.Marshal(&j)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		b := []byte(s)
		w.Write(b)
	})
	http.ListenAndServe(":8080", nil)
}
```

### get data type thing

```go
package main

import (
	"encoding/json"
	"net/http"
)

type thing struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := thing{
			Message: "Hello world",
		}
		s, err := json.Marshal(&t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		b := []byte(s)
		w.Write(b)
	})
	http.ListenAndServe(":8080", nil)
}
```

### post

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type thing struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			t := thing{
				Message: "Hello world",
			}
			s, err := json.Marshal(&t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			b := []byte(s)
			w.Write(b)
			return
		}
		if r.Method == http.MethodPut {
			t := thing{}
			err := json.NewDecoder(r.Body).Decode(&t)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			log.Printf("Got thing: %#v", t)
			return
		}
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	})
	http.ListenAndServe(":8080", nil)
}
```

PUT message format to test in Thunder Client: `{"message":"Hello server"}`

### db file

```go
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type thing struct {
	Message string `json:"message"`
}

const thingTXT = "thing.txt"

func main() {
	http.HandleFunc("/", handleIndex)
	http.ListenAndServe(":8080", nil)
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
```

## part 2

### test the endpoint

```go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_handleIndex(t *testing.T) {
	// given
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	// when
	handleIndex(w, r)

	//then
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != 200 {
		t.Errorf("Expected status code is 200, but got: %d", res.StatusCode)
	}
	b, _ := io.ReadAll(res.Body)
	wantResBody := `{"message":"Hello server"}`
	gotResBody := string(b)
	if gotResBody != wantResBody {
		t.Errorf("Want: %v, got: %v", wantResBody, gotResBody)
	}
}
```

### router

```go
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
```

### server

```go
//main.go

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
```

```go
// main_test.go

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
```

### server.db

```go
//main.go

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
```

```go
// main_test.go

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

```

### thing interface

### test thing
