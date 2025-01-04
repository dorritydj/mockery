package main

import (
	"fmt"
	"net/http"
	"path/filepath"

	"mockery.dorrity.dj/internals"
)

func main() {
	dir := "./routes"

	app := &internals.Server{
		Mux: http.NewServeMux(),
	}

	err := filepath.WalkDir(dir, app.ParseFile)

	if err != nil {
		fmt.Println("Error:", err)
	}

	app.Mux.HandleFunc("GET /{slug}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
	})

	err = http.ListenAndServe(":8080", app.Mux)
	if err != nil {
		fmt.Println("error starting server")
	}
}
