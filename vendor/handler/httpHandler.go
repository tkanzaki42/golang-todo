package handler

import (
	"fmt"
	"net/http"
)

func HTTPHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, HTTPサーバ")
}
