package main

import (
	"fmt"
	"net/http"

	"mockery.dorrity.dj/internals"
)

func main() {
	app := &internals.Server{}
	app.ReadFiles("./routes")

	err := http.ListenAndServe(":8080", app)
	if err != nil {
		fmt.Println("error starting server")
	}
}
