package internals

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
)

type Server struct {
	http.ServeMux
}

func (s *Server) ReadFiles(dir string) {
	err := filepath.WalkDir(dir, s.parseFile)

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func (s *Server) parseFile(path string, d fs.DirEntry, err error) error {
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

		req := fmt.Sprintf("%s %s", endpoint.Method, endpoint.ParseUri())

		def, err := endpoint.GetDefaultResponse()
		if err != nil {
			fmt.Printf("no default response available for %s\n", req)
		}

		fmt.Println(req)

		s.HandleFunc(req, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(def.Status)

			data, err := json.Marshal(def.Body)
			if err != nil {
				fmt.Println("error converting body to json")
			}

			w.Write(data)
		})
	}

	return nil
}
