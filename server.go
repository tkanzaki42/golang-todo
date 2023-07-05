package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/todos/{id}", GetTodo).Methods(http.MethodGet)
	router.HandleFunc("/api/todos", CreateTodo).Methods(http.MethodPost)
	router.HandleFunc("/api/todos/{id}", UpdateTodo).Methods(http.MethodPut)
	router.HandleFunc("/api/todos/{id}", DeleteTodo).Methods(http.MethodDelete)
	router.HandleFunc("/api/todos", GetTodos).Methods(http.MethodGet)
	router.HandleFunc("/api/todos", DeleteTodos).Methods(http.MethodDelete)
	router.HandleFunc("/", HTTPHandler)

	http.ListenAndServe(":4000", router)
}
