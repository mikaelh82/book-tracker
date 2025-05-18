package main

import (
	"fmt"
	"net/http"
)

func healthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "PONG")
}

func main() {

	http.HandleFunc("/ping", healthHandler)
	http.ListenAndServe(":8080", nil)

}
