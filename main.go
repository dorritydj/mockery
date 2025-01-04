package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type Endpoint struct {
	Id       string `json:"id"`
	Uri      string `json:"uri"`
	Method   string `json:"method"`
	Variants map[string]Variant
}

type Endpoints map[string]Endpoint

type Variant struct {
	Status  int  `json:"status"`
	Default bool `json:"default"`
	Body    any  `json:"body"`
}

func (e Endpoint) getDefaultResponse() Variant {
	var found Variant
	for _, v := range e.Variants {
		if v.Default {
			found = v
		}
	}

	return found
}

func (e Endpoint) parseUri() string {
	var parts []string
	res := strings.Fields(strings.ReplaceAll(e.Uri, "/", " "))

	for _, partition := range res {
		colon := string(partition[0])
		rest := string(partition[1:])

		if strings.Compare(colon, ":") == 0 {
			parts = append(parts, fmt.Sprintf("/{%s}", rest))
			continue
		}

		parts = append(parts, fmt.Sprintf("/%s", partition))
	}

	return strings.Join(parts, "")
}

func main() {
	// todo: config via env
	dir := "./routes"

	mux := http.NewServeMux()

	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			var endpoint Endpoint

			file, err := os.ReadFile(path)
			if err != nil {
				return err
			}

			err = json.Unmarshal(file, &endpoint)
			if err != nil {
				return err
			}

			// fmt.Printf("%#v\n", endpoint)

			req := fmt.Sprintf("%s %s", endpoint.Method, endpoint.parseUri())
			fmt.Println(req)
			def := endpoint.getDefaultResponse()

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
