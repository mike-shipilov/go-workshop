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

### db file

## part 2

### test endpoint

### router

### server

### dbFile

### thing interface

### test thing
