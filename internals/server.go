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
	endpoints Endpoints
}

func (s *Server) ParseConfigFiles(dir string) {
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			endpoint, err := getEndpointFromFile(path)
			if err != nil {
				return err
			}
			s.endpoints = append(s.endpoints, endpoint)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	s.addEndpointHandlers()
}

func (s *Server) addEndpointHandlers() {
	for _, e := range s.endpoints {
		req := fmt.Sprintf("%s %s", e.Method, e.ParseUri())

		def, err := e.GetDefaultResponse()
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
}

func getEndpointFromFile(path string) (Endpoint, error) {
	var endpoint Endpoint

	file, err := os.ReadFile(path)
	if err != nil {
		return endpoint, err
	}

	err = json.Unmarshal(file, &endpoint)
	if err != nil {
		return endpoint, err
	}

	return endpoint, nil
}
