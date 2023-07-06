package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
)

type TodoResult struct {
	ID   string `json:"id"`
	Task string `json:"task"`
}

type TodoResultWithoutID struct {
	Task string `json:"task"`
}

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "mypassword",
	DB:       0,
})

func GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	todoString, err := client.Get("todo:" + id).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var todo TodoResultWithoutID
	err = json.Unmarshal([]byte(todoString), &todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	getTodoResult := TodoResultWithoutID{
		Task: todo.Task,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(getTodoResult)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo TodoResult
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nextID, err := client.Incr("todo:next_id").Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.ID = strconv.FormatInt(nextID, 10)

	todoJson, err := json.Marshal(todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.Set("todo:"+todo.ID, string(todoJson), 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.RPush("todos", todo.ID).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var newTodo TodoResult
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTodo.ID = id

	newTodoJson, err := json.Marshal(newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = client.Set("todo:"+id, string(newTodoJson), 0).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := client.Del("todo:" + id).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = client.LRem("todos", 0, id).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetTodos(w http.ResponseWriter, r *http.Request) {
	todoIds, err := client.LRange("todos", 0, -1).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todos := []TodoResult{}
	for _, id := range todoIds {
		todo, err := client.Get("todo:" + id).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var t TodoResult
		err = json.Unmarshal([]byte(todo), &t)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		todos = append(todos, t)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func DeleteTodos(w http.ResponseWriter, r *http.Request) {
	err := client.Del("todos").Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, HTTPサーバ")
}
