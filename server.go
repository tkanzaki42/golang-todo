package main

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
)

func main() {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
	http.HandleFunc("/api/todo", HTTPHandler)
	http.ListenAndServe(":4000", nil)
}
