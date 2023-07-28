package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var client *redis.Client

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	router := mux.NewRouter()

	router.HandleFunc("/api/todos/{id}", GetTodo).Methods(http.MethodGet)
	router.HandleFunc("/api/todos", CreateTodo).Methods(http.MethodPost)
	router.HandleFunc("/api/todos/{id}", UpdateTodo).Methods(http.MethodPut)
	router.HandleFunc("/api/todos/{id}", DeleteTodo).Methods(http.MethodDelete)
	router.HandleFunc("/api/todos", GetTodos).Methods(http.MethodGet)
	router.HandleFunc("/api/todos", DeleteTodos).Methods(http.MethodDelete)
	router.HandleFunc("/", HTTPHandler)

	err = http.ListenAndServe(":4000", router)
	if err != nil {
		log.Fatal(err)
	}
}
