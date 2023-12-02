package main

import (
	"github.com/seanwiig/tictacgoserver/internal/server"
	"net/http"
)

func main() {
	_ = http.ListenAndServe(":8080", server.NewServer().Router())
}
