package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"mockery.dorrity.dj/internals"
)

func main() {
	dir := "./routes"

	mux := http.NewServeMux()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			var endpoint internals.Endpoint

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = json.Unmarshal(file, &endpoint)
			if err != nil {
				return err
			}

			// fmt.Printf("%#v\n", endpoint)

			req := fmt.Sprintf("%s %s", endpoint.Method, endpoint.ParseUri())
			fmt.Println(req)
			def := endpoint.GetDefaultResponse()

			mux.HandleFunc(req, func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(def.Status)

				data, err := json.Marshal(def.Body)
				if err != nil {
					fmt.Println("error converting body to json")
				}

				w.Write(data)
			})
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	mux.HandleFunc("GET /{slug}", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
	})

	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("error starting server")
	}
}
