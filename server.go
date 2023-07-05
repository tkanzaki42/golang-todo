package main

import (
	"handler"

	"net/http"
)

func main() {
	http.HandleFunc("/api/todo", handler.HTTPHandler)
	http.ListenAndServe(":4000", nil)
}
