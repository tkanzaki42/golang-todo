package main_test

import (
	"handler"

	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/api/todo", nil)
	handler.HTTPHandler(w, r)
	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		t.Fatal("failed test")
	}
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal("failed test")
	}
	const expected = "Hello, HTTPサーバ"
	if string(b) != expected {
		t.Fatalf("failed test. expected: %s, actual: %s", expected, string(b))
	}
}
