package main

import (
	"log"
	"net/http"
)

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	return mux
}

func main() {
	log.Fatal(http.ListenAndServe(":8080", NewHandler()))
}