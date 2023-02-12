package main

import (
	"fmt"
	"net/http"
)

const portNumber = ":8085"

func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about", About)

	fmt.Printf("Starting application on port %s", portNumber)

	http.ListenAndServe(portNumber, nil)
}
