package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func TestServer(t *testing.T) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// Hello
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/todo", nil)
	HTTPHandler(w, r)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatal("failed test")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("failed test")
	}
	var expected = "Hello, HTTPサーバ"
	if string(b) != expected {
		t.Fatalf("failed test. expected: %s, actual: %s", expected, string(b))
	}

	// POST /api/todos
	todo := map[string]string{
		"task": "test",
	}
	todoJson, _ := json.Marshal(todo)
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/api/todos", bytes.NewBuffer(todoJson))
	CreateTodo(w, r)
	res = w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Fatal("failed test: ", res.StatusCode)
	}

	// GET /api/todos/1
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/api/todos/1", nil)
	r = mux.SetURLVars(r, map[string]string{
		"id": "1",
	})
	GetTodo(w, r)
	res = w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatal("failed test: ", res.StatusCode)
	}
	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("failed test")
	}

	var todoResult TodoResultWithoutID
	err = json.Unmarshal(b, &todoResult)
	if err != nil {
		t.Fatalf("failed to unmarshal response: %s", err)
	}

	expected = "test"
	if todoResult.Task != expected {
		t.Fatalf("failed test. expected: %s, actual: %s", expected, todoResult.Task)
	}

	// GET /api/todos
}
